package cmd

import (
	"fmt"
	"os"
	"strconv"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// MergeCommand merge cmd struct
type MergeCommand struct {
	BaseCommand
}

// Init MergeCommand
func (mc *MergeCommand) Init() {
	mc.command = &cobra.Command{
		Use:     "merge",
		Short:   "Merge multiple kubeconfig files into one",
		Long:    `Merge multiple kubeconfig files into one`,
		Aliases: []string{"m"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return mc.runMerge(cmd, args)
		},
		Example: mergeExample(),
	}
	mc.command.Flags().StringP("folder", "f", "", "KubeConfig folder")
	mc.command.Flags().BoolP("assumeyes", "y", false, "skip interactive file overwrite confirmation")
	mc.command.Flags().String("context-prefix", "", "add a prefix before context name")
	mc.command.Flags().Bool("select-context", false, "select the context to be merged")
	mc.command.Flags().StringSlice("context-template", []string{"context"}, "define the attributes used for composing the context name, available values: filename, user, cluster, context, namespace")
	//_ = mc.command.MarkFlagRequired("folder")
	mc.AddCommands(&DocsCommand{})
}

func (mc MergeCommand) runMerge(command *cobra.Command, args []string) error {
	files := args
	folder, _ := mc.command.Flags().GetString("folder")
	contextPrefix, _ := mc.command.Flags().GetString("context-prefix")
	selectContext, _ := mc.command.Flags().GetBool("select-context")
	contextTemplate, _ := mc.command.Flags().GetStringSlice("context-template")

	err := validateContextTemplate(contextTemplate)
	if err != nil {
		return err
	}

	if folder != "" {
		folder, err = CheckAndTransformFilePath(folder)
		if err != nil {
			return err
		}
		files = append(files, listFile(folder)...)
	}
	if len(files) == 0 {
		return fmt.Errorf("please enter the files to be merged")
	}
	outConfigs := clientcmdapi.NewConfig()
	for _, yaml := range files {
		printString(os.Stdout, "Loading KubeConfig file: "+yaml+" \n")
		loadConfig, err := loadKubeConfig(yaml)
		if err != nil {
			// If an error is reported, the loading of this file is skipped.
			printWarning(os.Stdout, "File "+yaml+" is not kubeconfig\n")
			continue
		}
		kco := &KubeConfigOption{
			config:   loadConfig,
			fileName: getFileName(yaml),
		}
		outConfigs, err = kco.handleContexts(outConfigs, contextPrefix, selectContext, contextTemplate)
		if err != nil {
			return err
		}
	}

	if len(outConfigs.Contexts) == 0 {
		fmt.Println("No context to merge.")
		return nil
	}

	confirm, _ := mc.command.Flags().GetBool("assumeyes")
	if !confirm {
		cover := BoolUI(fmt.Sprintf("Are you sure you want to overwrite the 「%s」 file?", cfgFile))
		confirm, _ = strconv.ParseBool(cover)
	}
	err = WriteConfig(confirm, cfgFile, outConfigs)
	if err != nil {
		return err
	}
	return MacNotifier("Merge Successfully")
}

func loadKubeConfig(yaml string) (*clientcmdapi.Config, error) {
	loadConfig, err := clientcmd.LoadFromFile(yaml)
	if err != nil {
		return nil, err
	}
	if len(loadConfig.Contexts) == 0 {
		return nil, fmt.Errorf("no kubeconfig in %s ", yaml)
	}
	return loadConfig, err
}

func listFile(folder string) []string {
	files, _ := os.ReadDir(folder)
	var fileList []string
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		if file.IsDir() {
			listFile(folder + "/" + file.Name())
		} else {
			fileList = append(fileList, fmt.Sprintf("%s/%s", folder, file.Name()))
		}
	}
	return fileList
}

func mergeExample() string {
	return `
# Merge multiple kubeconfig
kubecm merge 1st.yaml 2nd.yaml 3rd.yaml
# Merge KubeConfig in the dir directory
kubecm merge -f dir
# Merge KubeConfig in the dir directory to the specified file.
kubecm merge -f dir --config kubecm.config
# Merge test.yaml with $HOME/.kube/config and add a prefix before context name
kubecm merge test.yaml --context-prefix test
# Merge test.yaml with $HOME/.kube/config and define the attributes used for composing the context name
kubecm merge test.yaml --context-template user,cluster
# Merge test.yaml with $HOME/.kube/config, define the attributes used for composing the context name and add a prefix before context name
kubecm merge test.yaml --context-template user,cluster --context-prefix demo
# Merge test.yaml with $HOME/.kube/config and select the context to be added in interactive mode
kubecm merge test.yaml --select-context
`
}
