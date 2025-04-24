package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// RangeCommand range subcommand for delete command
type RangeCommand struct {
	BaseCommand
	matchMode string // support prefix, suffix, contains
	yes       bool   // skip confirmation prompt
}

// Init RangeCommand
func (rc *RangeCommand) Init() {
	rc.command = &cobra.Command{
		Use:   "range",
		Short: "Delete contexts matching a pattern",
		Long:  `Delete all contexts that match a specified pattern from the kubeconfig`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rc.runRange(cmd, args)
		},
		Example: rangeExample(),
	}

	rc.command.Flags().StringVarP(&rc.matchMode, "mode", "", "prefix", "Match mode: prefix, suffix, or contains")
	rc.command.Flags().BoolVarP(&rc.yes, "yes", "y", false, "Skip confirmation prompt")
	rc.AddCommands(&DocsCommand{})
}

func (rc *RangeCommand) runRange(command *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("no pattern specified")
	}

	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig file %q: %w", cfgFile, err)
	}

	// Select contexts to delete
	needDeleteContexts, err := matchContexts(config.Contexts, args[0], rc.matchMode)
	if err != nil {
		return err
	}

	if len(needDeleteContexts) == 0 {
		return errors.New("no contexts matched the specified pattern")
	}

	// Confirm delete
	fmt.Printf("Found %d contexts matching %s mode with pattern %q:\n", len(needDeleteContexts), rc.matchMode, args[0])
	for _, ctx := range needDeleteContexts {
		fmt.Printf("  - %s\n", ctx)
	}

	// Skip confirmation if -y/--yes flag is set
	if !rc.yes {
		if !strings.EqualFold(BoolUI(fmt.Sprintf("Are you sure you want to delete these %d contexts?", len(needDeleteContexts))), "True") {
			return errors.New("range delete operation cancelled")
		}
	}

	if err := rangeDeleteContexts(needDeleteContexts, config); err != nil {
		return err
	}

	if err := WriteConfig(true, cfgFile, config); err != nil {
		return fmt.Errorf("failed to write kubeconfig file %q: %w", cfgFile, err)
	}

	return nil
}

// matchContexts selects contexts that match the given pattern and mode.
func matchContexts(contexts map[string]*clientcmdapi.Context, pattern, matchMode string) ([]string, error) {
	if pattern == "" {
		return nil, errors.New("pattern cannot be empty")
	}

	validModes := map[string]bool{"prefix": true, "suffix": true, "contains": true}
	if !validModes[matchMode] {
		return nil, fmt.Errorf("invalid match mode: %s, must be one of: prefix, suffix, contains", matchMode)
	}

	var matches []string
	for contextName := range contexts {
		var matched bool
		switch matchMode {
		case "prefix":
			matched = strings.HasPrefix(contextName, pattern)
		case "suffix":
			matched = strings.HasSuffix(contextName, pattern)
		case "contains":
			matched = strings.Contains(contextName, pattern)
		}
		if matched {
			matches = append(matches, contextName)
		}
	}

	return matches, nil
}

// rangeDeleteContexts deletes the specified contexts and their associated clusters and auth infos if not used elsewhere.
func rangeDeleteContexts(needDeleteContexts []string, config *clientcmdapi.Config) error {
	for _, ctx := range needDeleteContexts {
		if _, exists := config.Contexts[ctx]; !exists {
			return fmt.Errorf("context %q does not exist", ctx)
		}

		delContext := config.Contexts[ctx]
		isClusterNameExist, isUserNameExist := checkClusterAndUserNameExceptContextToDelete(config, delContext)

		if !isUserNameExist {
			delete(config.AuthInfos, delContext.AuthInfo)
		}
		if !isClusterNameExist {
			delete(config.Clusters, delContext.Cluster)
		}
		delete(config.Contexts, ctx)

		fmt.Printf("Context Delete:「%s」\n", ctx)
	}

	return nil
}

func rangeExample() string {
	return `
# Delete all contexts with prefix "dev-"
kubecm delete range dev-
or 
kubecm delete range --mode prefix dev-

# Delete all contexts with suffix "prod"
kubecm delete range --mode suffix prod

# Delete all contexts containing "staging"
kubecm delete range --mode contains staging

# Force delete all contexts with prefix "dev-" (skip confirmation)
kubecm delete range dev- -y
`
}
