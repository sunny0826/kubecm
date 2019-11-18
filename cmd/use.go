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
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
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
			config, err := LoadClientConfig(cfgFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			var currentContext bool
			for key, _ := range config.Contexts {
				if key == context {
					config.CurrentContext = key
					currentContext = true
					fmt.Println(fmt.Sprintf("Switched to context %s", context))
				}
			}
			if !currentContext {
				fmt.Println(fmt.Sprintf("no context exists with the name: %s", context))
				os.Exit(1)
			} else {
				err := ModifyKubeConfig(config)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				err = Formatable(nil)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
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

func ModifyKubeConfig(config *clientcmdapi.Config) error {
	commandLineFile, _ := ioutil.TempFile("", "")
	defer os.Remove(commandLineFile.Name())
	configType := clientcmdapi.Config{
		AuthInfos: config.AuthInfos,
		Clusters:  config.Clusters,
		Contexts:  config.Contexts,
	}
	_ = clientcmd.WriteToFile(configType, commandLineFile.Name())
	pathOptions := clientcmd.NewDefaultPathOptions()

	if err := clientcmd.ModifyConfig(pathOptions, *config, true); err != nil {
		fmt.Println("Unexpected error: %v", err)
		return err
	}
	return nil
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
