package cloud

// Cluster interface of cloud k8s cluster
type Cluster interface {
	ListCluster() (clusters []ClusterInfo, err error)
	GetKubeConfig(clusterID string) (kubeconfig string, err error)
}

// ClusterInfo ack cluster info
type ClusterInfo struct {
	Name     string
	ID       string
	RegionID string
}
