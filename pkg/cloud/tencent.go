package cloud

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

// TencentCloud struct of tencent cloud
type TencentCloud struct {
	SecretID  string
	SecretKey string
	RegionID  string
}

// getTenClient get tencent openapi client
func getTenClient(SecretID, SecretKey, RegionID string) (*tke.Client, error) {
	credential := common.NewCredential(
		SecretID,
		SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tke.tencentcloudapi.com"
	client, err := tke.NewClient(credential, RegionID, cpf)
	return client, err
}

// GetRegionID get region id of tke cluster
func (t *TencentCloud) GetRegionID() ([]string, error) {
	client, err := getTenClient(t.SecretID, t.SecretKey, "")
	if err != nil {
		return nil, err
	}
	request := tke.NewDescribeRegionsRequest()
	response, err := client.DescribeRegions(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	var regionList []string
	for _, region := range response.Response.RegionInstanceSet {
		regionList = append(regionList, *region.RegionName)
	}
	fmt.Printf("%s", regionList)
	return regionList, nil
}

// ListCluster list tke cluster info
func (t *TencentCloud) ListCluster() (clusters []ClusterInfo, err error) {
	client, err := getTenClient(t.SecretID, t.SecretKey, t.RegionID)
	if err != nil {
		return nil, err
	}
	request := tke.NewDescribeClustersRequest()
	response, err := client.DescribeClusters(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	var clusterList []ClusterInfo
	for _, cluster := range response.Response.Clusters {
		clusterList = append(clusterList, ClusterInfo{
			Name:       *cluster.ClusterName,
			ID:         *cluster.ClusterId,
			RegionID:   t.RegionID,
			K8sVersion: *cluster.ClusterVersion,
			ConsoleURL: fmt.Sprintf("https://console.cloud.tencent.com/tke2/cluster/sub/list/basic/info?clusterId=%s", *cluster.ClusterId),
		})
	}
	return clusterList, err
}

// GetKubeConfig get tke kubeConfig file
func (t *TencentCloud) GetKubeConfig(clusterID string) (string, error) {
	client, err := getTenClient(t.SecretID, t.SecretKey, t.RegionID)
	if err != nil {
		return "", err
	}
	request := tke.NewDescribeClusterKubeconfigRequest()
	request.ClusterId = common.StringPtr(clusterID)
	response, err := client.DescribeClusterKubeconfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return "", err
	}
	return *(response.Response.Kubeconfig), err
}
