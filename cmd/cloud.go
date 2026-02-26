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
	{
		Name:     "Azure",
		Alias:    []string{"azure", "aks"},
		HomePage: "https://portal.azure.com",
		Service:  "AKS",
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
	cc.command.PersistentFlags().String("aws_profile", "", "AWS profile name (from ~/.aws/config)")
	cc.AddCommands(&CloudAddCommand{})
	cc.AddCommands(&CloudListCommand{})
	cc.AddCommands(&DocsCommand{})
}

func getClusters(provider, regionID, awsProfile string, num int) ([]cloud.ClusterInfo, error) {
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
		awsProvider, err := buildAWSProvider(awsProfile, regionID)
		if err != nil {
			return nil, err
		}
		clusters, err = awsProvider.ListCluster()
		if err != nil {
			return nil, err
		}
	case 4:
		fmt.Println("⛅  Selected: Azure")
		authModes := []string{"Default (SDK Auth)", "Service Principal"}
		authMode := selectOption(nil, authModes, "Select Auth Type")
		azure := cloud.Azure{
			AuthMode:       cloud.AzureAuth(authMode),
			TenantID:       os.Getenv("AZURE_TENANT_ID"),
			SubscriptionID: os.Getenv("AZURE_SUBSCRIPTION_ID"),
		}

		if azure.TenantID == "" {
			azure.TenantID = PromptUI("Azure Tenant ID", "")
		}

		if azure.AuthMode == cloud.AuthModeServicePrincipal {
			azure.ClientID, azure.ClientSecret = checkEnvForSecret(4)
			azure.ObjectID = os.Getenv("AZURE_OBJECT_ID")

			if azure.ObjectID == "" {
				azure.ObjectID = PromptUI("Azure Object ID", "")
			}
		}

		subscriptionList, err := azure.ListSubscriptions()
		if err != nil {
			return nil, err
		}

		for _, subscription := range subscriptionList {
			if azure.SubscriptionID != "" && azure.SubscriptionID != subscription.ID {
				continue
			}

			subscriptionClusters, err := azure.ListCluster(subscription)
			if err != nil {
				return nil, err
			}
			clusters = append(clusters, subscriptionClusters...)
		}
	}

	return clusters, err
}

// buildAWSProvider creates an AWS provider with the appropriate auth mode.
// If awsProfile is set, it uses the default credential chain with that profile.
// Otherwise, it prompts the user to select an auth type.
func buildAWSProvider(awsProfile, regionID string) (cloud.AWS, error) {
	a := cloud.AWS{}

	if awsProfile != "" {
		// Profile flag provided: use default credential chain with profile
		a.AuthMode = cloud.AWSAuthDefault
		a.Profile = awsProfile
	} else {
		// Prompt for auth type (like Azure does)
		authModes := []string{"Default (Credential Chain)", "Static Credentials"}
		authMode := selectOption(nil, authModes, "Select Auth Type")
		switch authMode {
		case 0:
			a.AuthMode = cloud.AWSAuthDefault
			// Check for AWS_PROFILE env var
			if envProfile := os.Getenv("AWS_PROFILE"); envProfile != "" {
				a.Profile = envProfile
			}
		case 1:
			a.AuthMode = cloud.AWSAuthStaticCredentials
			a.AccessKeyID, a.AccessKeySecret = checkEnvForSecret(3)
		}
	}

	// Resolve region
	if regionID != "" {
		a.RegionID = regionID
	} else {
		// Check environment variables
		if envRegion := os.Getenv("AWS_REGION"); envRegion != "" {
			a.RegionID = envRegion
		} else if envRegion := os.Getenv("AWS_DEFAULT_REGION"); envRegion != "" {
			a.RegionID = envRegion
		} else {
			regionList, err := cloud.GetRegionID()
			if err != nil {
				return cloud.AWS{}, err
			}
			regionNum := selectRegion(regionList, "Select Region ID")
			a.RegionID = regionList[regionNum]
		}
	}

	return a, nil
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
	case 4:
		accessKeyID, id := os.LookupEnv("AZURE_CLIENT_ID")
		accessKeySecret, key := os.LookupEnv("AZURE_CLIENT_SECRET")
		if !id || !key {
			accessKeyID = PromptUI("Azure Client ID", "")
			accessKeySecret = PromptUI("Azure Client Secret", "")
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
{{- with .Account }}
{{ "Account:" | faint }}	{{ . }}
{{- end }}
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

	return selectOption(templates, regionList, label)
}

func selectOption(templates *promptui.SelectTemplates, options []string, label string) int {
	if templates == nil {
		templates = &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Inactive: "  {{ . | cyan }}",
			Selected: "\U0001F30D Selected: {{ . | green }}",
		}
	}

	searcher := func(input string, index int) bool {
		pepper := options[index]
		name := strings.ReplaceAll(strings.ToLower(pepper), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")
		return strings.Contains(name, input)
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     options,
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
