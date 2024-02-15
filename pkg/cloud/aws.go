package cloud

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/aws/aws-sdk-go/aws/session"
)

// AWS struct of aws cloud
type AWS struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
}

// getSession get session of aws cloud
func (a *AWS) getSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(a.RegionID),
		Credentials: credentials.NewStaticCredentials(a.AccessKeyID, a.AccessKeySecret, ""),
	})
	return sess, err
}

// GetRegionID get region id of aws
func GetRegionID() ([]string, error) {
	partitions := endpoints.DefaultPartitions()
	var regionList []string
	for _, p := range partitions {
		for id := range p.Regions() {
			// fmt.Printf("%s\n", id)
			regionList = append(regionList, id)
		}
	}
	return regionList, nil
}

// ListCluster list cluster info of aws
func (a *AWS) ListCluster() (clusters []ClusterInfo, err error) {
	sess, err := a.getSession()
	if err != nil {
		return nil, err
	}
	var clusterList []ClusterInfo
	svc := eks.New(sess)
	input := &eks.ListClustersInput{}

	result, err := svc.ListClusters(input)
	if err != nil {
		return nil, err
	}

	for _, clusterName := range result.Clusters {
		clusterInfo, err := a.getClusterInfo(*clusterName)
		if err != nil {
			return nil, err
		}
		clusterList = append(clusterList, clusterInfo)
	}
	return clusterList, nil
}

// GetClusterInfo get cluster info of aws eks
func (a *AWS) getClusterInfo(clusterName string) (clusterInfo ClusterInfo, err error) {
	sess, err := a.getSession()
	if err != nil {
		return ClusterInfo{}, err
	}

	svcSts := sts.New(sess)
	callerIdentity, err := svcSts.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return ClusterInfo{}, err
	}

	svc := eks.New(sess)
	input := &eks.DescribeClusterInput{
		Name: &clusterName,
	}
	cluster, err := svc.DescribeCluster(input)
	if err != nil {
		return ClusterInfo{}, err
	}
	return ClusterInfo{
		ID:         *cluster.Cluster.Name,
		Account:    *callerIdentity.Account,
		Name:       *cluster.Cluster.Name,
		RegionID:   a.RegionID,
		K8sVersion: *cluster.Cluster.Version,
		ConsoleURL: fmt.Sprintf("https://%s.console.aws.amazon.com/eks/home?region=%s#/clusters/%s", a.RegionID, a.RegionID, *cluster.Cluster.Name),
	}, err

}

// GetKubeConfigObj get aws eks kubeConfig file
func (a *AWS) GetKubeConfigObj(clusterID string) (*clientcmdapi.Config, error) {
	sess, err := a.getSession()
	if err != nil {
		return nil, err
	}

	svc := eks.New(sess)
	input := &eks.DescribeClusterInput{
		Name: &clusterID,
	}
	cluster, err := svc.DescribeCluster(input)
	if err != nil {
		return nil, err
	}

	decodePem, err := base64.StdEncoding.DecodeString(*cluster.Cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, err
	}

	kubeconfig := &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			*cluster.Cluster.Name: {
				Server:                   *cluster.Cluster.Endpoint,
				CertificateAuthorityData: decodePem,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			*cluster.Cluster.Name: {
				Exec: &clientcmdapi.ExecConfig{
					APIVersion: "client.authentication.k8s.io/v1beta1",
					Command:    "aws",
					Args: []string{
						"eks",
						"get-token",
						"--cluster-name",
						*cluster.Cluster.Name,
						"--region",
						a.RegionID,
						"--output",
						"json",
					},
				},
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			*cluster.Cluster.Name: {
				Cluster:  *cluster.Cluster.Name,
				AuthInfo: *cluster.Cluster.Name,
			},
		},
		CurrentContext: *cluster.Cluster.Name,
	}

	return kubeconfig, nil
}
