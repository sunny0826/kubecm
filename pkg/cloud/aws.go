package cloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// AWSAuth authentication mode for AWS
type AWSAuth int

const (
	// AWSAuthDefault uses the default credential chain (profile, SSO, env, IMDS)
	AWSAuthDefault AWSAuth = iota
	// AWSAuthStaticCredentials uses static access key + secret
	AWSAuthStaticCredentials
)

// AWS struct of aws cloud
type AWS struct {
	AuthMode        AWSAuth
	Profile         string // for AWSAuthDefault
	AccessKeyID     string // for AWSAuthStaticCredentials
	AccessKeySecret string // for AWSAuthStaticCredentials
	RegionID        string

	cfg *aws.Config // cached config
}

// getAWSConfig returns an AWS config, caching after first call.
func (a *AWS) getAWSConfig() (aws.Config, error) {
	if a.cfg != nil {
		return *a.cfg, nil
	}

	var (
		cfg aws.Config
		err error
	)

	ctx := context.Background()
	opts := []func(*config.LoadOptions) error{}

	if a.RegionID != "" {
		opts = append(opts, config.WithRegion(a.RegionID))
	}

	switch a.AuthMode {
	case AWSAuthDefault:
		if a.Profile != "" {
			opts = append(opts, config.WithSharedConfigProfile(a.Profile))
		}
		cfg, err = config.LoadDefaultConfig(ctx, opts...)
	case AWSAuthStaticCredentials:
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(a.AccessKeyID, a.AccessKeySecret, ""),
		))
		cfg, err = config.LoadDefaultConfig(ctx, opts...)
	default:
		return aws.Config{}, fmt.Errorf("invalid AWS auth mode: %d", a.AuthMode)
	}

	if err != nil {
		return aws.Config{}, err
	}

	a.cfg = &cfg
	return cfg, nil
}

// awsRegions is a static sorted list of AWS regions.
// This avoids a dependency on SDK v1's endpoints package.
var awsRegions = []string{
	"af-south-1",
	"ap-east-1",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-northeast-3",
	"ap-south-1",
	"ap-south-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-southeast-3",
	"ap-southeast-4",
	"ap-southeast-5",
	"ap-southeast-7",
	"ca-central-1",
	"ca-west-1",
	"eu-central-1",
	"eu-central-2",
	"eu-north-1",
	"eu-south-1",
	"eu-south-2",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"il-central-1",
	"me-central-1",
	"me-south-1",
	"mx-central-1",
	"sa-east-1",
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
}

// GetRegionID returns a sorted list of AWS region IDs.
func GetRegionID() ([]string, error) {
	regions := make([]string, len(awsRegions))
	copy(regions, awsRegions)
	sort.Strings(regions)
	return regions, nil
}

// ListCluster lists EKS clusters in the configured region.
func (a *AWS) ListCluster() ([]ClusterInfo, error) {
	cfg, err := a.getAWSConfig()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	svc := eks.NewFromConfig(cfg)
	input := &eks.ListClustersInput{}

	result, err := svc.ListClusters(ctx, input)
	if err != nil {
		return nil, err
	}

	var clusterList []ClusterInfo
	for _, clusterName := range result.Clusters {
		clusterInfo, err := a.getClusterInfo(clusterName)
		if err != nil {
			return nil, err
		}
		clusterList = append(clusterList, clusterInfo)
	}
	return clusterList, nil
}

// getClusterInfo returns info for a single EKS cluster.
func (a *AWS) getClusterInfo(clusterName string) (ClusterInfo, error) {
	cfg, err := a.getAWSConfig()
	if err != nil {
		return ClusterInfo{}, err
	}

	ctx := context.Background()

	stsSvc := sts.NewFromConfig(cfg)
	callerIdentity, err := stsSvc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return ClusterInfo{}, err
	}

	eksSvc := eks.NewFromConfig(cfg)
	cluster, err := eksSvc.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: &clusterName,
	})
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
	}, nil
}

// GetKubeConfigObj returns a kubeconfig object for the given EKS cluster.
// If a Profile is set, --profile is added to the exec args so the generated
// kubeconfig works without needing AWS_PROFILE in the environment.
func (a *AWS) GetKubeConfigObj(clusterID string) (*clientcmdapi.Config, error) {
	cfg, err := a.getAWSConfig()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	svc := eks.NewFromConfig(cfg)
	cluster, err := svc.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: &clusterID,
	})
	if err != nil {
		return nil, err
	}

	decodePem, err := base64.StdEncoding.DecodeString(*cluster.Cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, err
	}

	execArgs := []string{
		"eks",
		"get-token",
		"--cluster-name",
		*cluster.Cluster.Name,
		"--region",
		a.RegionID,
		"--output",
		"json",
	}
	if a.Profile != "" {
		execArgs = append(execArgs, "--profile", a.Profile)
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
					Args:       execArgs,
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
