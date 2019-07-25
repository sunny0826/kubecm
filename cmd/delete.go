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
			kubeYaml := Config{}
			kubeYaml.ReadYaml(cfgFile)
			err := kubeYaml.DeleteContext(args)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
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

func (c *Config) DeleteContext(ctxs []string) error {
	for i, ct := range c.Contexts {
		for _, ctx := range ctxs {
			if ct.Name == ctx {
				fmt.Println(fmt.Sprintf("delete: %s", ctx))
				c.Contexts = append(c.Contexts[:i], c.Contexts[i+1:]...)
				user := ct.Context.User
				cluster := ct.Context.Cluster
				for j, us := range c.Users {
					if us.Name == user {
						c.Users = append(c.Users[:j], c.Users[j+1:]...)
					}
				}
				for k, clu := range c.Clusters {
					if clu.Name == cluster {
						c.Clusters = append(c.Clusters[:k], c.Clusters[k+1:]...)
					}
				}
			}
		}
	}
	cover = true
	c.WriteYaml()
	return nil
}
