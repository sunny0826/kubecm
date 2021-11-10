package cloud

import (
	cs20151215 "github.com/alibabacloud-go/cs-20151215/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// AliCloud struct of alibaba cloud
type AliCloud struct {
	AccessKeyID     string
	AccessKeySecret string
}

// GetClient get aliyun openapi client
func getClient(accessKeyID, accessKeySecret string) (*cs20151215.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     &accessKeyID,
		AccessKeySecret: &accessKeySecret,
		RegionId:        tea.String("cn-hongkong"),
	}
	result, err := cs20151215.NewClient(config)
	return result, err
}

// ListCluster list cluster info
func (a *AliCloud) ListCluster() (clusters []ClusterInfo, err error) {
	client, err := getClient(a.AccessKeyID, a.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	describeClustersV1Request := &cs20151215.DescribeClustersV1Request{}
	v1, err := client.DescribeClustersV1(describeClustersV1Request)
	if err != nil {
		return nil, err
	}
	var clusterList []ClusterInfo
	for _, info := range v1.Body.Clusters {
		clusterList = append(clusterList, ClusterInfo{
			Name:     *info.Name,
			ID:       *info.ClusterId,
			RegionID: *info.RegionId,
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
	describeClusterUserKubeconfigRequest := &cs20151215.DescribeClusterUserKubeconfigRequest{}
	res, err := client.DescribeClusterUserKubeconfig(tea.String(clusterID), describeClusterUserKubeconfigRequest)
	if err != nil {
		return "", err
	}
	return *(res.Body.Config), err
}
