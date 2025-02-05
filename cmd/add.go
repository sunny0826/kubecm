package cmd

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// AddCommand add command struct
type AddCommand struct {
	BaseCommand
}

// KubeConfigOption kubeConfig option
type KubeConfigOption struct {
	config                *clientcmdapi.Config
	fileName              string
	insecureSkipTLSVerify bool
}

// Init AddCommand
func (ac *AddCommand) Init() {
	ac.command = &cobra.Command{
		Use:     "add",
		Short:   "Add KubeConfig to $HOME/.kube/config",
		Long:    "Add KubeConfig to $HOME/.kube/config",
		Aliases: []string{"a"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.runAdd(cmd, args)
		},
		Example: addExample(),
	}
	ac.command.PersistentFlags().BoolP("cover", "c", false, "overwrite local kubeconfig files")
	ac.command.Flags().StringP("file", "f", "", "path to merge kubeconfig files")
	ac.command.Flags().StringSlice("context", []string{}, "specify the context to be added")
	ac.command.Flags().String("context-prefix", "", "add a prefix before context name")
	ac.command.Flags().String("context-name", "", "override context name when add kubeconfig context, when context-name is set, context-prefix and context-template parameters will be ignored")
	ac.command.Flags().StringSlice("context-template", []string{"context"}, "define the attributes used for composing the context name, available values: filename, user, cluster, context, namespace")
	ac.command.Flags().Bool("select-context", false, "select the context to be added in interactive mode")
	ac.command.Flags().Bool("insecure-skip-tls-verify", false, "if true, the server's certificate will not be checked for validity")
	_ = ac.command.MarkFlagRequired("file")
	ac.AddCommands(&DocsCommand{})
}

func (ac *AddCommand) runAdd(cmd *cobra.Command, args []string) error {
	cover, _ := ac.command.Flags().GetBool("cover")
	file, _ := ac.command.Flags().GetString("file")
	context, _ := ac.command.Flags().GetStringSlice("context")
	contextPrefix, _ := ac.command.Flags().GetString("context-prefix")
	contextName, _ := ac.command.Flags().GetString("context-name")
	contextTemplate, _ := ac.command.Flags().GetStringSlice("context-template")
	selectContext, _ := ac.command.Flags().GetBool("select-context")
	insecureSkipTLSVerify, _ := ac.command.Flags().GetBool("insecure-skip-tls-verify")

	var newConfig *clientcmdapi.Config

	if contextName != "" {
		contextTemplate = []string{}
		contextPrefix = contextName
	}
	err := validateContextTemplate(contextTemplate)
	if err != nil {
		return err
	}

	if file == "-" {
		// from stdin
		contents, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		newConfig, err = clientcmd.Load(contents)
		if err != nil {
			return err
		}
	} else {
		// check path
		file, err := CheckAndTransformFilePath(file, cfgCreate)
		if err != nil {
			return err
		}
		newConfig, err = clientcmd.LoadFromFile(file)
		if err != nil {
			return err
		}
	}

	err = AddToLocal(newConfig, file, contextPrefix, cover, selectContext, contextTemplate, context, insecureSkipTLSVerify)
	if err != nil {
		return err
	}
	return nil
}

// AddToLocal add kubeConfig to local
func AddToLocal(newConfig *clientcmdapi.Config, path, contextPrefix string, cover bool, selectContext bool, contextTemplate []string, context []string, insecureSkipTLSVerify bool) error {
	oldConfig, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	kco := &KubeConfigOption{
		config:                newConfig,
		fileName:              getFileName(path),
		insecureSkipTLSVerify: insecureSkipTLSVerify,
	}
	// merge context loop
	outConfig, err := kco.handleContexts(oldConfig, contextPrefix, selectContext, contextTemplate, context)
	if err != nil {
		return err
	}
	if len(outConfig.Contexts) == 1 {
		for k := range outConfig.Contexts {
			outConfig.CurrentContext = k
		}
	}

	if reflect.DeepEqual(oldConfig, outConfig) {
		fmt.Println("No context to add.")
		return nil
	}

	if !cover {
		cover, err = strconv.ParseBool(BoolUI(fmt.Sprintf("Does it overwrite File 「%s」?", cfgFile)))
		if err != nil {
			return err
		}
	}
	err = WriteConfig(cover, path, outConfig)
	if err != nil {
		return err
	}
	return nil
}

