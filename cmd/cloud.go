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
	Alias    []string
	HomePage string
	Service  string
}

// Clouds date of clouds
var Clouds = []CloudInfo{
	{
		Name:     "AlibabaCloud",
		Alias:    []string{"alibabacloud", "alicloud", "aliyun"},
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
	cc.command.Flags().String("provider", "", "public cloud")
	cc.command.Flags().String("cluster_id", "", "kubernetes cluster id")
}

func (cc *CloudCommand) runCloud(cmd *cobra.Command, args []string) error {
	provider, _ := cc.command.Flags().GetString("provider")
	clusterID, _ := cc.command.Flags().GetString("cluster_id")
	var num int
	if provider == "" {
		num = selectCloud(Clouds, "Select Cloud")
	} else {
		num = cc.checkFlags(provider)
		fmt.Printf("⛅  Selected: %s\n", Clouds[num].Name)
	}
	switch num {
	case -1:
		var allAlias []string
		for _, cloud := range Clouds {
			for _, alias := range cloud.Alias {
				allAlias = append(allAlias, alias)
			}
		}
		fmt.Printf("'%s' is not supported, supported cloud alias are %v \n", provider, allAlias)
		return nil
	case 0:
		if clusterID == "" {
			clusters, err := aliyun.ListCluster()
			if err != nil {
				return err
			}
			clusterNum := selectCluster(clusters, "Select Cluster")
			kubeconfig, err := aliyun.GetKubeConfig(clusters[clusterNum].ID)
			if err != nil {
				return err
			}
			fmt.Println(kubeconfig)
		} else {
			kubeconfig, err := aliyun.GetKubeConfig(clusterID)
			if err != nil {
				return err
			}
			fmt.Println(kubeconfig)
		}
	}
	return nil
}

func (cc *CloudCommand) checkFlags(provider string) int {
	for i, cloud := range Clouds {
		for _, alias := range cloud.Alias {
			if alias == provider {
				return i
			}
		}
	}
	return -1
}

func selectCloud(clouds []CloudInfo, label string) int {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F680 {{ .Name | red }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U000026C5 Selected: {{ .Name | green }}",
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
		Selected: "\U0001F6A2 Selected: {{ .Name | green }}",
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
# Add kubeconfig from cloud
kubecm add cloud --provider alibabacloud --cluster_id=xxxxxx
`
}
