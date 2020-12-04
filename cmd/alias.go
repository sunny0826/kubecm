package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	_ = al.command.MarkFlagRequired("out")
}

const SourceCmd = "[[ ! -f ~/.kubecm ]] || source ~/.kubecm"

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
	err = updateFile(result, filepath.Join(homeDir(), ".kubecm"))
	if err != nil {
		return err
	}
	switch output {
	case "bash":
		err = writeAppend(SourceCmd, filepath.Join(homeDir(), ".bash_profile"))
		if err != nil {
			return err
		}
		al.command.Println("「.bash_profile」 write successful!\nPlease enter command `source .bash_profile`")
	case "zsh":
		err = writeAppend(SourceCmd, filepath.Join(homeDir(), ".zshrc"))
		if err != nil {
			return err
		}
		al.command.Println("「.zshrc」 write successful!\nPlease enter command `source .zshrc`")
	default:
		al.command.PrintErrln("Parameter error! Please input bash or zsh")
	}

	return nil
}

func updateFile(cxt, path string) error {
	err := ioutil.WriteFile(path, []byte(cxt), 0644)
	if err != nil {
		return err
	}
	return nil
}

func writeAppend(context, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	var exist bool
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if ok := strings.EqualFold(SourceCmd, line); ok {
			exist = true
			break
		}
	}
	if !exist {
		write := bufio.NewWriter(f)
		_, _ = write.WriteString(context+"\n")
		err = write.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func aliasExample() string {
	return `
$ kubecm alias -o zsh
# add context to ~/.zshrc
$ kubecm alias -o bash
# add context to ~/.bash_profile
`
}
