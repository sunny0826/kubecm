package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/sunny0826/kubecm/pkg/cloud"

	"k8s.io/client-go/tools/clientcmd"

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
	cover, _ := cc.command.Flags().GetBool("cover")
	var num int
	if provider == "" {
		num = selectCloud(Clouds, "Select Cloud")
	} else {
		num = checkFlags(provider)
	}
	switch num {
	case -1:
		var allAlias []string
		for _, cloud := range Clouds {
			allAlias = append(allAlias, cloud.Alias...)
		}
		fmt.Printf("'%s' is not supported, supported cloud alias are %v \n", provider, allAlias)
		return nil
	case 0:
		fmt.Println("⛅  Selected: AlibabaCloud")
		accessKeyID, accessKeySecret := checkEnvForSecret(0)
		ali := cloud.AliCloud{
			AccessKeyID:     accessKeyID,
			AccessKeySecret: accessKeySecret,
		}
		if clusterID == "" {
			clusters, err := ali.ListCluster()
			if err != nil {
				return err
			}
			clusterNum := selectCluster(clusters, "Select Cluster")
			kubeconfig, err := ali.GetKubeConfig(clusters[clusterNum].ID)
			if err != nil {
				return err
			}
			newConfig, err := clientcmd.Load([]byte(kubeconfig))
			if err != nil {
				return err
			}
			err = AddToLocal(newConfig, clusters[clusterNum].Name, cover)
			if err != nil {
				return err
			}
		} else {
			kubeconfig, err := ali.GetKubeConfig(clusterID)
			if err != nil {
				return err
			}
			newConfig, err := clientcmd.Load([]byte(kubeconfig))
			if err != nil {
				return err
			}
			err = AddToLocal(newConfig, fmt.Sprintf("alicloud-%s", clusterID), cover)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkEnvForSecret(num int) (string, string) {
	switch num {
	case 0:
		accessKeyID, id := os.LookupEnv("ACCESS_KEY_ID")
		accessKeySecret, sec := os.LookupEnv("ACCESS_KEY_SECRET")
		if !id || !sec {
			accessKeyID = PromptUI("access key id", "")
			accessKeySecret = PromptUI("access key secret", "")
		}
		return accessKeyID, accessKeySecret
	}
	return "", ""
}

func checkFlags(provider string) int {
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

func selectCluster(clouds []cloud.ClusterInfo, label string) int {
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
# Set env secret key
export ACCESS_KEY_ID=xxx
export ACCESS_KEY_SECRET=xxx
# Interaction: select kubeconfig from the cloud
kubecm add cloud
# Add kubeconfig from cloud
kubecm add cloud --provider alibabacloud --cluster_id=xxxxxx
`
}
