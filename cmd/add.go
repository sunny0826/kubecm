package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
	ac.command.PersistentFlags().BoolP("cover", "c", false, "Overwrite local kubeconfig files")
	_ = ac.command.MarkFlagRequired("file")
}

func (ac *AddCommand) runAdd(cmd *cobra.Command, args []string) error {
	file, _ := ac.command.Flags().GetString("file")
	cover, _ := ac.command.Flags().GetBool("cover")

	var newConfig *clientcmdapi.Config
	var err error

	if file == "-" {
		// from stdin
		contents, err := ioutil.ReadAll(os.Stdin)
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

	err = AddToLocal(newConfig, file, cover)
	if err != nil {
		return err
	}
	return nil
}

// AddToLocal add kubeConfig to local
func AddToLocal(newConfig *clientcmdapi.Config, path string, cover bool) error {
	oldConfig, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	kco := &KubeConfigOption{
		config:   newConfig,
		fileName: getFileName(path),
	}
	// merge context loop
	outConfig, err := kco.handleContexts(oldConfig)
	if err != nil {
		return err
	}
	if len(outConfig.Contexts) == 1 {
		for k := range outConfig.Contexts {
			outConfig.CurrentContext = k
		}
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

func (kc *KubeConfigOption) handleContexts(oldConfig *clientcmdapi.Config) (*clientcmdapi.Config, error) {
	newConfig := clientcmdapi.NewConfig()
	for name, ctx := range kc.config.Contexts {
		var newName string
		if len(kc.config.Contexts) >= 1 {
			newName = fmt.Sprintf("%s-%s", kc.fileName, HashSufString(name))
		} else {
			newName = kc.fileName
		}
		if checkContextName(newName, oldConfig) {
			nameConfirm := BoolUI(fmt.Sprintf("「%s」 Name already exists, do you want to rename it. (If you select `False`, this context will not be merged)", newName))
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
# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf |  kubecm add -f -
`
}
