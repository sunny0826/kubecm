package cmd

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

type DocsCommand struct {
	BaseCommand
}

var DOCS = "https://kubecm.cloud/"

// Init DocsCommand
func (dc *DocsCommand) Init() {
	dc.command = &cobra.Command{
		Use:   "docs",
		Short: "Open document website",
		Long:  "Open document website in your browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			return dc.runDocs(cmd, args)
		},
		Example: docsExample(),
	}
}

func (dc *DocsCommand) runDocs(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s#/en-us/cli/kubecm_%s", DOCS, cmd.Parent().Use)
	fmt.Printf("Opened %s in your browser.\n", url)
	return browser.OpenURL(url)
}

func docsExample() string {
	return `
# Open add command document page
kubecm add docs
`
}
