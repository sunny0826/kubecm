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

var rootCmd = &cobra.Command{
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
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	Info = log.New(os.Stdout, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "Error:", log.Ldate|log.Ltime|log.Lshortfile)
}

func initConfig() {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	rootCmd.Flags().StringVar(&cfgFile, "config", *kubeconfig, "config.yaml file (default is $HOME/.kubecm.yaml)")
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