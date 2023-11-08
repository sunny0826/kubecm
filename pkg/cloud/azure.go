package cloud

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
)

// Azure struct of azure cloud
type Azure struct {
	AuthMode       AzureAuth
	ClientID       string
	ClientSecret   string
	SubscriptionID string
	TenantID       string
	ObjectID       string

	client azcore.TokenCredential
}
type AzureAuth int

type AzureSubscription struct {
	ID          string
	DisplayName string
}

const (
	AuthModeDefault AzureAuth = iota
	AuthModeServicePrincipal
)

func (a *Azure) getAzureClient() (azcore.TokenCredential, error) {
	if a.client != nil {
		return a.client, nil
	}

	var (
		err    error
		client azcore.TokenCredential
	)

	switch a.AuthMode {
	case AuthModeDefault:
		client, err = azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
			TenantID: a.TenantID,
		})
	case AuthModeServicePrincipal:
		client, err = azidentity.NewClientSecretCredential(a.TenantID, a.ClientID, a.ClientSecret, nil)
	default:
		return nil, fmt.Errorf("invalid auth mode: %d", a.AuthMode)
	}

	a.client = client

	return client, err
}

// ListSubscriptions list subscriptions
func (a *Azure) ListSubscriptions() (subscription []AzureSubscription, err error) {
	client, err := a.getAzureClient()
	if err != nil {
		return nil, err
	}

	subscriptionClient, err := armsubscription.NewSubscriptionsClient(client, nil)
	if err != nil {
		return nil, err
	}

	var subscriptionList []AzureSubscription

	pager := subscriptionClient.NewListPager(nil)

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, subscription := range page.Value {
			subscriptionList = append(subscriptionList, AzureSubscription{
				DisplayName: *subscription.DisplayName,
				ID:          *subscription.SubscriptionID,
			})
		}
	}

	return subscriptionList, err
}

// ListCluster list cluster info
func (a *Azure) ListCluster(subscription AzureSubscription) (clusters []ClusterInfo, err error) {
	client, err := a.getAzureClient()
	if err != nil {
		return nil, err
	}

	aksClient, err := armcontainerservice.NewManagedClustersClient(subscription.ID, client, nil)
	if err != nil {
		return nil, err
	}

	var clusterList []ClusterInfo

	pager := aksClient.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, cluster := range page.Value {
			clusterList = append(clusterList, ClusterInfo{
				Name:       *cluster.Name,
				Account:    subscription.DisplayName,
				ID:         *cluster.ID,
				RegionID:   *cluster.Location,
				K8sVersion: *cluster.Properties.KubernetesVersion,
				ConsoleURL: "https://portal.azure.com",
				// Too long to CLI
				// ConsoleURL: fmt.Sprintf("https://portal.azure.com/#resource%s/overview", *cluster.ID),
			})
		}
	}

	return clusterList, err
}

// GetKubeConfig get kubeConfig file
func (a *Azure) GetKubeConfig(clusterName, resourceGroupName string) ([]byte, error) {
	client, err := a.getAzureClient()
	if err != nil {
		return nil, err
	}

	aksClient, err := armcontainerservice.NewManagedClustersClient(a.SubscriptionID, client, nil)
	if err != nil {
		return nil, err
	}

	res, err := aksClient.ListClusterUserCredentials(context.Background(), resourceGroupName, clusterName, nil)
	if err != nil {
		return nil, err
	}

	kubeconfig := res.Kubeconfigs
	for _, v := range kubeconfig {
		return v.Value, nil
	}
	return nil, nil
}

// GetAdminKubeConfig get kubeConfig file
func (a *Azure) GetAdminKubeConfig(clusterName, resourceGroupName string) ([]byte, error) {
	client, err := a.getAzureClient()
	if err != nil {
		return nil, err
	}

	aksClient, err := armcontainerservice.NewManagedClustersClient(a.SubscriptionID, client, nil)
	if err != nil {
		return nil, err
	}

	res, err := aksClient.ListClusterAdminCredentials(context.Background(), resourceGroupName, clusterName, nil)
	if err != nil {
		return nil, err
	}

	kubeconfig := res.Kubeconfigs
	for _, v := range kubeconfig {
		return v.Value, nil
	}
	return nil, nil
}
