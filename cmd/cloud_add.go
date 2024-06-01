package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mgutz/ansi"

	"github.com/spf13/cobra"
	"github.com/sunny0826/kubecm/pkg/cloud"
	"k8s.io/client-go/tools/clientcmd"
)

// CloudAddCommand add command struct
type CloudAddCommand struct {
	CloudCommand
}

// Init AddCommand
func (ca *CloudAddCommand) Init() {
	ca.command = &cobra.Command{
		Use:   "add",
		Short: "Add kubeconfig from cloud",
		Long:  "Add kubeconfig from cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ca.runCloudAdd(cmd, args)
		},
		Example: cloudAddExample(),
	}
}

func (ca *CloudAddCommand) runCloudAdd(cmd *cobra.Command, args []string) error {
	provider, _ := ca.command.Flags().GetString("provider")
	clusterID, _ := ca.command.Flags().GetString("cluster_id")
	regionID, _ := ca.command.Flags().GetString("region_id")
	cover, _ := ca.command.Flags().GetBool("cover")
	context, _ := ca.command.Flags().GetStringSlice("context")
	selectContext, _ := ca.command.Flags().GetBool("select-context")
	contextTemplate, _ := ca.command.Flags().GetStringSlice("context-template")
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
			err = AddToLocal(newConfig, clusters[clusterNum].Name, "", cover, selectContext, contextTemplate, context)
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
			err = AddToLocal(newConfig, fmt.Sprintf("alicloud-%s", clusterID), "", cover, selectContext, contextTemplate, context)
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
			err = AddToLocal(newConfig, clusters[clusterNum].Name, "", cover, selectContext, contextTemplate, context)
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
			err = AddToLocal(newConfig, fmt.Sprintf("tencent-%s", clusterID), "", cover, selectContext, contextTemplate, context)
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
			err = AddToLocal(newConfig, clusters[clusterNum].Name, "", cover, selectContext, contextTemplate, context)
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
			err = AddToLocal(newConfig, fmt.Sprintf("rancher-%s", clusterID), "", cover, selectContext, contextTemplate, context)
			if err != nil {
				return err
			}
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
				return err
			}
			regionNum := selectRegion(regionList, "Select Region ID")
			aws.RegionID = regionList[regionNum]
		} else {
			aws.RegionID = regionID
		}
		if clusterID == "" {
			clusters, err := aws.ListCluster()
			if err != nil {
				return err
			}
			if len(clusters) == 0 {
				return errors.New("no clusters found")
			}
			clusterNum := selectCluster(clusters, "Select Cluster")
			clusterID = clusters[clusterNum].ID
		}
		newConfig, err := aws.GetKubeConfigObj(clusterID)
		if err != nil {
			return err
		}
		err = AddToLocal(newConfig, fmt.Sprintf("aws-%s", clusterID), "", cover, selectContext, contextTemplate, context)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n",
			ansi.Color("Note", "blue"),
			ansi.Color(" please install the AWS CLI before normal use.", "white+h"))
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

		if clusterID != "" {
			clusterIDParts := strings.Split(clusterID, "/")
			if len(clusterIDParts) != 9 {
				return fmt.Errorf("invalid id %s", clusterID)
			}
			azure.SubscriptionID = clusterIDParts[2]
			resourceGroup := clusterIDParts[4]
			clusterName := clusterIDParts[8]

			var (
				kubeConfig []byte
				err        error
			)
			kubeConfigType := selectOption(nil, []string{"User Config", "Admin Config"}, "Select Config Type")
			switch kubeConfigType {
			case 0:
				kubeConfig, err = azure.GetKubeConfig(clusterName, resourceGroup)
			case 1:
				kubeConfig, err = azure.GetAdminKubeConfig(clusterName, resourceGroup)
			default:
				return fmt.Errorf("invalid config type %d", kubeConfigType)
			}

			if err != nil {
				return err
			}
			newConfig, err := clientcmd.Load(kubeConfig)
			if err != nil {
				return err
			}
			return AddToLocal(newConfig, fmt.Sprintf("azure-%s", clusterID), "", cover, selectContext, contextTemplate, context)
		}

		subscriptionList, err := azure.ListSubscriptions()
		if err != nil {
			return err
		}

		var clusters []cloud.ClusterInfo

		for _, subscription := range subscriptionList {
			if azure.SubscriptionID != "" && azure.SubscriptionID != subscription.ID {
				continue
			}

			subscriptionClusters, err := azure.ListCluster(subscription)
			if err != nil {
				return err
			}
			clusters = append(clusters, subscriptionClusters...)
		}

		if len(clusters) == 0 {
			return errors.New("no clusters found")
		}

		clusterNum := selectCluster(clusters, "Select Cluster")
		cluster := clusters[clusterNum]

		if azure.SubscriptionID == "" {
			azure.SubscriptionID = strings.Split(cluster.ID, "/")[2]
		}

		resourceGroup := strings.Split(cluster.ID, "/")[4]

		var kubeConfig []byte
		kubeConfigType := selectOption(nil, []string{"User Config", "Admin Config"}, "Select Config Type")
		switch kubeConfigType {
		case 0:
			kubeConfig, err = azure.GetKubeConfig(cluster.Name, resourceGroup)
		case 1:
			kubeConfig, err = azure.GetAdminKubeConfig(cluster.Name, resourceGroup)
		default:
			return fmt.Errorf("invalid config type %d", kubeConfigType)
		}

		if err != nil {
			return err
		}
		newConfig, err := clientcmd.Load(kubeConfig)
		if err != nil {
			return err
		}
		return AddToLocal(newConfig, fmt.Sprintf("azure-%s", clusterID), "", cover, selectContext, contextTemplate, context)

	}
	return nil
}

func cloudAddExample() string {
	return `
# Supports AWS, Ali Cloud, Tencent Cloud and Rancher
# The AK/AS of the cloud platform will be retrieved directly 
# if it exists in the environment variable, 
# otherwise a prompt box will appear asking for it.

# Set env AliCloud secret key
export ACCESS_KEY_ID=YOUR_AKID
export ACCESS_KEY_SECRET=YOUR_SECRET_KEY

# Set env Tencent secret key
export TENCENTCLOUD_SECRET_ID=YOUR_SECRET_ID
export TENCENTCLOUD_SECRET_KEY=YOUR_SECRET_KEY

# Set env Rancher secret key
export RANCHER_SERVER_URL=https://xxx
export RANCHER_API_KEY=YOUR_API_KEY

# Set env AWS secret key
# Note: Please install the AWS CLI before normal use.
export AWS_ACCESS_KEY_ID=YOUR_AKID
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY

# Set env Azure secret key
export AZURE_SUBSCRIPTION_ID=YOUR_SUBSCRIPTION_ID
export AZURE_CLIENT_ID=YOUR_CLIENT_ID
export AZURE_CLIENT_SECRET=YOUR_CLIENT_SECRET
export AZURE_TENANT_ID=YOUR_TENANT_ID
export AZURE_OBJECT_ID=YOUR_OBJECT_ID

# Interaction: select kubeconfig from the cloud
kubecm cloud add
# Add kubeconfig from cloud
kubecm cloud add --provider alibabacloud --cluster_id=xxxxxx
`
}
