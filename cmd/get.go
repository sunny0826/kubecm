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
	"github.com/bndr/gotabulate"
	"github.com/spf13/cobra"
	"os"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Displays one or many contexts from the kubeconfig file",
	Long:  `Displays one or many contexts from the kubeconfig file`,
	Example: `
  # List all the contexts in your kubeconfig file
  kubecm get

  # Describe one context in your kubeconfig file.
  kubecm get my-context
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := Formatable(nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			err := Formatable(args)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		err := ClusterStatus()
		if err != nil {
			fmt.Printf("Cluster check failure!\n%v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.SetArgs([]string{""})
}

func Formatable(args []string) error {
	config, err := LoadClientConfig(cfgFile)
	if err != nil {
		return err
	}
	var table [][]string
	if args == nil {
		for key, obj := range config.Contexts {
			var tmp []string
			if config.CurrentContext == key {
				tmp = append(tmp, "*")
			} else {
				tmp = append(tmp, "")
			}
			tmp = append(tmp, key)
			tmp = append(tmp, obj.Cluster)
			tmp = append(tmp, obj.AuthInfo)
			tmp = append(tmp, obj.Namespace)
			table = append(table, tmp)
		}
	} else {
		for key, obj := range config.Contexts {
			var tmp []string
			if config.CurrentContext == key {
				tmp = append(tmp, "*")
				tmp = append(tmp, key)
				tmp = append(tmp, obj.Cluster)
				tmp = append(tmp, obj.AuthInfo)
				tmp = append(tmp, obj.Namespace)
				table = append(table, tmp)
			}
		}
	}

	if table != nil {
		tabulate := gotabulate.Create(table)
		tabulate.SetHeaders([]string{"CURRENT", "NAME", "CLUSTER", "USER", "Namespace"})
		// Turn On String Wrapping
		tabulate.SetWrapStrings(true)
		// Render the table
		tabulate.SetAlign("center")
		fmt.Println(tabulate.Render("grid", "left"))
	} else {
		return fmt.Errorf("context %v not found", args)
	}
	return nil
}
