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
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type pepper struct {
	Name     string
	HeatUnit int
	Peppers  int
}

type needle struct {
	Name    string
	Cluster string
	User    string
}

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch Kube Context interactively.",
	Example: `
# Switch Kube Context interactively
kubecm switch
`,
	Long: `
Switch Kube Context interactively.
`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadClientConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		var kubeItems []needle
		current := config.CurrentContext
		for key, obj := range config.Contexts {
			if key != current {
				kubeItems = append(kubeItems, needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
			} else {
				kubeItems = append([]needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo}}, kubeItems...)
			}
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\U0001F63C {{ .Name | red }}",
			Inactive: "  {{ .Name | cyan }}",
			Selected: "\U0001F638 {{ .Name | green }}",
			Details: `
--------- Info ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Cluster:" | faint }}	{{ .Cluster }}
{{ "User:" | faint }}	{{ .User }}`,
		}

		searcher := func(input string, index int) bool {
			pepper := kubeItems[index]
			name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, input)
		}

		prompt := promptui.Select{
			Label:     "Select Kube Context",
			Items:     kubeItems,
			Templates: templates,
			Size:      4,
			Searcher:  searcher,
		}

		i, _, err := prompt.Run()
		kubeName := kubeItems[i].Name
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		config.CurrentContext = kubeName

		err = ModifyKubeConfig(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Switched to context %s\n", kubeName)
		err = Formatable(nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
