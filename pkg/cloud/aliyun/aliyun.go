package aliyun

import (
	"fmt"

	cs20151215 "github.com/alibabacloud-go/cs-20151215/v2/client"
	env "github.com/alibabacloud-go/darabonba-env/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// ClusterInfo ack cluster info
type ClusterInfo struct {
	Name     string
	ID       string
	RegionID string
}

// GetClient get aliyun openapi client
func GetClient() (_result *cs20151215.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     env.GetEnv(tea.String("ACCESS_KEY_ID")),
		AccessKeySecret: env.GetEnv(tea.String("ACCESS_KEY_SECRET")),
		RegionId:        tea.String("cn-hongkong"),
	}
	_result = &cs20151215.Client{}
	_result, _err = cs20151215.NewClient(config)
	return _result, _err
}

// ListCluster list cluster info
func ListCluster() (clusters []ClusterInfo, _err error) {
	client, _err := GetClient()
	if _err != nil {
		return nil, _err
	}
	describeClustersV1Request := &cs20151215.DescribeClustersV1Request{}
	v1, _err := client.DescribeClustersV1(describeClustersV1Request)
	if _err != nil {
		return nil, _err
	}
	var clusterList []ClusterInfo
	for _, info := range v1.Body.Clusters {
		clusterList = append(clusterList, ClusterInfo{
			Name:     *info.Name,
			ID:       *info.ClusterId,
			RegionID: *info.RegionId,
		})
	}
	fmt.Println(clusterList)
	return clusterList, _err
}
