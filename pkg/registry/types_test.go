package registry

import (
	"testing"
)

func TestRole_NormalizedContexts_Fragments(t *testing.T) {
	role := &Role{
		Fragments: []string{"cluster-a", "cluster-b"},
	}

	got := role.NormalizedContexts()
	if len(got) != 2 {
		t.Fatalf("expected 2 contexts, got %d", len(got))
	}
	if got[0].ClusterRef() != "cluster-a" || got[0].User != "" {
		t.Errorf("context[0] = %+v, want cluster-a", got[0])
	}
	if got[1].ClusterRef() != "cluster-b" || got[1].User != "" {
		t.Errorf("context[1] = %+v, want cluster-b", got[1])
	}
}

func TestRole_NormalizedContexts_Contexts(t *testing.T) {
	role := &Role{
		Contexts: []RoleContext{
			{Cluster: "prod", User: "admin", Name: "prod-admin"},
			{Cluster: "prod", User: "readonly", Name: "prod-ro"},
		},
	}

	got := role.NormalizedContexts()
	if len(got) != 2 {
		t.Fatalf("expected 2 contexts, got %d", len(got))
	}
	if got[0].User != "admin" || got[0].Name != "prod-admin" {
		t.Errorf("context[0] = %+v", got[0])
	}
	if got[1].User != "readonly" || got[1].Name != "prod-ro" {
		t.Errorf("context[1] = %+v", got[1])
	}
}

func TestRole_NormalizedContexts_Both(t *testing.T) {
	// When both are set, contexts: takes priority
	role := &Role{
		Fragments: []string{"old-a", "old-b"},
		Contexts: []RoleContext{
			{Cluster: "new-a"},
		},
	}

	got := role.NormalizedContexts()
	if len(got) != 1 {
		t.Fatalf("expected 1 context (contexts takes priority), got %d", len(got))
	}
	if got[0].ClusterRef() != "new-a" {
		t.Errorf("context[0].ClusterRef() = %q, want %q", got[0].ClusterRef(), "new-a")
	}
}

func TestRoleContext_ClusterRef(t *testing.T) {
	t.Run("cluster field", func(t *testing.T) {
		rc := RoleContext{Cluster: "prod"}
		if rc.ClusterRef() != "prod" {
			t.Errorf("ClusterRef() = %q, want %q", rc.ClusterRef(), "prod")
		}
	})

	t.Run("fragment field (deprecated)", func(t *testing.T) {
		rc := RoleContext{Fragment: "legacy"}
		if rc.ClusterRef() != "legacy" {
			t.Errorf("ClusterRef() = %q, want %q", rc.ClusterRef(), "legacy")
		}
	})

	t.Run("cluster takes priority over fragment", func(t *testing.T) {
		rc := RoleContext{Cluster: "new", Fragment: "old"}
		if rc.ClusterRef() != "new" {
			t.Errorf("ClusterRef() = %q, want %q", rc.ClusterRef(), "new")
		}
	})
}

func TestValidateRoleContexts_OK(t *testing.T) {
	role := &Role{
		Metadata: RegistryMetadata{Name: "devops"},
		Contexts: []RoleContext{
			{Cluster: "prod", User: "admin", Name: "prod-admin"},
			{Cluster: "prod", User: "readonly", Name: "prod-ro"},
			{Cluster: "staging"},
		},
	}
	if err := ValidateRoleContexts(role); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateRoleContexts_DuplicateNames(t *testing.T) {
	role := &Role{
		Metadata: RegistryMetadata{Name: "devops"},
		Contexts: []RoleContext{
			{Cluster: "prod", User: "admin"},
			{Cluster: "prod", User: "readonly"},
		},
	}
	err := ValidateRoleContexts(role)
	if err == nil {
		t.Error("expected error for duplicate context names")
	}
}

func TestValidateRoleContexts_DuplicateWithPrefix(t *testing.T) {
	role := &Role{
		Metadata:      RegistryMetadata{Name: "devops"},
		ContextPrefix: "acme",
		Contexts: []RoleContext{
			{Cluster: "prod", User: "admin"},
			{Cluster: "prod", User: "readonly"},
		},
	}
	err := ValidateRoleContexts(role)
	if err == nil {
		t.Error("expected error for duplicate context names with prefix")
	}
}

func TestValidateRoleContexts_Empty(t *testing.T) {
	role := &Role{
		Metadata: RegistryMetadata{Name: "empty"},
	}
	err := ValidateRoleContexts(role)
	if err == nil {
		t.Error("expected error for empty role")
	}
}

func TestValidateRoleContexts_LegacyFragmentField(t *testing.T) {
	// Validate works with deprecated fragment: field too
	role := &Role{
		Metadata: RegistryMetadata{Name: "legacy"},
		Contexts: []RoleContext{
			{Fragment: "prod"},
			{Fragment: "staging"},
		},
	}
	if err := ValidateRoleContexts(role); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
