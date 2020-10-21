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
)

type SwitchCommand struct {
	baseCommand
}

func (sc *SwitchCommand) Init() {
	sc.command = &cobra.Command{
		Use:   "switch",
		Short: "Switch Kube Context interactively",
		Long: `
Switch Kube Context interactively
`,
		Aliases: []string{"s"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return sc.runSwitch(cmd, args)
		},
		Example: switchExample(),
	}
}

func (sc *SwitchCommand) runSwitch(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
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
	// exit option
	kubeItems, err = ExitOption(kubeItems)
	if err != nil {
		return err
	}
	num := SelectUI(kubeItems, "Select Kube Context")
	kubeName := kubeItems[num].Name
	config.CurrentContext = kubeName
	err = sc.WriteConfig(true,cfgFile,config)
	if err != nil {
		return err
	}
	sc.command.Printf("Switched to context 「%s」\n", config.CurrentContext)
	err = Formatable()
	if err != nil {
		return err
	}
	return nil
}

func switchExample() string {
	return `
# Switch Kube Context interactively
kubecm switch
`
}
