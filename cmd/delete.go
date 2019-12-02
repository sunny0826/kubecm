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
	"log"
	"os"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the specified context from the kubeconfig",
	Long:  `Delete the specified context from the kubeconfig`,
	Example: `
# Delete the context interactively
kubecm delete
# Delete the context
kubecm delete my-context
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadClientConfig(cfgFile)
		if len(args) != 0 {
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			err = deleteContext(args, config)
			if err != nil {
				Error.Println(err)
				os.Exit(-1)
			}
		} else if len(args) == 0 {
			var kubeItems []needle
			for key, obj := range config.Contexts {
				if key != config.CurrentContext {
					kubeItems = append(kubeItems, needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
				} else {
					kubeItems = append([]needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
				}
			}
			num := SelectUI(kubeItems, "Select The Delete Kube Context")
			kubeName := kubeItems[num].Name
			err = deleteContext([]string{kubeName}, config)
			if err != nil {
				Error.Println(err)
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
	for _, ctx := range ctxs {
		if _, ok := config.Contexts[ctx]; ok {
			delete(config.Contexts, ctx)
			log.Printf("Context Delete: %s \n", ctx)
		} else {
			Error.Printf("「%s」do not exit.", ctx)
		}
	}
	err := ModifyKubeConfig(config)
	if err != nil {
		Error.Println(err)
		os.Exit(1)
	}
	return nil
}
