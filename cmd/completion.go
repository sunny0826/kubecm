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
	zsh "github.com/rsteube/cobra-zsh-gen"
	"github.com/spf13/cobra"
	"os"
)

type CompletionCommand struct {
	baseCommand
}

func (cc *CompletionCommand) Init() {
	cc.command = &cobra.Command{
		Use:     "completion",
		Short:   "Generates bash/zsh completion scripts",
		Long:    `Output shell completion code for the specified shell (bash or zsh).`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.runCompletion(cmd, args)
		},
		Example: completionExample(),
	}
}

func (cc *CompletionCommand) runCompletion(command *cobra.Command, args []string) error {
	complet := args[0]
	switch complet {
	case "bash":
		err := command.Root().GenBashCompletion(os.Stdout)
		if err != nil {
			return err
		}
	case "zsh":
		err := zsh.Wrap(cc.command).GenZshCompletion(os.Stdout)
		if err != nil {
			return err
		}
	default:
		cc.command.PrintErrln("Parameter error! Please input bash or zsh")
	}
	return nil
}

func completionExample() string {
	return `
# bash
kubecm completion bash > ~/.kube/kubecm.bash.inc
printf "
# kubecm shell completion
source '$HOME/.kube/kubecm.bash.inc'
" >> $HOME/.bash_profile
source $HOME/.bash_profile

# add to $HOME/.zshrc
source <(kubecm completion zsh)
# or
kubecm completion zsh > "${fpath[1]}/_kubecm"
`
}
