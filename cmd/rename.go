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
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type RenameCommand struct {
	baseCommand
}

func (rc *RenameCommand) Init() {
	rc.command = &cobra.Command{
		Use:     "rename",
		Short:   "Rename the contexts of kubeconfig",
		Long:    "Rename the contexts of kubeconfig",
		Aliases: []string{"r"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return rc.runRename(cmd, args)
		},
		Example: renameExample(),
	}
}

func (rc *RenameCommand) runRename(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	var kubeItems []needle
	for key, obj := range config.Contexts {
		if key != config.CurrentContext {
			kubeItems = append(kubeItems, needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
		} else {
			kubeItems = append([]needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
		}
	}
	// exit option
	kubeItems, err = ExitOption(kubeItems)
	if err != nil {
		return err
	}
	num := SelectUI(kubeItems, "Select The Rename Kube Context")
	kubeName := kubeItems[num].Name
	rename := PromptUI("Rename", kubeName)
	if rename != kubeName {
		if _, ok := config.Contexts[rename]; ok {
			log.Fatalf("Name: %s already exists", rename)
		} else {
			if obj, ok := config.Contexts[kubeName]; ok {
				config.Contexts[rename] = obj
				delete(config.Contexts, kubeName)
				if config.CurrentContext == kubeName {
					config.CurrentContext = rename
				}
			}
			err = rc.WriteConfig(true,cfgFile,config)
			if err != nil {
				return err
			}
		}
	} else {
		rc.command.Printf("No name: %s changes\n", rename)
	}
	return nil
}

func renameExample() string {
	return `
# Renamed the context interactively
kubecm rename
`
}
