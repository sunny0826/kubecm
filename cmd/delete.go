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
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the specified context from the kubeconfig",
	Long:  `Delete the specified context from the kubeconfig`,
	Example: `
  # Delete the context
  kubecm delete my-context
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			config, err := LoadClientConfig(cfgFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			err = deleteContext(args, config)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		} else {
			fmt.Println("Please enter the context you want to delete.")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.SetArgs([]string{""})
}

func deleteContext(ctxs []string, config *clientcmdapi.Config) error {
	for key, _ := range config.Contexts {
		for _, ctx := range ctxs {
			if ctx == key {
				delete(config.Contexts, key)
				break
			}
		}
	}
	err := ModifyKubeConfig(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Context delete succeeded!\nDelete: %v \n", ctxs)
	return nil
}
