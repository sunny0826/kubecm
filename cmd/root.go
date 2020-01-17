/*
Copyright © 2019 Guo Xudong

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
	"flag"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var cfgFile string

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type Cli struct {
	rootCmd *cobra.Command
}

//NewCli returns the cli instance used to register and execute command
func NewCli() *Cli {
	cli := &Cli{
		rootCmd: &cobra.Command{
			Use:   "kubecm",
			Short: "KubeConfig Manager.",
			Long: `
KubeConfig Manager
 _          _
| | ___   _| |__   ___  ___ _ __ ___
| |/ / | | | '_ \ / _ \/ __| '_ \ _ \
|   <| |_| | |_) |  __/ (__| | | | | |
|_|\_\\__,_|_.__/ \___|\___|_| |_| |_|

Find more information at: https://github.com/sunny0826/kubecm
`,
			Example: cliExample(),
			Run: func(cmd *cobra.Command, args []string) {
				runCli(cmd, args)
			},
		},
	}
	cli.rootCmd.SetOutput(os.Stdout)
	cli.setFlags()
	return cli
}

func (cli *Cli) setFlags() {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flags := cli.rootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", *kubeconfig, "path of kubeconfig")
}

//Run command
func (cli *Cli) Run() error {
	return cli.rootCmd.Execute()
}

func runCli(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		err := Formatable(nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := Formatable(args)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := ClusterStatus()
	if err != nil {
		log.Fatalf("Cluster check failure!\n%v", err)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func cliExample() string {
	return `
# List all the contexts in your kubeconfig file
kubecm
`
}
