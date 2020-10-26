package cmd

import (
	"fmt"
	"io/ioutil"

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
		Short:   "Merge the kubeconfig files in the specified directory",
		Long:    `Merge the kubeconfig files in the specified directory`,
		Aliases: []string{"m"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return mc.runMerge(cmd, args)
		},
		Example: mergeExample(),
	}
	mc.command.Flags().StringP("folder", "f", "", "Kubeconfig folder")
	mc.command.Flags().BoolP("cover", "c", false, "Overwrite the original kubeconfig file")
	_ = mc.command.MarkFlagRequired("folder")
}

func (mc MergeCommand) runMerge(command *cobra.Command, args []string) error {
	folder, _ := mc.command.Flags().GetString("folder")
	files := listFile(folder)
	mc.command.Printf("Loading kubeconfig file: %v \n", files)
	configs := clientcmdapi.NewConfig()
	for _, yaml := range files {
		config, err := formatNewConfig(yaml, "")
		if err != nil {
			return err
		}
		configs = appendConfig(configs, config)
	}
	cover, _ := mc.command.Flags().GetBool("cover")
	err := WriteConfig(cover, folder, configs)
	if err != nil {
		return err
	}
	return nil
}

func listFile(folder string) []string {
	files, _ := ioutil.ReadDir(folder)
	var flist []string
	for _, file := range files {
		if file.IsDir() {
			listFile(folder + "/" + file.Name())
		} else {
			flist = append(flist, fmt.Sprintf("%s/%s", folder, file.Name()))
		}
	}
	return flist
}

func mergeExample() string {
	return `
# Merge kubeconfig in the dir directory
kubecm merge -f dir

# Merge kubeconfig in the directory and overwrite the original kubeconfig file
kubecm merge -f dir -c
`
}
