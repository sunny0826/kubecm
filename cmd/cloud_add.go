package cmd

import (
	"errors"
	"fmt"

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
			err = AddToLocal(newConfig, clusters[clusterNum].Name, "", cover)
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
			err = AddToLocal(newConfig, fmt.Sprintf("alicloud-%s", clusterID), "", cover)
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
			err = AddToLocal(newConfig, clusters[clusterNum].Name, "", cover)
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
			err = AddToLocal(newConfig, fmt.Sprintf("tencent-%s", clusterID), "", cover)
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
			err = AddToLocal(newConfig, clusters[clusterNum].Name, "", cover)
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
			err = AddToLocal(newConfig, fmt.Sprintf("rancher-%s", clusterID), "", cover)
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
		return AddToLocal(newConfig, fmt.Sprintf("aws-%s", clusterID), "", cover)
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
export AWS_ACCESS_KEY_ID=YOUR_AKID
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY

# Interaction: select kubeconfig from the cloud
kubecm cloud add
# Add kubeconfig from cloud
kubecm cloud add --provider alibabacloud --cluster_id=xxxxxx
`
}
