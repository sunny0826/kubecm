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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"log"
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
				log.Fatal(err)
			}
			err = deleteContext(args, config)
			if err != nil {
				log.Fatal(err)
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
			confirm := BoolUI(fmt.Sprintf("Are you sure you want to delete「%s」?", kubeName))
			if confirm == "True" {
				err = deleteContext([]string{kubeName}, config)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				cmd.Println("Nothing deleted！")
			}
		} else {
			cmd.Println("Please enter the context you want to delete.")
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
			log.Printf("Context Delete:「%s」\n", ctx)
		} else {
			Error.Printf("「%s」do not exit.", ctx)
		}
	}
	err := ModifyKubeConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func BoolUI(label string) string {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F37A {{ . | red }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F47B {{ . | green }}",
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     []string{"True", "False"},
		Templates: templates,
		Size:      2,
	}
	_, obj, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return obj
}
