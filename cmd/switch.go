package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// SwitchCommand switch cmd struct
type SwitchCommand struct {
	BaseCommand
}

// Init SwitchCommand
func (sc *SwitchCommand) Init() {
	sc.command = &cobra.Command{
		Use:   "switch",
		Short: "Switch Kube Context interactively",
		Long: `
Switch Kube Context interactively
`,
		Aliases: []string{"s"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return sc.runSwitch(cmd, args)
		},
		Example: switchExample(),
	}
}

func (sc *SwitchCommand) runSwitch(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	var kubeItems []Needle
	current := config.CurrentContext
	for key, obj := range config.Contexts {
		if key != current {
			kubeItems = append(kubeItems, Needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
		} else {
			kubeItems = append([]Needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
		}
	}
	// exit option
	kubeItems, err = ExitOption(kubeItems)
	if err != nil {
		return err
	}
	num := SelectUI(kubeItems, "Select Kube Context")
	kubeName := kubeItems[num].Name
	config.CurrentContext = kubeName
	err = WriteConfig(true, cfgFile, config)
	if err != nil {
		return err
	}
	sc.command.Printf("Switched to context 「%s」\n", config.CurrentContext)
	err = Formatable()
	if err != nil {
		return err
	}
	return nil
}

func switchExample() string {
	return `
# Switch Kube Context interactively
kubecm switch
`
}
