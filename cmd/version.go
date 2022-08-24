package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/savioxavier/termlink"

	"github.com/cli/safeexec"
	"github.com/sunny0826/kubecm/pkg/update"

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
			kubecmVersion := getVersion().KubecmVersion

			updateMessageChan := make(chan *update.ReleaseInfo)
			go func() {
				rel, _ := update.CheckForUpdate("sunny0826/kubecm", kubecmVersion)
				updateMessageChan <- rel
			}()
			fmt.Printf("%s: %s\n",
				ansi.Color("Version", "blue"),
				ansi.Color(strings.TrimPrefix(getVersion().KubecmVersion, "v"), "white+h"))
			fmt.Printf("%s: %s\n",
				ansi.Color("GoOs", "blue"),
				ansi.Color(getVersion().GoOs, "white+h"))
			fmt.Printf("%s: %s\n",
				ansi.Color("GoArch", "blue"),
				ansi.Color(getVersion().GoArch, "white+h"))
			newRelease := <-updateMessageChan
			if newRelease != nil {
				fmt.Printf("\n\n%s %s → %s\n",
					ansi.Color("A new release of kubecm is available:", "yellow"),
					ansi.Color(strings.TrimPrefix(kubecmVersion, "v"), "cyan"),
					ansi.Color(strings.TrimPrefix(newRelease.Version, "v"), "green"))
				if isUnderHomebrew() {
					fmt.Printf("To upgrade, run: %s\n", "brew update && brew upgrade kubecm")
				}
				fmt.Printf("%s\n\n",
					termlink.ColorLink("Click into the release page", newRelease.URL, "yellow"))
				//ansi.Color(newRelease.URL, "yellow"))
			}
		},
	}
}

// getVersion returns version.
func getVersion() version {
	return version{
		v.Version,
		v.GoOs,
		v.GoArch,
	}
}

// Check whether the gh binary was found under the Homebrew prefix
func isUnderHomebrew() bool {
	brewExe, err := safeexec.LookPath("brew")
	if err != nil {
		return false
	}

	brewPrefixBytes, err := exec.Command(brewExe, "--prefix").Output()
	if err != nil {
		return false
	}

	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return false
	}
	kubecmBinary, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	brewBinPrefix := filepath.Join(strings.TrimSpace(string(brewPrefixBytes)), "bin") + string(filepath.Separator)
	return strings.HasPrefix(kubecmBinary, brewBinPrefix)
}