func (kc *KubeConfigOption) handleContexts(oldConfig *clientcmdapi.Config, contextPrefix string, selectContext bool, contextTemplate []string, context []string) (*clientcmdapi.Config, error) {
	newConfig := clientcmdapi.NewConfig()
	var newName string
	generatedName := make(map[string]int)

	for name, ctx := range kc.config.Contexts {
		newName = kc.generateContextName(name, ctx, contextTemplate)
		if contextPrefix != "" {
			newName = strings.TrimSuffix(fmt.Sprintf("%s-%s", contextPrefix, newName), "-")
		}

		// prevent generate duplicate context name
		// for example: set --context-template user,cluster, the context1 and context2 have the same user and cluster
		generatedName[newName]++
		if generatedName[newName] > 1 {
			newName = fmt.Sprintf("%s-%d", newName, generatedName[newName])
		}

		if len(context) > 0 {
			if !slices.Contains(context, newName) {
				continue
			}
		} else if selectContext {
			importContext := BoolUI(fmt.Sprintf("Do you want to add context「%s」? (If you select `False`, this context will not be merged)", newName))
			if importContext == "False" {
				continue
			}
		}

		var quitNewName bool
		for checkContextName(newName, oldConfig) {
			nameConfirm := BoolUI(fmt.Sprintf("「%s」 Name already exists, do you want to rename it? (If you select `False`, this context will not be merged)", newName))
			if nameConfirm == "True" {
				newName = PromptUI("Rename", newName)
				continue
			} else {
				quitNewName = true
				break
			}
		}
		if quitNewName {
			continue
		}
		itemConfig := kc.handleContext(oldConfig, newName, ctx)
		newConfig = appendConfig(newConfig, itemConfig)
		fmt.Printf("Add Context: %s \n", newName)
	}
	outConfig := appendConfig(oldConfig, newConfig)
	return outConfig, nil
}

func (kc *KubeConfigOption) generateContextName(name string, ctx *clientcmdapi.Context, contextTemplate []string) string {
	valueMap := map[string]string{
		Filename:  kc.fileName,
		Context:   name,
		User:      ctx.AuthInfo,
		Cluster:   ctx.Cluster,
		Namespace: ctx.Namespace,
	}

	var contextValues []string
	for _, value := range contextTemplate {
		if v, ok := valueMap[value]; ok {
			if v != "" {
				contextValues = append(contextValues, v)
			}
		}
	}

	return strings.Join(contextValues, "-")
}

func checkContextName(name string, oldConfig *clientcmdapi.Config) bool {
	if _, ok := oldConfig.Contexts[name]; ok {
		return true
	}
	return false
}

// checkClusterAndUserName check if the cluster and user exist in the same context, or just cluster or user exist
func checkClusterAndUserName(oldConfig *clientcmdapi.Config, newClusterName, newUserName string) (bool, bool, bool) {
	var (
		clusterAndUserNameExistInSameContext bool
		justClusterNameExist                 bool
		justUserNameExist                    bool
	)

	for _, ctx := range oldConfig.Contexts {
		if ctx.Cluster == newClusterName && ctx.AuthInfo == newUserName {
			clusterAndUserNameExistInSameContext = true
		}
		if ctx.Cluster == newClusterName {
			justClusterNameExist = true
		}
		if ctx.AuthInfo == newUserName {
			justUserNameExist = true
		}
	}

	return clusterAndUserNameExistInSameContext, justClusterNameExist, justUserNameExist
}

