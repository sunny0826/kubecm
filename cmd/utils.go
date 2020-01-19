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
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/bndr/gotabulate"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	"log"
	"os"
	"strings"
)

// ModifyKubeConfig modify kubeconfig
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

// Copied from https://github.com/kubernetes/kubernetes
// /blob/master/pkg/kubectl/util/hash/hash.go
func hEncode(hex string) (string, error) {
	if len(hex) < 10 {
		return "", fmt.Errorf(
			"input length must be at least 10")
	}
	enc := []rune(hex[:10])
	for i := range enc {
		switch enc[i] {
		case '0':
			enc[i] = 'g'
		case '1':
			enc[i] = 'h'
		case '3':
			enc[i] = 'k'
		case 'a':
			enc[i] = 'm'
		case 'e':
			enc[i] = 't'
		}
	}
	return string(enc), nil
}

// Hash returns the hex form of the sha256 of the argument.
func Hash(data string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}

// HashSuffix return the string of kubeconfig.
func HashSuf(config *clientcmdapi.Config) string {
	re_json, err := runtime.Encode(clientcmdlatest.Codec, config)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
	sum, _ := hEncode(Hash(string(re_json)))
	return sum
}

// Formatable generate table
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

// SelectUI output select ui
func SelectUI(kubeItems []needle, label string) int {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F63C {{ .Name | red }}{{ .Center | red}}",
		Inactive: "  {{ .Name | cyan }}{{ .Center | red}}",
		Selected: "\U0001F638 Select:{{ .Name | green }}",
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
		Label:     label,
		Items:     kubeItems,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}
	i, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return i
}

// PromptUI output prompt ui
func PromptUI(label string, name string) string {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Context name must have more than 3 characters")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  name,
	}
	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// BoolUI output bool ui
func BoolUI(label string) string {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F37A {{ . | red }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F47B {{ . | green }}",
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     []string{"True", "False"},
		Templates: templates,
		Size:      2,
	}
	_, obj, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return obj
}

// ClusterStatus output cluster status
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

// WriteConfig write kubeconfig
func WriteConfig(config []byte) error {
	if cover {
		err := ioutil.WriteFile(cfgFile, config, 0777)
		if err != nil {
			return err
		}
	} else {
		err := ioutil.WriteFile("./config.yaml", config, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}