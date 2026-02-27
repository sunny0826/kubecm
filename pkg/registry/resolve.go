package registry

import (
	"fmt"

	"github.com/sunny0826/kubecm/pkg/cloud"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd"
)

// ResolveClusterWithUser resolves a cluster, optionally overriding
// credentials with a User definition. If user is nil, behaves like ResolveCluster.
func ResolveClusterWithUser(cl *Cluster, user *User) (*clientcmdapi.Config, error) {
	if user == nil {
		return ResolveCluster(cl)
	}
	if user.Provider != cl.Provider {
		return nil, fmt.Errorf("provider mismatch: cluster %q is %q but user %q is %q",
			cl.Metadata.Name, cl.Provider, user.Metadata.Name, user.Provider)
	}
	merged := cloneCluster(cl)
	switch user.Provider {
	case "aws":
		if user.AWS != nil {
			if merged.AWS == nil {
				merged.AWS = &AWSClusterConfig{}
			}
			merged.AWS.Profile = user.AWS.Profile
		}
	case "azure":
		if user.Azure != nil {
			if merged.Azure == nil {
				merged.Azure = &AzureClusterConfig{}
			}
			merged.Azure.TenantID = user.Azure.TenantID
		}
	}
	return ResolveCluster(merged)
}

// cloneCluster returns a shallow copy of the cluster with deep-copied
// provider sections so mutations don't affect the original.
func cloneCluster(cl *Cluster) *Cluster {
	c := *cl
	if cl.AWS != nil {
		cp := *cl.AWS
		c.AWS = &cp
	}
	if cl.Azure != nil {
		cp := *cl.Azure
		c.Azure = &cp
	}
	return &c
}

// ResolveCluster fetches the kubeconfig for a cluster by calling
// the appropriate cloud provider API or parsing static kubeconfig.
func ResolveCluster(cl *Cluster) (*clientcmdapi.Config, error) {
	switch cl.Provider {
	case "aws":
		return resolveAWS(cl)
	case "azure":
		return resolveAzure(cl)
	case "static":
		return resolveStatic(cl)
	default:
		return nil, fmt.Errorf("unsupported provider %q in cluster %q", cl.Provider, cl.Metadata.Name)
	}
}

func resolveAWS(cl *Cluster) (*clientcmdapi.Config, error) {
	if cl.AWS == nil {
		return nil, fmt.Errorf("cluster %q: provider is aws but aws section is missing", cl.Metadata.Name)
	}

	a := cloud.AWS{
		AuthMode: cloud.AWSAuthDefault,
		Profile:  cl.AWS.Profile,
		RegionID: cl.AWS.Region,
	}

	cfg, err := a.GetKubeConfigObj(cl.AWS.Cluster)
	if err != nil {
		return nil, fmt.Errorf("cluster %q: aws: %w", cl.Metadata.Name, err)
	}
	return cfg, nil
}

func resolveAzure(cl *Cluster) (*clientcmdapi.Config, error) {
	if cl.Azure == nil {
		return nil, fmt.Errorf("cluster %q: provider is azure but azure section is missing", cl.Metadata.Name)
	}

	a := cloud.Azure{
		AuthMode:       cloud.AuthModeDefault,
		SubscriptionID: cl.Azure.SubscriptionID,
		TenantID:       cl.Azure.TenantID,
	}

	data, err := a.GetKubeConfig(cl.Azure.Cluster, cl.Azure.ResourceGroup)
	if err != nil {
		return nil, fmt.Errorf("cluster %q: azure: %w", cl.Metadata.Name, err)
	}

	cfg, err := clientcmd.Load(data)
	if err != nil {
		return nil, fmt.Errorf("cluster %q: parsing azure kubeconfig: %w", cl.Metadata.Name, err)
	}
	return cfg, nil
}

func resolveStatic(cl *Cluster) (*clientcmdapi.Config, error) {
	if cl.Kubeconfig == "" {
		return nil, fmt.Errorf("cluster %q: provider is static but kubeconfig is empty", cl.Metadata.Name)
	}

	cfg, err := clientcmd.Load([]byte(cl.Kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("cluster %q: parsing static kubeconfig: %w", cl.Metadata.Name, err)
	}
	return cfg, nil
}

// Deprecated aliases for backward compatibility.
var ResolveFragment = ResolveCluster
var ResolveFragmentWithUser = ResolveClusterWithUser
