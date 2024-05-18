package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// AddCommand add command struct
type AddCommand struct {
	BaseCommand
}

// KubeConfigOption kubeConfig option
type KubeConfigOption struct {
	config   *clientcmdapi.Config
	fileName string
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
	ac.command.Flags().StringP("file", "f", "", "Path to merge kubeconfig files")
	ac.command.Flags().String("context-name", "", "override context name when add kubeconfig context")
	ac.command.PersistentFlags().BoolP("cover", "c", false, "Overwrite local kubeconfig files")
	ac.command.PersistentFlags().Bool("select-context", false, "select the context to be added")
	_ = ac.command.MarkFlagRequired("file")
	ac.AddCommands(&DocsCommand{})
}

func (ac *AddCommand) runAdd(cmd *cobra.Command, args []string) error {
	file, _ := ac.command.Flags().GetString("file")
	cover, _ := ac.command.Flags().GetBool("cover")
	contextName, _ := ac.command.Flags().GetString("context-name")
	selectContext, _ := ac.command.Flags().GetBool("select-context")

	var newConfig *clientcmdapi.Config
	var err error

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
		file, err := CheckAndTransformFilePath(file)
		if err != nil {
			return err
		}
		newConfig, err = clientcmd.LoadFromFile(file)
		if err != nil {
			return err
		}
	}

	err = AddToLocal(newConfig, file, contextName, cover, selectContext)
	if err != nil {
		return err
	}
	return nil
}

// AddToLocal add kubeConfig to local
func AddToLocal(newConfig *clientcmdapi.Config, path, newName string, cover bool, selectContext bool) error {
	oldConfig, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	kco := &KubeConfigOption{
		config:   newConfig,
		fileName: getFileName(path),
	}
	// merge context loop
	outConfig, err := kco.handleContexts(oldConfig, newName, selectContext)
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

func (kc *KubeConfigOption) handleContexts(oldConfig *clientcmdapi.Config, contextName string, selectContext bool) (*clientcmdapi.Config, error) {
	newConfig := clientcmdapi.NewConfig()
	var index int
	var newName string
	for name, ctx := range kc.config.Contexts {
		if len(kc.config.Contexts) > 1 {
			if contextName == "" {
				newName = fmt.Sprintf("%s-%s", kc.fileName, HashSufString(name))
			} else {
				newName = fmt.Sprintf("%s-%d", contextName, index)
			}
		} else if contextName == "" {
			newName = name
		} else {
			newName = contextName
		}

		if selectContext {
			importContext := BoolUI(fmt.Sprintf("Do you want to add context「%s」? (If you select `False`, this context will not be merged)", newName))
			if importContext == "False" {
				continue
			}
		}

		if checkContextName(newName, oldConfig) {
			nameConfirm := BoolUI(fmt.Sprintf("「%s」 Name already exists, do you want to rename it? (If you select `False`, this context will not be merged)", newName))
			if nameConfirm == "True" {
				newName = PromptUI("Rename", newName)
				if newName == kc.fileName {
					return nil, errors.New("need to rename")
				}
			} else {
				continue
			}
		}
		itemConfig := kc.handleContext(oldConfig, newName, ctx)
		newConfig = appendConfig(newConfig, itemConfig)
		fmt.Printf("Add Context: %s \n", newName)
		index++
	}
	outConfig := appendConfig(oldConfig, newConfig)
	return outConfig, nil
}

func checkContextName(name string, oldConfig *clientcmdapi.Config) bool {
	if _, ok := oldConfig.Contexts[name]; ok {
		return true
	}
	return false
}

func checkClusterAndUserName(oldConfig *clientcmdapi.Config, newClusterName, newUserName string) (bool, bool) {
	var (
		isClusterNameExist bool
		isUserNameExist    bool
	)

	for _, ctx := range oldConfig.Contexts {
		if ctx.Cluster == newClusterName {
			isClusterNameExist = true
		}
		if ctx.AuthInfo == newUserName {
			isUserNameExist = true
		}
	}

	return isClusterNameExist, isUserNameExist
}

func (kc *KubeConfigOption) handleContext(oldConfig *clientcmdapi.Config,
	name string, ctx *clientcmdapi.Context) *clientcmdapi.Config {

	var (
		clusterNameSuffix string
		userNameSuffix    string
	)

	isClusterNameExist, isUserNameExist := checkClusterAndUserName(oldConfig, ctx.Cluster, ctx.AuthInfo)
	newConfig := clientcmdapi.NewConfig()
	suffix := HashSufString(name)

	if isClusterNameExist {
		clusterNameSuffix = "-" + suffix
	}
	if isUserNameExist {
		userNameSuffix = "-" + suffix
	}

	userName := fmt.Sprintf("%v%v", ctx.AuthInfo, userNameSuffix)
	clusterName := fmt.Sprintf("%v%v", ctx.Cluster, clusterNameSuffix)
	newCtx := ctx.DeepCopy()
	newConfig.AuthInfos[userName] = kc.config.AuthInfos[newCtx.AuthInfo]
	newConfig.Clusters[clusterName] = kc.config.Clusters[newCtx.Cluster]
	newConfig.Contexts[name] = newCtx
	newConfig.Contexts[name].AuthInfo = userName
	newConfig.Contexts[name].Cluster = clusterName

	return newConfig
}

func addExample() string {
	return `
Note: If -c is set and more than one context is added to the kubeconfig file, the following will occur:
- If --context-name is set, the context will be generated as <context-name-0>, <context-name-1> ...
- If --context-name is not set, it will be generated as <file-name-{hash}> where {hash} is the MD5 hash of the file name.

# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
# Merge test.yaml with $HOME/.kube/config and rename context name
kubecm add -cf test.yaml --context-name test
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf |  kubecm add -f -
`
}
