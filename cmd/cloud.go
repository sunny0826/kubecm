package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/sunny0826/kubecm/pkg/cloud"
)

// CloudCommand cloud command struct
type CloudCommand struct {
	BaseCommand
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
	{
		Name:     "AWS",
		Alias:    []string{"aws", "eks"},
		HomePage: "https://console.aws.amazon.com/eks/home",
		Service:  "EKS",
	},
}

// Init CloudCommand
func (cc *CloudCommand) Init() {
	cc.command = &cobra.Command{
		Use:   "cloud [COMMANDS]",
		Short: "Manage kubeconfig from cloud",
		Long:  "Manage kubeconfig from cloud",
	}
	cc.command.PersistentFlags().String("provider", "", "public cloud")
	cc.command.PersistentFlags().String("cluster_id", "", "kubernetes cluster id")
	cc.command.PersistentFlags().String("region_id", "", "cloud region id")
	cc.AddCommands(&CloudAddCommand{})
	cc.AddCommands(&CloudListCommand{})
}

func getClusters(provider, regionID string, num int) ([]cloud.ClusterInfo, error) {
	var clusters []cloud.ClusterInfo
	var err error
	switch num {
	case -1:
		var allAlias []string
		for _, cloudInfo := range Clouds {
			allAlias = append(allAlias, cloudInfo.Alias...)
		}
		return nil, fmt.Errorf("'%s' is not supported, supported cloud alias are %v", provider, allAlias)
	case 0:
		accessKeyID, accessKeySecret := checkEnvForSecret(0)
		ali := cloud.AliCloud{
			AccessKeyID:     accessKeyID,
			AccessKeySecret: accessKeySecret,
		}
		clusters, err = ali.ListCluster()
		if err != nil {
			return nil, err
		}
	case 1:
		secretID, secretKey := checkEnvForSecret(1)
		ten := cloud.TencentCloud{
			SecretID:  secretID,
			SecretKey: secretKey,
		}
		if regionID == "" {
			regionList, err := ten.GetRegionID()
			if err != nil {
				return nil, err
			}
			regionNum := selectRegion(regionList, "Select Region Name")
			ten.RegionID = regionList[regionNum]
		} else {
			ten.RegionID = regionID
		}
		clusters, err = ten.ListCluster()
		if err != nil {
			return nil, err
		}
	case 2:
		fmt.Println("⛅  Selected: Rancher")
		serverURL, apiKey := checkEnvForSecret(2)
		rancher := cloud.Rancher{
			ServerURL: serverURL,
			APIKey:    apiKey,
		}
		clusters, err = rancher.ListCluster()
		if err != nil {
			return nil, err
		}
	case 3:
		fmt.Println("⛅  Selected: AWS")
		accessKeyID, accessKeySecret := checkEnvForSecret(3)
		aws := cloud.AWS{
			AccessKeyID:     accessKeyID,
			AccessKeySecret: accessKeySecret,
		}
		if regionID == "" {
			regionList, err := cloud.GetRegionID()
			if err != nil {
				return nil, err
			}
			regionNum := selectRegion(regionList, "Select Region ID")
			aws.RegionID = regionList[regionNum]
		} else {
			aws.RegionID = regionID
		}
		clusters, err = aws.ListCluster()
		if err != nil {
			return nil, err
		}
	}
	return clusters, err
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
	case 3:
		accessKeyID, id := os.LookupEnv("AWS_ACCESS_KEY_ID")
		accessKeySecret, key := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
		if !id || !key {
			accessKeyID = PromptUI("AWS Access Key ID", "")
			accessKeySecret = PromptUI("AWS Access Key Secret", "")
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
