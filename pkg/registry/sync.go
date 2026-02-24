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
// It loads the role, resolves each fragment, and returns
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

	for _, fragName := range role.Fragments {
		frag, err := LoadFragment(repoDir, fragName)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("loading fragment %q: %v", fragName, err))
			continue
		}

		// Apply template variables
		if err := ResolveFragmentTemplates(frag, entry.Variables); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("template %q: %v", fragName, err))
			continue
		}

		// Resolve fragment -> kubeconfig
		fragConfig, err := ResolveFragment(frag)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("resolving %q: %v", fragName, err))
			continue
		}

		// Merge fragment kubeconfig into current config with prefix
		mergeFragmentConfig(currentConfig, fragConfig, role.ContextPrefix, frag.Metadata.Name, managedSet, newContexts, result, dryRun)
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

// mergeFragmentConfig merges a single fragment's kubeconfig into the current config.
func mergeFragmentConfig(
	current *clientcmdapi.Config,
	fragConfig *clientcmdapi.Config,
	contextPrefix string,
	fragmentName string,
	managedSet map[string]bool,
	newContexts map[string]bool,
	result *SyncResult,
	dryRun bool,
) {
	for origCtxName, ctx := range fragConfig.Contexts {
		// Build prefixed context name
		ctxName := buildContextName(contextPrefix, fragmentName, origCtxName)
		newContexts[ctxName] = true

		// Build prefixed cluster and user names
		clusterName := ctxName
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
		if origCluster, ok := fragConfig.Clusters[ctx.Cluster]; ok {
			current.Clusters[clusterName] = origCluster
		}

		// Copy user data with new name
		if origUser, ok := fragConfig.AuthInfos[ctx.AuthInfo]; ok {
			current.AuthInfos[userName] = origUser
		}

		// Create context with new names
		current.Contexts[ctxName] = &clientcmdapi.Context{
			Cluster:   clusterName,
			AuthInfo:  userName,
			Namespace: ctx.Namespace,
		}
	}
}

// buildContextName creates the full context name with prefix.
// Format: <prefix>-<fragmentName>, or just <fragmentName> if no prefix.
func buildContextName(prefix, fragmentName, _ string) string {
	if prefix == "" {
		return fragmentName
	}
	return prefix + "-" + fragmentName
}

// removeContext removes a context and its cluster/user if not referenced elsewhere.
func removeContext(config *clientcmdapi.Config, ctxName string) {
	ctx, ok := config.Contexts[ctxName]
	if !ok {
		return
	}

	clusterName := ctx.Cluster
	userName := ctx.AuthInfo

	delete(config.Contexts, ctxName)

	// Only delete cluster/user if no other context references them
	clusterUsed := false
	userUsed := false
	for _, c := range config.Contexts {
		if c.Cluster == clusterName {
			clusterUsed = true
		}
		if c.AuthInfo == userName {
			userUsed = true
		}
	}
	if !clusterUsed {
		delete(config.Clusters, clusterName)
	}
	if !userUsed {
		delete(config.AuthInfos, userName)
	}
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
