package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// RangeDeleteCommand range delete cmd struct
type RangeDeleteCommand struct {
	BaseCommand
	matchMode string // support prefix, suffix, contains
}

// Init RangeDeleteCommand
func (rdc *RangeDeleteCommand) Init() {
	rdc.command = &cobra.Command{
		Use:     "range-delete",
		Short:   "Delete contexts matching a pattern",
		Long:    `Delete all contexts that match a specified pattern from the kubeconfig`,
		Aliases: []string{"rd"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return rdc.runRangeDelete(cmd, args)
		},
		Example: rangeDeleteExample(),
	}

	rdc.command.Flags().StringVarP(&rdc.matchMode, "mode", "", "prefix", "Match mode: prefix, suffix, or contains")
	rdc.AddCommands(&DocsCommand{})
}

func (rdc *RangeDeleteCommand) runRangeDelete(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("nothing can be deleted, because no pattern specified")
	}

	if rdc.matchMode != "prefix" && rdc.matchMode != "suffix" && rdc.matchMode != "contains" {
		return fmt.Errorf("invalid match mode: %s, must be one of: prefix, suffix, contains", rdc.matchMode)
	}

	err = rangeDeleteContexts(args[0], rdc.matchMode, config)
	if err != nil {
		return err
	}

	err = WriteConfig(false, cfgFile, config)
	if err != nil {
		return err
	}
	return nil
}

func rangeDeleteContexts(pattern string, matchMode string, config *clientcmdapi.Config) error {
	var needDeleteContexts []string

	for contextName := range config.Contexts {
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
			needDeleteContexts = append(needDeleteContexts, contextName)
		}
	}

	if len(needDeleteContexts) == 0 {
		return errors.New("nothing can be deleted, because no contexts matched")
	}

	// confirm delete
	fmt.Printf("Found %d contexts matching %s mode with pattern %q:\n", len(needDeleteContexts), matchMode, pattern)
	for _, ctx := range needDeleteContexts {
		fmt.Printf("  - %s\n", ctx)
	}

	confirm := BoolUI(fmt.Sprintf("Are you sure you want to delete these %d contexts?", len(needDeleteContexts)))
	if confirm != "True" {
		return errors.New("range delete operation cancelled")
	}

	// delete all matched contexts
	for _, ctx := range needDeleteContexts {
		delContext := config.Contexts[ctx]
		isClusterNameExist, isUserNameExist := checkClusterAndUserNameExceptContextToDelete(config, config.Contexts[ctx])
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

func rangeDeleteExample() string {
	return `
# Delete all contexts with prefix "dev-"
kubecm range-delete dev-
or 
kubecm range-delete -m prefix dev-

# Delete all contexts with suffix "-prod"
kubecm range-delete -m suffix -prod

# Delete all contexts containing "staging"
kubecm range-delete -m contains staging
`
}
