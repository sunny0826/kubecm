package registry

import (
	"testing"
)

func TestResolveTemplate(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		vars    map[string]string
		want    string
		wantErr bool
	}{
		{
			name: "simple substitution",
			text: "hello {{ .Username }}",
			vars: map[string]string{"Username": "clark"},
			want: "hello clark",
		},
		{
			name: "no variables",
			text: "hello world",
			vars: nil,
			want: "hello world",
		},
		{
			name: "multiple variables",
			text: "{{ .Env }}-{{ .Region }}",
			vars: map[string]string{"Env": "prod", "Region": "eu-central-1"},
			want: "prod-eu-central-1",
		},
		{
			name:    "missing variable",
			text:    "hello {{ .Missing }}",
			vars:    map[string]string{"Username": "clark"},
			wantErr: true,
		},
		{
			name: "empty vars map",
			text: "no template here",
			vars: map[string]string{},
			want: "no template here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveTemplate(tt.text, tt.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ResolveTemplate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveFragmentTemplates(t *testing.T) {
	t.Run("aws fragment", func(t *testing.T) {
		frag := &Fragment{
			Provider: "aws",
			AWS: &AWSFragment{
				Region:  "{{ .Region }}",
				Cluster: "eks-{{ .Env }}",
				Profile: "{{ .Profile }}",
			},
		}
		vars := map[string]string{
			"Region":  "eu-central-1",
			"Env":     "prod",
			"Profile": "admin",
		}

		if err := ResolveFragmentTemplates(frag, vars); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if frag.AWS.Region != "eu-central-1" {
			t.Errorf("region = %q, want %q", frag.AWS.Region, "eu-central-1")
		}
		if frag.AWS.Cluster != "eks-prod" {
			t.Errorf("cluster = %q, want %q", frag.AWS.Cluster, "eks-prod")
		}
		if frag.AWS.Profile != "admin" {
			t.Errorf("profile = %q, want %q", frag.AWS.Profile, "admin")
		}
	})

	t.Run("static fragment with kubeconfig template", func(t *testing.T) {
		frag := &Fragment{
			Provider:   "static",
			Kubeconfig: "token: {{ .Username }}-token",
		}
		vars := map[string]string{"Username": "clark"}

		if err := ResolveFragmentTemplates(frag, vars); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if frag.Kubeconfig != "token: clark-token" {
			t.Errorf("kubeconfig = %q, want %q", frag.Kubeconfig, "token: clark-token")
		}
	})

	t.Run("nil vars is no-op", func(t *testing.T) {
		frag := &Fragment{
			Provider: "aws",
			AWS: &AWSFragment{
				Region:  "eu-west-1",
				Cluster: "test",
			},
		}
		if err := ResolveFragmentTemplates(frag, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if frag.AWS.Region != "eu-west-1" {
			t.Errorf("region changed unexpectedly")
		}
	})
}
