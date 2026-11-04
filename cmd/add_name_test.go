package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestAddCommandContextNaming(t *testing.T) {
	const originalName = "kubernetes-admin@kubernetes"
	tests := []struct {
		name     string
		args     []string
		multiple bool
		want     []string
	}{
		{
			name: "override name",
			args: []string{"--context-name", "dev-85"},
			want: []string{"dev-85"},
		},
		{
			name: "override ignores prefix and template",
			args: []string{"--context-name", "dev-85", "--context-prefix", "ignored", "--context-template", "invalid"},
			want: []string{"dev-85"},
		},
		{
			name: "override preserves trailing hyphen",
			args: []string{"--context-name", "dev-85-"},
			want: []string{"dev-85-"},
		},
		{
			name:     "override disambiguates multiple contexts",
			args:     []string{"--context-name", "dev-85"},
			multiple: true,
			want:     []string{"dev-85", "dev-85-2"},
		},
		{
			name: "filter uses overridden name",
			args: []string{"--context-name", "dev-85", "--context", "dev-85"},
			want: []string{"dev-85"},
		},
		{
			name: "default preserves original name",
			want: []string{originalName},
		},
		{
			name: "empty template preserves original name",
			args: []string{"--context-template="},
			want: []string{originalName},
		},
		{
			name: "prefix with empty template",
			args: []string{"--context-prefix", "dev-85", "--context-template="},
			want: []string{"dev-85-" + originalName},
		},
		{
			name: "prefix with custom template",
			args: []string{"--context-prefix", "dev-85", "--context-template", "cluster,user"},
			want: []string{"dev-85-kubernetes-kubernetes-admin"},
		},
	}

	for _, source := range []string{"file", "stdin"} {
		for _, tt := range tests {
			t.Run(source+"/"+tt.name, func(t *testing.T) {
				oldCfgFile, oldSilenceTable, oldCfgCreate := cfgFile, silenceTable, cfgCreate
				t.Cleanup(func() {
					cfgFile, silenceTable, cfgCreate = oldCfgFile, oldSilenceTable, oldCfgCreate
				})
				cfgFile = filepath.Join(t.TempDir(), "target.yaml")
				silenceTable, cfgCreate = true, false
				require.NoError(t, clientcmd.WriteToFile(*clientcmdapi.NewConfig(), cfgFile))

				input := clientcmdapi.NewConfig()
				input.Clusters["kubernetes"] = &clientcmdapi.Cluster{Server: "https://example.invalid:6443"}
				input.AuthInfos["kubernetes-admin"] = &clientcmdapi.AuthInfo{}
				input.Contexts[originalName] = &clientcmdapi.Context{
					Cluster: "kubernetes", AuthInfo: "kubernetes-admin", Namespace: "default",
				}
				input.CurrentContext = originalName
				if tt.multiple {
					input.Contexts["second"] = &clientcmdapi.Context{
						Cluster: "kubernetes", AuthInfo: "kubernetes-admin", Namespace: "second",
					}
				}
				inputPath := filepath.Join(t.TempDir(), "input.yaml")
				require.NoError(t, clientcmd.WriteToFile(*input, inputPath))

				fileArg := inputPath
				if source == "stdin" {
					stdin, err := os.Open(inputPath)
					require.NoError(t, err)
					oldStdin := os.Stdin
					os.Stdin = stdin
					t.Cleanup(func() {
						os.Stdin = oldStdin
						require.NoError(t, stdin.Close())
					})
					fileArg = "-"
				}

				command := &AddCommand{}
				command.Init()
				command.command.SetArgs(append([]string{"-c", "-f", fileArg}, tt.args...))
				require.NoError(t, command.command.Execute())

				got, err := clientcmd.LoadFromFile(cfgFile)
				require.NoError(t, err)
				var names, namespaces []string
				for name, ctx := range got.Contexts {
					names = append(names, name)
					namespaces = append(namespaces, ctx.Namespace)
					require.Equal(t, "https://example.invalid:6443", got.Clusters[ctx.Cluster].Server)
					require.Contains(t, got.AuthInfos, ctx.AuthInfo)
				}
				sort.Strings(names)
				require.Equal(t, tt.want, names)
				if tt.multiple {
					require.ElementsMatch(t, []string{"default", "second"}, namespaces)
				} else {
					require.Equal(t, tt.want[0], got.CurrentContext)
				}
			})
		}
	}
}
