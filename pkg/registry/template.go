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

	tmpl, err := template.New("fragment").Option("missingkey=error").Parse(text)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}
	return buf.String(), nil
}

// ResolveFragmentTemplates applies template variables to all string fields in a Fragment.
func ResolveFragmentTemplates(frag *Fragment, vars map[string]string) error {
	if len(vars) == 0 {
		return nil
	}

	var err error

	if frag.AWS != nil {
		if frag.AWS.Region, err = ResolveTemplate(frag.AWS.Region, vars); err != nil {
			return fmt.Errorf("aws.region: %w", err)
		}
		if frag.AWS.Cluster, err = ResolveTemplate(frag.AWS.Cluster, vars); err != nil {
			return fmt.Errorf("aws.cluster: %w", err)
		}
		if frag.AWS.Profile, err = ResolveTemplate(frag.AWS.Profile, vars); err != nil {
			return fmt.Errorf("aws.profile: %w", err)
		}
	}

	if frag.Azure != nil {
		if frag.Azure.SubscriptionID, err = ResolveTemplate(frag.Azure.SubscriptionID, vars); err != nil {
			return fmt.Errorf("azure.subscriptionId: %w", err)
		}
		if frag.Azure.ResourceGroup, err = ResolveTemplate(frag.Azure.ResourceGroup, vars); err != nil {
			return fmt.Errorf("azure.resourceGroup: %w", err)
		}
		if frag.Azure.Cluster, err = ResolveTemplate(frag.Azure.Cluster, vars); err != nil {
			return fmt.Errorf("azure.cluster: %w", err)
		}
		if frag.Azure.TenantID, err = ResolveTemplate(frag.Azure.TenantID, vars); err != nil {
			return fmt.Errorf("azure.tenantId: %w", err)
		}
	}

	if frag.Kubeconfig != "" {
		if frag.Kubeconfig, err = ResolveTemplate(frag.Kubeconfig, vars); err != nil {
			return fmt.Errorf("kubeconfig: %w", err)
		}
	}

	return nil
}
