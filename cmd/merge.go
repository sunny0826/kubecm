package cmd

import (
	"fmt"
	"io/ioutil"
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
	//_ = mc.command.MarkFlagRequired("folder")
}

func (mc MergeCommand) runMerge(command *cobra.Command, args []string) error {
	files := args
	folder, _ := mc.command.Flags().GetString("folder")
	if folder != "" {
		folder, err := CheckAndTransformFilePath(folder)
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
		outConfigs, err = kco.handleContexts(outConfigs)
		if err != nil {
			return err
		}
	}
	confirm, _ := mc.command.Flags().GetBool("assumeyes")
	if !confirm {
		cover := BoolUI(fmt.Sprintf("Are you sure you want to overwrite the 「%s」 file?", cfgFile))
		confirm, _ = strconv.ParseBool(cover)
	}
	err := WriteConfig(confirm, cfgFile, outConfigs)
	if err != nil {
		return err
	}
	return nil
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
	files, _ := ioutil.ReadDir(folder)
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
`
}
