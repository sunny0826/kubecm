package cloud

import (
	ack "github.com/alibabacloud-go/cs-20151215/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// AliCloud struct of alibaba cloud
type AliCloud struct {
	AccessKeyID     string
	AccessKeySecret string
}

// getClient get aliyun openapi client
func getClient(accessKeyID, accessKeySecret string) (*ack.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     &accessKeyID,
		AccessKeySecret: &accessKeySecret,
		RegionId:        tea.String("cn-hongkong"),
	}
	result, err := ack.NewClient(config)
	return result, err
}

// GetRegionID get region id of ack cluster
func (a *AliCloud) GetRegionID() ([]string, error) {
	// ack No RegionID required, return nil
	return nil, nil
}

// ListCluster list cluster info
func (a *AliCloud) ListCluster() (clusters []ClusterInfo, err error) {
	client, err := getClient(a.AccessKeyID, a.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	describeClustersV1Request := &ack.DescribeClustersV1Request{}
	v1, err := client.DescribeClustersV1(describeClustersV1Request)
	if err != nil {
		return nil, err
	}
	var clusterList []ClusterInfo
	for _, info := range v1.Body.Clusters {
		clusterList = append(clusterList, ClusterInfo{
			Name:       *info.Name,
			ID:         *info.ClusterId,
			RegionID:   *info.RegionId,
			K8sVersion: *info.CurrentVersion,
		})
	}
	return clusterList, err
}

// GetKubeConfig get kubeConfig file
func (a *AliCloud) GetKubeConfig(clusterID string) (string, error) {
	client, err := getClient(a.AccessKeyID, a.AccessKeySecret)
	if err != nil {
		return "", err
	}
	describeClusterUserKubeconfigRequest := &ack.DescribeClusterUserKubeconfigRequest{}
	res, err := client.DescribeClusterUserKubeconfig(tea.String(clusterID), describeClusterUserKubeconfigRequest)
	if err != nil {
		return "", err
	}
	return *(res.Body.Config), err
}
