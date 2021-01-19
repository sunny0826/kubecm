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
		Short:   "Merge the KubeConfig files in the specified directory",
		Long:    `Merge the KubeConfig files in the specified directory`,
		Aliases: []string{"m"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return mc.runMerge(cmd, args)
		},
		Example: mergeExample(),
	}
	mc.command.Flags().StringP("folder", "f", "", "KubeConfig folder")
	_ = mc.command.MarkFlagRequired("folder")
}

func (mc MergeCommand) runMerge(command *cobra.Command, args []string) error {
	folder, _ := mc.command.Flags().GetString("folder")
	folder, err := CheckAndTransformFilePath(folder)
	if err != nil {
		return err
	}
	files := listFile(folder)
	outConfigs := clientcmdapi.NewConfig()
	for _, yaml := range files {
		printString(os.Stdout, "Loading KubeConfig file:"+yaml+" \n")
		loadConfig, err := clientcmd.LoadFromFile(yaml)
		if err != nil {
			return err
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
	cover := BoolUI(fmt.Sprintf("Are you sure you want to overwrite the「%s」file?", cfgFile))
	confirm, err := strconv.ParseBool(cover)
	if err != nil {
		return err
	}
	err = WriteConfig(confirm, folder, outConfigs)
	if err != nil {
		return err
	}
	return nil
}

func listFile(folder string) []string {
	files, _ := ioutil.ReadDir(folder)
	var fileList []string
	for _, file := range files {
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
# Merge KubeConfig in the dir directory
kubecm merge -f dir
`
}
