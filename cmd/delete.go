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
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"log"
)

type DeleteCommand struct {
	baseCommand
}

func (dc *DeleteCommand) Init() {
	dc.command = &cobra.Command{
		Use:     "delete",
		Short:   "Delete the specified context from the kubeconfig",
		Long:    `Delete the specified context from the kubeconfig`,
		Aliases: []string{"d"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return dc.runDelete(cmd, args)
		},
		Example: deleteExample(),
	}

}

func (dc *DeleteCommand) runDelete(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if len(args) != 0 {
		if err != nil {
			log.Fatal(err)
		}
		err = dc.deleteContext(args, config)
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
		// exit option
		kubeItems, err := ExitOption(kubeItems)
		if err != nil {
			return err
		}
		num := SelectUI(kubeItems, "Select The Delete Kube Context")
		kubeName := kubeItems[num].Name
		confirm := BoolUI(fmt.Sprintf("Are you sure you want to delete「%s」?", kubeName))
		if confirm == "True" {
			err = dc.deleteContext([]string{kubeName}, config)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			dc.command.Println("Nothing deleted！")
		}
	} else {
		dc.command.Println("Please enter the context you want to delete.")
	}
	return nil
}

func (dc *DeleteCommand) deleteContext(ctxs []string, config *clientcmdapi.Config) error {
	for _, ctx := range ctxs {
		if _, ok := config.Contexts[ctx]; ok {
			delete(config.Contexts, ctx)
			dc.command.Printf("Context Delete:「%s」\n", ctx)
		} else {
			Error.Printf("「%s」do not exit.", ctx)
		}
	}
	err := dc.WriteConfig(true,cfgFile,config)
	if err != nil {
		return err
	}
	return nil
}

func deleteExample() string {
	return `
# Delete the context interactively
kubecm delete
# Delete the context
kubecm delete my-context
`
}
