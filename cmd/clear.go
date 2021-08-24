package cmd

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// ClearCommand clean command struct
type ClearCommand struct {
	BaseCommand
}

// Init ClearCommand
func (cl *ClearCommand) Init() {
	cl.command = &cobra.Command{
		Use:   "clear",
		Short: "Clear lapsed context, cluster and user",
		Long:  "Clear lapsed context, cluster and user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runClear(cmd, args)
		},
		Example: clearExample(),
	}
	cl.command.DisableFlagsInUseLine = true
}

func (cl *ClearCommand) runClear(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		if ok, err := clearContext(cfgFile); ok {
			printString(os.Stdout, fmt.Sprintf("There is nothing to clean in 「%s」\n", cfgFile))
			return nil
		} else if err != nil {
			return err
		}
	} else {
		for _, file := range args {
			if ok, err := clearContext(file); ok {
				printString(os.Stdout, fmt.Sprintf("There is nothing to clean in 「%s」\n", file))
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}

func clearContext(file string) (bool, error) {
	config, err := clientcmd.LoadFromFile(file)
	if err != nil {
		return false, err
	}
	outConfig := config.DeepCopy()
	outConfig = CheckValidContext(true, outConfig)
	if reflect.DeepEqual(config, outConfig) {
		return true, nil
	}
	err = UpdateConfigFile(file, outConfig)
	if err != nil {
		return false, err
	}
	return false, nil
}

func clearExample() string {
	return fmt.Sprintf(`
# Clear lapsed context, cluster and user (default is %s)
kubecm clear
# Customised clear lapsed files
kubecm clear config.yaml test.yaml
`, cfgFile)
}
