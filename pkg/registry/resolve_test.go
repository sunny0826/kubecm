package registry

import (
	"testing"
)

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
