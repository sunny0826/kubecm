package registry

import (
	"fmt"

	"github.com/sunny0826/kubecm/pkg/cloud"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd"
)

// ResolveFragment fetches the kubeconfig for a fragment by calling
// the appropriate cloud provider API or parsing static kubeconfig.
func ResolveFragment(frag *Fragment) (*clientcmdapi.Config, error) {
	switch frag.Provider {
	case "aws":
		return resolveAWS(frag)
	case "azure":
		return resolveAzure(frag)
	case "static":
		return resolveStatic(frag)
	default:
		return nil, fmt.Errorf("unsupported provider %q in fragment %q", frag.Provider, frag.Metadata.Name)
	}
}

func resolveAWS(frag *Fragment) (*clientcmdapi.Config, error) {
	if frag.AWS == nil {
		return nil, fmt.Errorf("fragment %q: provider is aws but aws section is missing", frag.Metadata.Name)
	}

	a := cloud.AWS{
		AuthMode: cloud.AWSAuthDefault,
		Profile:  frag.AWS.Profile,
		RegionID: frag.AWS.Region,
	}

	cfg, err := a.GetKubeConfigObj(frag.AWS.Cluster)
	if err != nil {
		return nil, fmt.Errorf("fragment %q: aws: %w", frag.Metadata.Name, err)
	}
	return cfg, nil
}

func resolveAzure(frag *Fragment) (*clientcmdapi.Config, error) {
	if frag.Azure == nil {
		return nil, fmt.Errorf("fragment %q: provider is azure but azure section is missing", frag.Metadata.Name)
	}

	a := cloud.Azure{
		AuthMode:       cloud.AuthModeDefault,
		SubscriptionID: frag.Azure.SubscriptionID,
		TenantID:       frag.Azure.TenantID,
	}

	data, err := a.GetKubeConfig(frag.Azure.Cluster, frag.Azure.ResourceGroup)
	if err != nil {
		return nil, fmt.Errorf("fragment %q: azure: %w", frag.Metadata.Name, err)
	}

	cfg, err := clientcmd.Load(data)
	if err != nil {
		return nil, fmt.Errorf("fragment %q: parsing azure kubeconfig: %w", frag.Metadata.Name, err)
	}
	return cfg, nil
}

func resolveStatic(frag *Fragment) (*clientcmdapi.Config, error) {
	if frag.Kubeconfig == "" {
		return nil, fmt.Errorf("fragment %q: provider is static but kubeconfig is empty", frag.Metadata.Name)
	}

	cfg, err := clientcmd.Load([]byte(frag.Kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("fragment %q: parsing static kubeconfig: %w", frag.Metadata.Name, err)
	}
	return cfg, nil
}
