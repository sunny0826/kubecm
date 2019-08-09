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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Sets the current-context in a kubeconfig file",
	Example: `
# Use the context for the test cluster
kubecm use test
`,
	Long: `
Sets the current-context in a kubeconfig file
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			context := args[0]
			kubeYaml := Config{}
			kubeYaml.ReadYaml(cfgFile)
			if kubeYaml.CheckContext(context) {
				tmpContext := kubeYaml.CurrentContext
				kubeYaml.CurrentContext = context
				cover = true
				kubeYaml.WriteYaml()
				err := ClusterStatus()
				if err != nil {
					fmt.Printf("Cluster check failure! Please check your kubeconfig.\n%v", err)
					kubeYaml.CurrentContext = tmpContext
					kubeYaml.WriteYaml()
					os.Exit(1)
				} else {
					fmt.Println(fmt.Sprintf("Switched to context %s", context))
					err := Formatable(nil)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			} else {
				fmt.Println(fmt.Sprintf("no context exists with the name: %s", context))
				os.Exit(1)
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

func (c *Config) CheckContext(name string) bool {
	for _, con := range c.Contexts {
		if con.Name == name {
			return true
		}
	}
	return false
}

func ClusterStatus() error {
	config, err := clientcmd.BuildConfigFromFlags("", cfgFile)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	cus, err := clientset.CoreV1().ComponentStatuses().List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	var names []string
	for _, k := range cus.Items {
		names = append(names, k.Name)
	}
	fmt.Printf("Cluster check succeeded!\nContains components: %v \n", names)
	return nil
}
