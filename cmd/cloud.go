package cmd

import (
	"fmt"
	"log"

	"github.com/sunny0826/kubecm/pkg/cloud/aliyun"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// CloudCommand add command struct
type CloudCommand struct {
	AddCommand
}

// CloudInfo Public cloud info
type CloudInfo struct {
	Name     string
	HomePage string
	Service  string
}

// Clouds date of clouds
var Clouds = []CloudInfo{
	{
		Name:     "Aliyun",
		HomePage: "https://cs.console.aliyun.com",
		Service:  "ACK",
	},
}

// Init AddCommand
func (cc *CloudCommand) Init() {
	cc.command = &cobra.Command{
		Use:   "cloud",
		Short: "Manage kubeconfig with public cloud",
		Long:  "Manage kubeconfig with public cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.runCloud(cmd, args)
		},
		Example: addCloudExample(),
	}
}

func (cc *CloudCommand) runCloud(cmd *cobra.Command, args []string) error {
	num := selectCloud(Clouds, "Select Cloud")
	fmt.Println(Clouds[num].Name)
	clusters, err := aliyun.ListCluster()
	if err != nil {
		return err
	}
	clusterNum := selectCluster(clusters, "Select Cluster")
	fmt.Println(clusters[clusterNum].Name)
	return nil
}

func selectCloud(clouds []CloudInfo, label string) int {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F680 {{ .Name | red }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U000026C5 Select: {{ .Name | green }}",
		Details: `
--------- Info ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "HomePage:" | faint }}	{{ .HomePage }}
{{ "Service:" | faint }}	{{ .Service }}`,
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     clouds,
		Templates: templates,
		Size:      4,
	}
	i, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return i
}

func selectCluster(clouds []aliyun.ClusterInfo, label string) int {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F680 {{ .Name | red }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U000026C5 Select: {{ .Name | green }}",
		Details: `
--------- Info ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "RegionID:" | faint }}	{{ .RegionID }}
{{ "ID:" | faint }}	{{ .ID }}`,
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     clouds,
		Templates: templates,
		Size:      4,
	}
	i, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return i
}

func addCloudExample() string {
	return `
# Select kubeconfig from the cloud
kubecm add cloud
`
}