// isSameKubeConfigAlreadyExist return true if the same kubeconfig is already
// exist, assert by cluster, user name and corresponding cluster and user info
func (kc *KubeConfigOption) isSameKubeConfigAlreadyExist(oldConfig *clientcmdapi.Config, ctx *clientcmdapi.Context) bool {
	oldCluster, ok := oldConfig.Clusters[ctx.Cluster]
	if !ok {
		return false
	}

	oldUser, ok := oldConfig.AuthInfos[ctx.AuthInfo]
	if !ok {
		return false
	}

	return reflect.DeepEqual(oldCluster, kc.config.Clusters[ctx.Cluster]) && reflect.DeepEqual(oldUser, kc.config.AuthInfos[ctx.AuthInfo])
}

func (kc *KubeConfigOption) handleContext(oldConfig *clientcmdapi.Config,
	name string, ctx *clientcmdapi.Context) *clientcmdapi.Config {

	var (
		clusterNameSuffix string
		userNameSuffix    string
	)

	newConfig := clientcmdapi.NewConfig()
	bothExistInSameContext, justClusterNameExist, justUserNameExist := checkClusterAndUserName(oldConfig, ctx.Cluster, ctx.AuthInfo)
	// if same kubeconfig is already exist, skip it
	if bothExistInSameContext {
		if kc.isSameKubeConfigAlreadyExist(oldConfig, ctx) {
			return newConfig
		}
	}
	suffix := rand.String(10)

	if justClusterNameExist {
		clusterNameSuffix = "-" + suffix
	}
	if justUserNameExist {
		userNameSuffix = "-" + suffix
	}

	userName := fmt.Sprintf("%v%v", ctx.AuthInfo, userNameSuffix)
	clusterName := fmt.Sprintf("%v%v", ctx.Cluster, clusterNameSuffix)
	newCtx := ctx.DeepCopy()

	// deep copy and clear CA data
	cluster := kc.config.Clusters[newCtx.Cluster].DeepCopy()
	if kc.insecureSkipTLSVerify {
		cluster.InsecureSkipTLSVerify = true
		cluster.CertificateAuthority = ""
		cluster.CertificateAuthorityData = nil
	}

	var clusterInfoExist, userInfoExist bool
	for _, oldCluster := range oldConfig.Clusters {
		if reflect.DeepEqual(oldCluster, cluster) {
			clusterInfoExist = true
		}
	}
	if !clusterInfoExist {
		newConfig.Clusters[clusterName] = cluster
	}

	for _, oldUser := range oldConfig.AuthInfos {
		if reflect.DeepEqual(oldUser, kc.config.AuthInfos[newCtx.AuthInfo]) {
			userInfoExist = true
		}
	}
	if !userInfoExist {
		newConfig.AuthInfos[userName] = kc.config.AuthInfos[newCtx.AuthInfo]
	}

	newConfig.Contexts[name] = newCtx
	newConfig.Contexts[name].AuthInfo = userName
	newConfig.Contexts[name].Cluster = clusterName

	return newConfig
}

func addExample() string {
	return `
# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml
# Merge test.yaml with $HOME/.kube/config and add a prefix before context name
kubecm add -cf test.yaml --context-prefix test
# Merge test.yaml with $HOME/.kube/config and define the attributes used for composing the context name
kubecm add -f test.yaml --context-template user,cluster
# Merge test.yaml with $HOME/.kube/config, define the attributes used for composing the context name and add a prefix before context name
kubecm add -f test.yaml --context-template user,cluster --context-prefix demo
# Merge test.yaml with $HOME/.kube/config and override context name, it's useful if there is only one context in the kubeconfig file
kubecm add -f test.yaml --context-name test
# Merge test.yaml with $HOME/.kube/config and select the context to be added in interactive mode
kubecm add -f test.yaml --select-context
# Merge test.yaml with $HOME/.kube/config and specify the context to be added
kubecm add -f test.yaml --context context1,context2
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf | kubecm add -f -
# Merge test.yaml with $HOME/.kube/config and skip TLS certificate verification
kubecm add -f test.yaml --insecure-skip-tls-verify
`
}
