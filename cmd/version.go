package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
)

// VersionCommand version cmd struct
type VersionCommand struct {
	BaseCommand
}

var (
	kubecmVersion = "unknown"
	goos          = runtime.GOOS
	goarch        = runtime.GOARCH
)

// version returns the version of kustomize.
type version struct {
	// KustomizeVersion is a kustomize binary version.
	KubecmVersion string `json:"kubecmVersion"`
	// GoOs holds OS name.
	GoOs string `json:"goOs"`
	// GoArch holds architecture name.
	GoArch string `json:"goArch"`
}

// Init VersionCommand
func (vc *VersionCommand) Init() {
	vc.command = &cobra.Command{
		Use:     "version",
		Short:   "Print version info",
		Long:    "Print version info",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version: %s\n", getVersion().KubecmVersion)
			cmd.Printf("GoOs: %s\n", getVersion().GoOs)
			cmd.Printf("GoArch: %s\n", getVersion().GoArch)
		},
	}
}

// getVersion returns version.
func getVersion() version {
	return version{
		kubecmVersion,
		goos,
		goarch,
	}
}
