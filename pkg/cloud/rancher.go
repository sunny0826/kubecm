package cloud

import (
	"strings"

	"github.com/rancher/norman/clientbase"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Rancher struct of rancher
type Rancher struct {
	ServerURL string
	APIKey    string
}

// getRancherClient get rancher client
func getRancherClient(serverURL, apiKey string) (*managementClient.Client, error) {
	if !strings.HasSuffix(serverURL, "/v3") {
		serverURL += "/v3"
	}

	options := &clientbase.ClientOpts{
		URL:      serverURL,
		TokenKey: apiKey,
		Insecure: true,
	}

	result, err := managementClient.NewClient(options)
	if err != nil {
		return nil, err
	}
	return result, err
}

// GetRegionID get region id of rancher cluster
func (r *Rancher) GetRegionID() ([]string, error) {
	// rancher No RegionID required, return nil
	return nil, nil
}

// ListCluster list cluster info
func (r *Rancher) ListCluster() (clusters []ClusterInfo, err error) {
	client, err := getRancherClient(r.ServerURL, r.APIKey)
	if err != nil {
		return nil, err
	}
	clusterCollection, err := client.Cluster.ListAll(nil)
	if err != nil {
		return nil, err
	}
	var clusterList []ClusterInfo
	for _, info := range clusterCollection.Data {
		clusterList = append(clusterList, ClusterInfo{
			Name:       info.Name,
			ID:         info.ID,
			RegionID:   "",
			K8sVersion: info.Version.GitVersion,
		})
	}
	return clusterList, err
}

// GetKubeConfig get kubeConfig file
func (r *Rancher) GetKubeConfig(clusterID string) (string, error) {
	client, err := getRancherClient(r.ServerURL, r.APIKey)
	if err != nil {
		return "", err
	}

	cluster, err := client.Cluster.ByID(clusterID)
	if err != nil {
		return "", err
	}

	config, err := client.Cluster.ActionGenerateKubeconfig(cluster)
	if err != nil {
		return "", err
	}
	return config.Config, nil
}
