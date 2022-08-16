package cmd

import (
	"fmt"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	v "github.com/sunny0826/kubecm/version"
)

// VersionCommand version cmd struct
type VersionCommand struct {
	BaseCommand
}

// version returns the version of kubecm.
type version struct {
	// kubecmVersion is a kubecm binary version.
	KubecmVersion string `json:"kubecmVersion"`
	// GitRevision is the commit of repo
	GitRevision string `json:"gitRevision"`
	// BuildDate is the build date of kubecm binary.
	BuildDate string `json:"buildDate"`
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
			fmt.Printf("%s: %s\n",
				ansi.Color("Version:", "blue"),
				ansi.Color(strings.TrimPrefix(getVersion().KubecmVersion+fmt.Sprintf("(%s)", getVersion().BuildDate), "v"), "white+h"))
			fmt.Printf("%s: %s\n",
				ansi.Color("GitRevision:", "blue"),
				ansi.Color(getVersion().GitRevision, "white+h"))
			fmt.Printf("%s: %s\n",
				ansi.Color("GoOs:", "blue"),
				ansi.Color(getVersion().GoOs, "white+h"))
			fmt.Printf("%s: %s\n",
				ansi.Color("GoArch:", "blue"),
				ansi.Color(getVersion().GoArch, "white+h"))
		},
	}
}

// getVersion returns version.
func getVersion() version {
	return version{
		v.Version,
		v.GitRevision,
		v.BuildDate,
		v.GoOs,
		v.GoArch,
	}
}
