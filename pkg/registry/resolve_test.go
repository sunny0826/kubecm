package registry

import (
	"testing"
)

// --- Tests using deprecated Fragment aliases (backward compat) ---

func TestResolveFragment_Static(t *testing.T) {
	frag := &Fragment{
		Metadata: RegistryMetadata{Name: "test-static"},
		Provider: "static",
		Kubeconfig: `apiVersion: v1
kind: Config
clusters:
  - cluster:
      server: https://k8s.internal:6443
    name: onprem
contexts:
  - context:
      cluster: onprem
      user: onprem
    name: onprem
users:
  - name: onprem
    user:
      token: test-token
current-context: onprem
`,
	}

	cfg, err := ResolveFragment(frag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := cfg.Clusters["onprem"]; !ok {
		t.Error("expected cluster 'onprem' in config")
	}
	if _, ok := cfg.Contexts["onprem"]; !ok {
		t.Error("expected context 'onprem' in config")
	}
	if _, ok := cfg.AuthInfos["onprem"]; !ok {
		t.Error("expected authinfo 'onprem' in config")
	}
	if cfg.Clusters["onprem"].Server != "https://k8s.internal:6443" {
		t.Errorf("server = %q, want %q", cfg.Clusters["onprem"].Server, "https://k8s.internal:6443")
	}
}

func TestResolveFragment_StaticEmpty(t *testing.T) {
	frag := &Fragment{
		Metadata:   RegistryMetadata{Name: "empty"},
		Provider:   "static",
		Kubeconfig: "",
	}

	_, err := ResolveFragment(frag)
	if err == nil {
		t.Error("expected error for empty static kubeconfig")
	}
}

func TestResolveFragment_UnsupportedProvider(t *testing.T) {
	frag := &Fragment{
		Metadata: RegistryMetadata{Name: "unknown"},
		Provider: "gcp",
	}

	_, err := ResolveFragment(frag)
	if err == nil {
		t.Error("expected error for unsupported provider")
	}
}

func TestResolveFragment_AWSMissingSection(t *testing.T) {
	frag := &Fragment{
		Metadata: RegistryMetadata{Name: "bad-aws"},
		Provider: "aws",
	}

	_, err := ResolveFragment(frag)
	if err == nil {
		t.Error("expected error for aws fragment without aws section")
	}
}

func TestResolveFragment_AzureMissingSection(t *testing.T) {
	frag := &Fragment{
		Metadata: RegistryMetadata{Name: "bad-azure"},
		Provider: "azure",
	}

	_, err := ResolveFragment(frag)
	if err == nil {
		t.Error("expected error for azure fragment without azure section")
	}
}

func TestResolveFragmentWithUser_Nil(t *testing.T) {
	frag := &Fragment{
		Metadata:   RegistryMetadata{Name: "test"},
		Provider:   "static",
		Kubeconfig: staticKubeconfig("https://k8s:6443", "original-token"),
	}

	cfg, err := ResolveFragmentWithUser(frag, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AuthInfos["test-cluster"].Token != "original-token" {
		t.Errorf("token = %q, want %q", cfg.AuthInfos["test-cluster"].Token, "original-token")
	}
}

func TestResolveFragmentWithUser_StaticOverride(t *testing.T) {
	frag := &Fragment{
		Metadata:   RegistryMetadata{Name: "test"},
		Provider:   "static",
		Kubeconfig: staticKubeconfig("https://k8s:6443", "tok"),
	}
	user := &User{
		Metadata: RegistryMetadata{Name: "test-user"},
		Provider: "static",
	}

	cfg, err := ResolveFragmentWithUser(frag, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Clusters["test-cluster"]; !ok {
		t.Error("expected cluster 'test-cluster' in config")
	}
}

func TestResolveFragmentWithUser_ProviderMismatch(t *testing.T) {
	frag := &Fragment{
		Metadata: RegistryMetadata{Name: "aws-frag"},
		Provider: "aws",
		AWS:      &AWSFragment{Region: "eu-west-1", Cluster: "test"},
	}
	user := &User{
		Metadata: RegistryMetadata{Name: "azure-user"},
		Provider: "azure",
	}

	_, err := ResolveFragmentWithUser(frag, user)
	if err == nil {
		t.Error("expected error for provider mismatch")
	}
}

func TestCloneFragment_NoMutation(t *testing.T) {
	orig := &Fragment{
		Metadata: RegistryMetadata{Name: "orig"},
		Provider: "aws",
		AWS: &AWSFragment{
			Region:  "eu-west-1",
			Cluster: "test",
			Profile: "original-profile",
		},
	}

	clone := cloneCluster(orig)
	clone.AWS.Profile = "new-profile"
	clone.Metadata.Name = "clone"

	if orig.AWS.Profile != "original-profile" {
		t.Errorf("original AWS profile mutated: got %q", orig.AWS.Profile)
	}
	if orig.Metadata.Name != "orig" {
		t.Errorf("original metadata mutated: got %q", orig.Metadata.Name)
	}
}

func TestCloneFragment_AzureNoMutation(t *testing.T) {
	orig := &Fragment{
		Metadata: RegistryMetadata{Name: "orig"},
		Provider: "azure",
		Azure: &AzureFragment{
			SubscriptionID: "sub-123",
			ResourceGroup:  "rg",
			Cluster:        "aks",
			TenantID:       "original-tenant",
		},
	}

	clone := cloneCluster(orig)
	clone.Azure.TenantID = "new-tenant"

	if orig.Azure.TenantID != "original-tenant" {
		t.Errorf("original Azure tenantID mutated: got %q", orig.Azure.TenantID)
	}
}

// --- Tests using new Cluster types (new API) ---

func TestResolveCluster_Static(t *testing.T) {
	cl := &Cluster{
		Metadata: RegistryMetadata{Name: "test-static"},
		Provider: "static",
		Kubeconfig: staticKubeconfig("https://k8s.internal:6443", "test-token"),
	}

	cfg, err := ResolveCluster(cl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Clusters["test-cluster"]; !ok {
		t.Error("expected cluster 'test-cluster' in config")
	}
}

func TestResolveClusterWithUser_Nil(t *testing.T) {
	cl := &Cluster{
		Metadata:   RegistryMetadata{Name: "test"},
		Provider:   "static",
		Kubeconfig: staticKubeconfig("https://k8s:6443", "original-token"),
	}

	cfg, err := ResolveClusterWithUser(cl, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AuthInfos["test-cluster"].Token != "original-token" {
		t.Errorf("token = %q, want %q", cfg.AuthInfos["test-cluster"].Token, "original-token")
	}
}

func TestResolveClusterWithUser_ProviderMismatch(t *testing.T) {
	cl := &Cluster{
		Metadata: RegistryMetadata{Name: "aws-cl"},
		Provider: "aws",
		AWS:      &AWSClusterConfig{Region: "eu-west-1", Cluster: "test"},
	}
	user := &User{
		Metadata: RegistryMetadata{Name: "azure-user"},
		Provider: "azure",
	}

	_, err := ResolveClusterWithUser(cl, user)
	if err == nil {
		t.Error("expected error for provider mismatch")
	}
}

func TestCloneCluster_NoMutation(t *testing.T) {
	orig := &Cluster{
		Metadata: RegistryMetadata{Name: "orig"},
		Provider: "aws",
		AWS: &AWSClusterConfig{
			Region:  "eu-west-1",
			Cluster: "test",
			Profile: "original-profile",
		},
	}

	clone := cloneCluster(orig)
	clone.AWS.Profile = "new-profile"

	if orig.AWS.Profile != "original-profile" {
		t.Errorf("original AWS profile mutated: got %q", orig.AWS.Profile)
	}
}

// staticKubeconfig returns a minimal valid kubeconfig YAML for testing.
func staticKubeconfig(server, token string) string {
	return `apiVersion: v1
kind: Config
clusters:
  - cluster:
      server: ` + server + `
    name: test-cluster
contexts:
  - context:
      cluster: test-cluster
      user: test-cluster
    name: test-cluster
users:
  - name: test-cluster
    user:
      token: ` + token + `
`
}
