package cloud

// Cluster interface of cloud k8s cluster
type Cluster interface {
	GetRegionID() ([]string, error)
	ListCluster() (clusters []ClusterInfo, err error)
	GetKubeConfig(clusterID string) (kubeconfig string, err error)
}

// ClusterInfo ack cluster info
type ClusterInfo struct {
	Name       string
	Account    string
	ID         string
	RegionID   string
	K8sVersion string
	ConsoleURL string
}
