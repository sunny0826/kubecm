package registry

import (
	"fmt"
	"strings"
	"time"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// SyncResult holds the outcome of a sync operation.
type SyncResult struct {
	Added    []string
	Updated  []string
	Removed  []string
	Skipped  []string // conflicts: context exists but not managed
	Errors   []string
}

// Sync performs a full registry sync for the given entry.
// It loads the role, resolves each cluster, and returns
// the merged kubeconfig changes + updated managed contexts list.
func Sync(repoDir string, entry *RegistryEntry, currentConfig *clientcmdapi.Config, dryRun bool) (*SyncResult, error) {
	role, err := LoadRole(repoDir, entry.Role)
	if err != nil {
		return nil, fmt.Errorf("loading role: %w", err)
	}

	result := &SyncResult{}
	newContexts := make(map[string]bool)
	managedSet := make(map[string]bool)
	for _, ctx := range entry.ManagedContexts {
		managedSet[ctx] = true
	}

	// Validate role contexts before processing
	if err := ValidateRoleContexts(role); err != nil {
		return nil, fmt.Errorf("validating role: %w", err)
	}

	for _, rc := range role.NormalizedContexts() {
		clusterRef := rc.ClusterRef()

		cl, err := LoadCluster(repoDir, clusterRef)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("loading cluster %q: %v", clusterRef, err))
			continue
		}

		// Apply template variables to cluster
		if err := ResolveClusterTemplates(cl, entry.Variables); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("template %q: %v", clusterRef, err))
			continue
		}

		// Load and template user if specified
		var user *User
		if rc.User != "" {
			user, err = LoadUser(repoDir, rc.User)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("loading user %q: %v", rc.User, err))
				continue
			}
			if err := ResolveUserTemplates(user, entry.Variables); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("template user %q: %v", rc.User, err))
				continue
			}
		}

		// Resolve cluster with optional user override
		clConfig, err := ResolveClusterWithUser(cl, user)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("resolving %q: %v", clusterRef, err))
			continue
		}

		// Use explicit name if provided, otherwise cluster ref
		contextName := clusterRef
		if rc.Name != "" {
			contextName = rc.Name
		}

		// Merge cluster kubeconfig into current config with prefix
		mergeClusterConfig(currentConfig, clConfig, role.ContextPrefix, contextName, managedSet, newContexts, result, dryRun)
	}

	// Remove stale managed contexts (in managedSet but not in newContexts)
	for _, ctx := range entry.ManagedContexts {
		if !newContexts[ctx] {
			if _, exists := currentConfig.Contexts[ctx]; exists {
				if !dryRun {
					removeContext(currentConfig, ctx)
				}
				result.Removed = append(result.Removed, ctx)
			}
		}
	}

	// Update entry
	if !dryRun {
		var managed []string
		for ctx := range newContexts {
			managed = append(managed, ctx)
		}
		entry.ManagedContexts = managed
		now := time.Now().UTC()
		entry.LastSync = &now
	}

	return result, nil
}

// mergeClusterConfig merges a single cluster's kubeconfig into the current config.
func mergeClusterConfig(
	current *clientcmdapi.Config,
	clConfig *clientcmdapi.Config,
	contextPrefix string,
	clusterName string,
	managedSet map[string]bool,
	newContexts map[string]bool,
	result *SyncResult,
	dryRun bool,
) {
	for origCtxName, ctx := range clConfig.Contexts {
		// Build prefixed context name
		ctxName := buildContextName(contextPrefix, clusterName, origCtxName)
		newContexts[ctxName] = true

		// Build prefixed cluster and user names
		clName := ctxName
		userName := ctxName

		// Check for conflicts
		if _, exists := current.Contexts[ctxName]; exists {
			if !managedSet[ctxName] {
				// Context exists but is not managed by this registry -> skip
				result.Skipped = append(result.Skipped, ctxName)
				continue
			}
			// Managed context -> update (overwrite)
			result.Updated = append(result.Updated, ctxName)
		} else {
			result.Added = append(result.Added, ctxName)
		}

		if dryRun {
			continue
		}

		// Copy cluster data with new name
		if origCluster, ok := clConfig.Clusters[ctx.Cluster]; ok {
			current.Clusters[clName] = origCluster
		}

		// Copy user data with new name
		if origUser, ok := clConfig.AuthInfos[ctx.AuthInfo]; ok {
			current.AuthInfos[userName] = origUser
		}

		// Create context with new names
		current.Contexts[ctxName] = &clientcmdapi.Context{
			Cluster:   clName,
			AuthInfo:  userName,
			Namespace: ctx.Namespace,
		}
	}
}

// buildContextName creates the full context name with prefix.
// Format: <prefix>-<name>, or just <name> if no prefix.
func buildContextName(prefix, name, _ string) string {
	if prefix == "" {
		return name
	}
	return prefix + "-" + name
}

// removeContext removes a context and its cluster/user if not referenced elsewhere.
func removeContext(config *clientcmdapi.Config, ctxName string) {
	ctx, ok := config.Contexts[ctxName]
	if !ok {
		return
	}

	clName := ctx.Cluster
	userName := ctx.AuthInfo

	delete(config.Contexts, ctxName)

	// Only delete cluster/user if no other context references them
	clusterUsed := false
	userUsed := false
	for _, c := range config.Contexts {
		if c.Cluster == clName {
			clusterUsed = true
		}
		if c.AuthInfo == userName {
			userUsed = true
		}
	}
	if !clusterUsed {
		delete(config.Clusters, clName)
	}
	if !userUsed {
		delete(config.AuthInfos, userName)
	}
}

// ValidateRoleContexts checks that a role's contexts produce unique names.
// It returns an error if the same context name would appear twice.
func ValidateRoleContexts(role *Role) error {
	contexts := role.NormalizedContexts()
	if len(contexts) == 0 {
		return fmt.Errorf("role %q has no clusters or contexts defined", role.Metadata.Name)
	}
	seen := make(map[string]bool)
	for _, rc := range contexts {
		name := rc.ClusterRef()
		if rc.Name != "" {
			name = rc.Name
		}
		fullName := buildContextName(role.ContextPrefix, name, "")
		if seen[fullName] {
			return fmt.Errorf("role %q: duplicate context name %q (use the name: field to disambiguate)", role.Metadata.Name, fullName)
		}
		seen[fullName] = true
	}
	return nil
}

// FormatSyncResult returns a human-readable summary.
func FormatSyncResult(r *SyncResult) string {
	var sb strings.Builder

	if len(r.Added) > 0 {
		sb.WriteString("  Added:\n")
		for _, c := range r.Added {
			fmt.Fprintf(&sb, "    + %s\n", c)
		}
	}
	if len(r.Updated) > 0 {
		sb.WriteString("  Updated:\n")
		for _, c := range r.Updated {
			fmt.Fprintf(&sb, "    ~ %s\n", c)
		}
	}
	if len(r.Removed) > 0 {
		sb.WriteString("  Removed:\n")
		for _, c := range r.Removed {
			fmt.Fprintf(&sb, "    - %s\n", c)
		}
	}
	if len(r.Skipped) > 0 {
		sb.WriteString("  Skipped (conflict):\n")
		for _, c := range r.Skipped {
			fmt.Fprintf(&sb, "    ! %s\n", c)
		}
	}
	if len(r.Errors) > 0 {
		sb.WriteString("  Errors:\n")
		for _, e := range r.Errors {
			fmt.Fprintf(&sb, "    ERROR: %s\n", e)
		}
	}

	if sb.Len() == 0 {
		sb.WriteString("  No changes.\n")
	}
	return sb.String()
}
