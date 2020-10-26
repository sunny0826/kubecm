package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// AliasCommand alias command struct
type AliasCommand struct {
	BaseCommand
}

// Init AliasCommand
func (al *AliasCommand) Init() {
	al.command = &cobra.Command{
		Use:     "alias",
		Short:   "Generate alias for all contexts",
		Long:    "Generate alias for all contexts",
		Aliases: []string{"al"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return al.runAlias(cmd, args)
		},
		Example: aliasExample(),
	}
	al.command.DisableFlagsInUseLine = true
	al.command.Flags().StringP("out", "o", "", "output to ~/.zshrc or ~/.bash_profile")
}

func (al *AliasCommand) runAlias(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	allTemp := `
## KubeCm Alias Start%s
## KubeCm Alias End
`
	aliasTemp := `
# %s
alias %s='kubectl --context %s'`
	var tmp string
	for key := range config.Contexts {
		tmp += fmt.Sprintf(aliasTemp, key, "k-"+key, key)
	}
	output, _ := al.command.Flags().GetString("out")
	result := fmt.Sprintf(allTemp, tmp)
	switch output {
	case "bash":
		err = writeAppend(result, filepath.Join(homeDir(), ".bash_profile"))
		if err != nil {
			return err
		}
		al.command.Println("「.bash_profile」 write successful!")
	case "zsh":
		err = writeAppend(result, filepath.Join(homeDir(), ".zshrc"))
		if err != nil {
			return err
		}
		al.command.Println("「.zshrc」 write successful!")
	default:
		al.command.Print(result)
	}
	return nil
}

func writeAppend(context, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	write := bufio.NewWriter(f)
	//strings.TrimSpace(context)
	_, _ = write.WriteString(context)
	err = write.Flush()
	if err != nil {
		return err
	}
	return nil
}

func aliasExample() string {
	return `
# dev 
alias k-dev='kubectl --context dev'
# test
alias k-test='kubectl --context test'
# prod
alias k-prod='kubectl --context prod'
$ kubecm alias -o zsh
# add context to ~/.zshrc
$ kubecm alias -o bash
# add context to ~/.bash_profile
`
}
