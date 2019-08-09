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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"strings"
)

var file string
var name string
var cover bool

type (
	Config struct {
		ApiVersion     string     `yaml:"apiVersion"`
		Kind           string     `yaml:"kind"`
		Clusters       []Clusters `yaml:"clusters"`
		Contexts       []Contexts `yaml:"contexts"`
		CurrentContext string     `yaml:"current-context"`
		Users          []Users    `yaml:"users"`
	}
	Clusters struct {
		Cluster Cluster `yaml:"cluster"`
		Name    string  `yaml:"name"`
	}
	Cluster struct {
		Server                   string `yaml:"server"`
		CertificateAuthorityData string `yaml:"certificate-authority-data"`
	}
	Contexts struct {
		Context Context `yaml:"context"`
		Name    string  `yaml:"name"`
	}
	Context struct {
		Cluster   string `yaml:"cluster"`
		User      string `yaml:"user"`
		NameSpace string `yaml:"namespace,omitempty"`
	}
	Users struct {
		Name string `yaml:"name"`
		User User   `yaml:"user"`
	}
	User struct {
		ClientCertificateData string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
	}
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Merge configuration file with ./kube/config",
	Example: `
# Merge example.yaml with ./kube/config
kubecm add -f example.yaml 

# Merge example.yaml and name contexts test with ./kube/config
kubecm add -f example.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f example.yaml -c
`,
	Run: func(cmd *cobra.Command, args []string) {
		if FileExists(file) {
			err := ConfigCheck(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cover, _ = cmd.Flags().GetBool("cover")
			oldYaml := Config{}
			oldYaml.ReadYaml(cfgFile)
			addYaml := Config{}
			addYaml.ReadYaml(file)
			err = oldYaml.MergeConfig(addYaml)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("%s file does not exist", file)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&file, "file", "f", "", "Path to merge kubeconfig files")
	addCmd.Flags().StringVarP(&name, "name", "n", "", "The name of contexts. if this field is null,it will be named with file name.")
	addCmd.Flags().BoolP("cover", "c", false, "Overwrite the original kubeconfig file")
	addCmd.MarkFlagRequired("file")
}

func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (c *Config) ReadYaml(f string) {
	buffer, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = yaml.Unmarshal(buffer, &c)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (c *Config) WriteYaml() {
	buffer, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if cover {
		err = ioutil.WriteFile(cfgFile, buffer, 0777)
	} else {
		err = ioutil.WriteFile("./config.yaml", buffer, 0777)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetName() string {
	if name == "" {
		n := strings.Split(file, "/")
		result := strings.Split(n[len(n)-1], ".")
		return result[0]
	} else {
		return name
	}
}

func (c *Config) MergeConfig(a Config) error {
	name := GetName()
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
	c.WriteYaml()
	return nil
}

func (c *Config) Check(name string) error {
	for _, old := range c.Contexts {
		if old.Name == name {
			return fmt.Errorf("The name: %s already exists, please replace it.", name)
		}
	}
	return nil
}

func ConfigCheck(kubeconfigPath string) error {
	_, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}
	return nil
}
