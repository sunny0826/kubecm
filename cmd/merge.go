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
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var folder string

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge the kubeconfig files in the specified directory.",
	Long:  `Merge the kubeconfig files in the specified directory.`,
	Example: `
# Merge kubeconfig in the test directory
kubecm merge -f test 

# Merge kubeconfig in the test directory and overwrite the original kubeconfig file
kubecm merge -f test -c
`,
	Run: func(cmd *cobra.Command, args []string) {
		cover, _ = cmd.Flags().GetBool("cover")
		files := listFile(folder)
		fmt.Printf("Loading kubeconfig file: %v \n", files)
		mergeYaml := Config{}
		for _, yaml := range files {
			err := ConfigCheck(yaml)
			if err != nil {
				fmt.Printf("Please check kubeconfig file: %v \n%s", yaml, err)
				os.Exit(1)
			}
			tmpYaml := Config{}
			tmpYaml.ReadYaml(yaml)
			err = mergeYaml.MergeAllConfig(tmpYaml, yaml)
			if err != nil {
				fmt.Println(err)
			}
		}
		if cover {
			fmt.Println("Overwrite the origin kubeconfig file.")
		}
		mergeYaml.ApiVersion = "v1"
		mergeYaml.Kind = "Config"
		mergeYaml.CurrentContext = mergeYaml.Contexts[0].Name
		mergeYaml.WriteYaml()
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().StringVarP(&folder, "folder", "f", "", "Kubeconfig folder")
	mergeCmd.Flags().BoolP("cover", "c", false, "Overwrite the original kubeconfig file")
	mergeCmd.MarkFlagRequired("folder")
}

func listFile(folder string) []string {
	files, _ := ioutil.ReadDir(folder)
	var flist []string
	for _, file := range files {
		if file.IsDir() {
			listFile(folder + "/" + file.Name())
		} else {
			flist = append(flist, fmt.Sprintf("%s/%s", folder, file.Name()))
			//fmt.Println(folder + "/" + file.Name())
		}
	}
	return flist
}

func (c *Config) MergeAllConfig(a Config, n string) error {
	name := NameHandle(n)
	err := c.Check(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	suffix := HashSuffix(a)
	for _, obj := range a.Clusters {
		obj.Name = fmt.Sprintf("cluster-%v", suffix)
		c.Clusters = append(c.Clusters, obj)
	}
	for _, obj := range a.Contexts {
		obj.Name = fmt.Sprintf("%s", name)
		obj.Context.Cluster = fmt.Sprintf("cluster-%v", suffix)
		obj.Context.User = fmt.Sprintf("user-%v", suffix)
		c.Contexts = append(c.Contexts, obj)
	}
	for _, obj := range a.Users {
		obj.Name = fmt.Sprintf("user-%v", suffix)
		c.Users = append(c.Users, obj)
	}
	return nil
}

func NameHandle(path string) string {
	n := strings.Split(path, "/")
	result := strings.Split(n[len(n)-1], ".")
	return result[0]
}
