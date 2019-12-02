/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Sets the current-context in a kubeconfig file(Will be removed in new version, please use kubecm swtich)",
	Example: `
# Use the context for the test cluster
kubecm use test
`,
	Long: `
Sets the current-context in a kubeconfig file(Will be removed in new version, please use kubecm swtich)
`,
	Run: func(cmd *cobra.Command, args []string) {
		Warning.Println("This command Will be removed in new version, please use kubecm swtich.")
		if len(args) == 1 {
			context := args[0]
			config, err := LoadClientConfig(cfgFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			var currentContext bool
			for key, _ := range config.Contexts {
				if key == context {
					config.CurrentContext = key
					currentContext = true
					fmt.Println(fmt.Sprintf("Switched to context %s", context))
				}
			}
			if !currentContext {
				fmt.Println(fmt.Sprintf("no context exists with the name: %s", context))
				os.Exit(1)
			} else {
				err := ModifyKubeConfig(config)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				err = Formatable(nil)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("Please input a CONTEXT_NAME.")
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
	useCmd.SetArgs([]string{""})
}
