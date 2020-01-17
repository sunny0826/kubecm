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
	"log"
	"os"
)

type needle struct {
	Name    string
	Cluster string
	User    string
	Center  string
}

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch Kube Context interactively.",
	Long: `
Switch Kube Context interactively.
`,
	Example: switchExample(),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadClientConfig(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		var kubeItems []needle
		current := config.CurrentContext
		for key, obj := range config.Contexts {
			if key != current {
				kubeItems = append(kubeItems, needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
			} else {
				kubeItems = append([]needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
			}
		}
		num := SelectUI(kubeItems, "Select Kube Context")
		kubeName := kubeItems[num].Name
		config.CurrentContext = kubeName
		err = ModifyKubeConfig(config)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Printf("Switched to context 「%s」\n", config.CurrentContext)
		err = Formatable(nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
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
		log.Println("Unexpected error: %v", err)
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
	log.Printf("Cluster check succeeded!\nContains components: %v \n", names)
	return nil
}

func switchExample() string {
	return `
# Switch Kube Context interactively
kubecm switch
`
}
