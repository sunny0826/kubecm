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
	"fmt"
	"github.com/bndr/gotabulate"
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
	Example: `
# List all the contexts in your kubeconfig file
kubecm
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := Formatable(nil)
			if err != nil {
				Error.Println(err)
				os.Exit(1)
			}
		} else {
			err := Formatable(args)
			if err != nil {
				Error.Println(err)
				os.Exit(1)
			}
		}
		err := ClusterStatus()
		if err != nil {
			log.Printf("Cluster check failure!\n%v", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
