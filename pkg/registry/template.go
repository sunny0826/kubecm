package registry

import (
	"bytes"
	"fmt"
	"text/template"
)

// ResolveTemplate applies Go template variables to a string.
// Variables are accessed as {{ .Key }}.
func ResolveTemplate(text string, vars map[string]string) (string, error) {
	if len(vars) == 0 {
		return text, nil
	}

	tmpl, err := template.New("cluster").Option("missingkey=error").Parse(text)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}
	return buf.String(), nil
}

// ResolveClusterTemplates applies template variables to all string fields in a Cluster.
func ResolveClusterTemplates(cl *Cluster, vars map[string]string) error {
	if len(vars) == 0 {
		return nil
	}

	var err error

	if cl.AWS != nil {
		if cl.AWS.Region, err = ResolveTemplate(cl.AWS.Region, vars); err != nil {
			return fmt.Errorf("aws.region: %w", err)
		}
		if cl.AWS.Cluster, err = ResolveTemplate(cl.AWS.Cluster, vars); err != nil {
			return fmt.Errorf("aws.cluster: %w", err)
		}
		if cl.AWS.Profile, err = ResolveTemplate(cl.AWS.Profile, vars); err != nil {
			return fmt.Errorf("aws.profile: %w", err)
		}
	}

	if cl.Azure != nil {
		if cl.Azure.SubscriptionID, err = ResolveTemplate(cl.Azure.SubscriptionID, vars); err != nil {
			return fmt.Errorf("azure.subscriptionId: %w", err)
		}
		if cl.Azure.ResourceGroup, err = ResolveTemplate(cl.Azure.ResourceGroup, vars); err != nil {
			return fmt.Errorf("azure.resourceGroup: %w", err)
		}
		if cl.Azure.Cluster, err = ResolveTemplate(cl.Azure.Cluster, vars); err != nil {
			return fmt.Errorf("azure.cluster: %w", err)
		}
		if cl.Azure.TenantID, err = ResolveTemplate(cl.Azure.TenantID, vars); err != nil {
			return fmt.Errorf("azure.tenantId: %w", err)
		}
	}

	if cl.Kubeconfig != "" {
		if cl.Kubeconfig, err = ResolveTemplate(cl.Kubeconfig, vars); err != nil {
			return fmt.Errorf("kubeconfig: %w", err)
		}
	}

	return nil
}

// ResolveUserTemplates applies template variables to all string fields in a User.
func ResolveUserTemplates(u *User, vars map[string]string) error {
	if len(vars) == 0 {
		return nil
	}

	var err error

	if u.AWS != nil {
		if u.AWS.Profile, err = ResolveTemplate(u.AWS.Profile, vars); err != nil {
			return fmt.Errorf("aws.profile: %w", err)
		}
	}

	if u.Azure != nil {
		if u.Azure.TenantID, err = ResolveTemplate(u.Azure.TenantID, vars); err != nil {
			return fmt.Errorf("azure.tenantId: %w", err)
		}
	}

	return nil
}

// ResolveFragmentTemplates is a deprecated alias for ResolveClusterTemplates.
func ResolveFragmentTemplates(cl *Cluster, vars map[string]string) error {
	return ResolveClusterTemplates(cl, vars)
}
