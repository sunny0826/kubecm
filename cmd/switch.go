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
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch Kube Context interactively.",
	Example: `
# Switch Kube Context interactively
kubecm switch
`,
	Long: `
Switch Kube Context interactively.
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadClientConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		var kubeItems []string
		current := config.CurrentContext
		for key, _ := range config.Contexts {
			if key != current {
				kubeItems = append(kubeItems, key)
			} else {
				kubeItems = append([]string{current}, kubeItems...)
			}
		}
		prompt := promptui.Select{
			Label: "Select Kube Context",
			Items: kubeItems,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		config.CurrentContext = result
		err = ModifyKubeConfig(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Switched to context %s\n", result)
		err = Formatable(nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
