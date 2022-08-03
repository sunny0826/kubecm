package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sunny0826/kubecm/pkg/cloud"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// AddCloudCommand add command struct
type AddCloudCommand struct {
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
		Alias:    []string{"alibabacloud", "alicloud", "aliyun", "ack"},
		HomePage: "https://cs.console.aliyun.com",
		Service:  "ACK",
	},
	{
		Name:     "TencentCloud",
		Alias:    []string{"tencentcloud", "tencent", "tke"},
		HomePage: "https://console.cloud.tencent.com/tke",
		Service:  "TKE",
	},
	{
		Name:     "Rancher",
		Alias:    []string{"rancher"},
		HomePage: "https://rancher.com",
		Service:  "Rancher",
	},
}

// Init AddCommand
func (cc *AddCloudCommand) Init() {
	cc.command = &cobra.Command{
		Use:   "cloud",
		Short: "Add kubeconfig from public cloud",
		Long:  "Add kubeconfig from public cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.runCloud(cmd, args)
		},
		Example: addCloudExample(),
	}
	cc.command.Flags().String("provider", "", "public cloud")
	cc.command.Flags().String("cluster_id", "", "kubernetes cluster id")
	cc.command.Flags().String("region_id", "", "cloud region id")
}

func (cc *AddCloudCommand) runCloud(cmd *cobra.Command, args []string) error {
	provider, _ := cc.command.Flags().GetString("provider")
	clusterID, _ := cc.command.Flags().GetString("cluster_id")
	regionID, _ := cc.command.Flags().GetString("region_id")
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
			if len(clusters) == 0 {
				return errors.New("no clusters found")
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
	case 1:
		fmt.Println("⛅  Selected: TencentCloud")
		secretID, secretKey := checkEnvForSecret(1)
		ten := cloud.TencentCloud{
			SecretID:  secretID,
			SecretKey: secretKey,
		}
		if regionID == "" {
			regionList, err := ten.GetRegionID()
			if err != nil {
				return err
			}
			regionNum := selectRegion(regionList, "Select Region Name")
			ten.RegionID = regionList[regionNum]
		} else {
			ten.RegionID = regionID
		}

		if clusterID == "" {
			clusters, err := ten.ListCluster()
			if err != nil {
				return err
			}
			if len(clusters) == 0 {
				return errors.New("no clusters found")
			}
			clusterNum := selectCluster(clusters, "Select Cluster")
			kubeconfig, err := ten.GetKubeConfig(clusters[clusterNum].ID)
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
			kubeconfig, err := ten.GetKubeConfig(clusterID)
			if err != nil {
				return err
			}
			newConfig, err := clientcmd.Load([]byte(kubeconfig))
			if err != nil {
				return err
			}
			err = AddToLocal(newConfig, fmt.Sprintf("tencent-%s", clusterID), cover)
			if err != nil {
				return err
			}
		}
	case 2:
		fmt.Println("⛅  Selected: Rancher")
		serverURL, apiKey := checkEnvForSecret(2)
		rancher := cloud.Rancher{
			ServerURL: serverURL,
			APIKey:    apiKey,
		}
		if clusterID == "" {
			clusters, err := rancher.ListCluster()
			if err != nil {
				return err
			}
			if len(clusters) == 0 {
				return errors.New("no clusters found")
			}
			clusterNum := selectCluster(clusters, "Select Cluster")
			kubeconfig, err := rancher.GetKubeConfig(clusters[clusterNum].ID)
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
			kubeconfig, err := rancher.GetKubeConfig(clusterID)
			if err != nil {
				return err
			}
			newConfig, err := clientcmd.Load([]byte(kubeconfig))
			if err != nil {
				return err
			}
			err = AddToLocal(newConfig, fmt.Sprintf("rancher-%s", clusterID), cover)
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
			accessKeyID = PromptUI("AlibabaCloud Access Key ID", "")
			accessKeySecret = PromptUI("AlibabaCloud Access Key Secret", "")
		}
		return accessKeyID, accessKeySecret
	case 1:
		secretID, id := os.LookupEnv("TENCENTCLOUD_SECRET_ID")
		secretKey, key := os.LookupEnv("TENCENTCLOUD_SECRET_KEY")
		if !id || !key {
			secretID = PromptUI("TencentCloud API secretId", "")
			secretKey = PromptUI("TencentCloud API secretKey", "")
		}
		return secretID, secretKey
	case 2:
		serverURL, su := os.LookupEnv("RANCHER_SERVER_URL")
		apiKey, key := os.LookupEnv("RANCHER_API_KEY")
		if !su || !key {
			serverURL = PromptUI("Rancher API serverURL", "")
			apiKey = PromptUI("Rancher API key", "")
		}
		return serverURL, apiKey
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
		Size:      uiSize,
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
{{ "Version:" | faint }}	{{ .K8sVersion }}
{{ "ID:" | faint }}	{{ .ID }}`,
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     clouds,
		Templates: templates,
		Size:      uiSize,
	}
	i, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return i
}

func selectRegion(regionList []string, label string) int {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F680 {{ .Name | red }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F30D Selected: {{ . | green }}",
	}
	searcher := func(input string, index int) bool {
		pepper := regionList[index]
		name := strings.Replace(strings.ToLower(pepper), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     regionList,
		Templates: templates,
		Size:      uiSize,
		Searcher:  searcher,
	}
	i, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return i
}

func addCloudExample() string {
	return `
# Supports Ali Cloud and Tencent Cloud
# The AK/AS of the cloud platform will be retrieved directly 
# if it exists in the environment variable, 
# otherwise a prompt box will appear asking for it.

# Set env AliCloud secret key
export ACCESS_KEY_ID=xxx
export ACCESS_KEY_SECRET=xxx
# Set env Tencent secret key
export TENCENTCLOUD_SECRET_ID=xxx
export TENCENTCLOUD_SECRET_KEY=xxx
# Set env Rancher secret key
export RANCHER_SERVER_URL=https://xxx
export RANCHER_API_KEY=xxx
# Interaction: select kubeconfig from the cloud
kubecm add cloud
# Add kubeconfig from cloud
kubecm add cloud --provider alibabacloud --cluster_id=xxxxxx
`
}
