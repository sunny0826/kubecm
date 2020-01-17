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
	"errors"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var oldName string
var newName string

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename the contexts of kubeconfig",
	Long: `
# Renamed the context interactively
kubecm rename
# Renamed dev to test
kubecm rename -o dev -n test
# Renamed current-context name to dev
kubecm rename -n dev -c
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadClientConfig(cfgFile)
		if newName == "" && oldName == "" {
			var kubeItems []needle
			for key, obj := range config.Contexts {
				if key != config.CurrentContext {
					kubeItems = append(kubeItems, needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
				} else {
					kubeItems = append([]needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
				}
			}
			num := SelectUI(kubeItems, "Select The Rename Kube Context")
			kubeName := kubeItems[num].Name
			rename := InputStr(kubeName)
			if rename != kubeName {
				if _, ok := config.Contexts[rename]; ok {
					log.Fatal("Name: %s already exists", rename)
				} else {
					if obj, ok := config.Contexts[kubeName]; ok {
						config.Contexts[rename] = obj
						delete(config.Contexts, kubeName)
						if config.CurrentContext == kubeName {
							config.CurrentContext = rename
						}
					}
					err = ModifyKubeConfig(config)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				log.Fatalf("No name: %s changes", rename)
			}
		} else {
			cover, _ = cmd.Flags().GetBool("cover")
			if cover && oldName != "" {
				log.Fatalln("parameter `-c` and `-n` cannot be set at the same time")
			} else {
				if err != nil {
					log.Fatal(err)
				}
				if _, ok := config.Contexts[newName]; ok {
					log.Fatalf("the name:%s is exit.", newName)
				}
				if cover {
					for key, obj := range config.Contexts {
						if current := config.CurrentContext; key == current {
							config.Contexts[newName] = obj
							delete(config.Contexts, key)
							config.CurrentContext = newName
							cmd.Printf("Rename %s to %s", key, newName)
							break
						}
					}
				} else {
					if obj, ok := config.Contexts[oldName]; ok {
						config.Contexts[newName] = obj
						delete(config.Contexts, oldName)
						if config.CurrentContext == oldName {
							config.CurrentContext = newName
						}
					} else {
						cmd.Printf("Can not find context: %s", oldName)
						err := Formatable(nil)
						if err != nil {
							log.Fatal(err)
						}
						os.Exit(-1)
					}
				}
				err = ModifyKubeConfig(config)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		err = Formatable(nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&oldName, "old", "o", "", "Old context name")
	renameCmd.Flags().StringVarP(&newName, "new", "n", "", "New context name")
	renameCmd.Flags().BoolP("cover", "c", false, "")
}

func InputStr(name string) string {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Context name must have more than 3 characters")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Rename",
		Validate: validate,
		Default:  name,
	}
	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}
