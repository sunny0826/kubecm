package cloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// Azure struct of azure cloud
type Azure struct {
	ClientID       string
	ClientSecret   string
	SubscriptionID string
	TenantID       string
	ObjectID       string
}

func getAzureClient(clientID, clientSecret, TenantID string) (autorest.Authorizer, error) {
	credConfig := auth.NewClientCredentialsConfig(clientID, clientSecret, TenantID)
	return credConfig.Authorizer()
}

// ListCluster list cluster info
func (a *Azure) ListCluster() (clusters []ClusterInfo, err error) {
	authorizer, err := getAzureClient(a.ClientID, a.ClientSecret, a.TenantID)
	if err != nil {
		return nil, err
	}

	aksClient := containerservice.NewManagedClustersClient(a.SubscriptionID)
	aksClient.Authorizer = authorizer
	list, err := aksClient.List(context.Background())
	if err != nil {
		return nil, err
	}
	var clusterList []ClusterInfo
	for _, cluster := range list.Values() {
		clusterList = append(clusterList, ClusterInfo{
			Name:       *cluster.Name,
			ID:         getResourceGroupName(*cluster.ID),
			RegionID:   *cluster.Location,
			K8sVersion: *cluster.ManagedClusterProperties.KubernetesVersion,
			ConsoleURL: "https://portal.azure.com",
		})
	}
	return clusterList, err
}

// GetKubeConfig get kubeConfig file
func (a *Azure) GetKubeConfig(clusterName, resourceGroupName string) ([]byte, error) {
	authorizer, err := getAzureClient(a.ClientID, a.ClientSecret, a.TenantID)
	if err != nil {
		return nil, err
	}

	aksClient := containerservice.NewManagedClustersClient(a.SubscriptionID)
	aksClient.Authorizer = authorizer
	res, err := aksClient.ListClusterAdminCredentials(context.Background(), resourceGroupName, clusterName, "")
	if err != nil {
		return nil, err
	}
	kubeconfig := *(res.Kubeconfigs)
	for _, v := range kubeconfig {
		return (*(v.Value)), nil
	}
	return nil, nil
}

func getResourceGroupName(id string) string {
	segments := strings.Split(id, "/")
	resourceGroupIndex := -1

	for i, segment := range segments {
		if segment == "resourcegroups" {
			resourceGroupIndex = i
			break
		}
	}

	if resourceGroupIndex != -1 && resourceGroupIndex+1 < len(segments) {
		resourceGroupName := segments[resourceGroupIndex+1]
		return resourceGroupName
	}
	fmt.Println("can not found resourcegroups")
	return ""
}
