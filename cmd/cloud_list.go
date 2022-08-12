package cmd

import (
	"errors"
	"fmt"

	"github.com/bndr/gotabulate"
	"github.com/spf13/cobra"
	"github.com/sunny0826/kubecm/pkg/cloud"
)

// CloudListCommand add command struct
type CloudListCommand struct {
	CloudCommand
}

// Init AddCommand
func (cl *CloudListCommand) Init() {
	cl.command = &cobra.Command{
		Use:   "list",
		Short: "list Cluster from cloud",
		Long:  "list Cluster from cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runCloudList(cmd, args)
		},
		Example: cloudListExample(),
	}
}

func (cl *CloudListCommand) runCloudList(cmd *cobra.Command, args []string) error {
	provider, err := cl.command.Flags().GetString("provider")
	regionID, err := cl.command.Flags().GetString("region_id")
	var num int
	if provider == "" {
		num = selectCloud(Clouds, "Select Cloud")
	} else {
		num = checkFlags(provider)
	}
	clusters, err := getClusters(provider, regionID, num)
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		return errors.New("no clusters found")
	}
	return printListTable(clusters)
}

// PrintTable generate table
func printListTable(clusters []cloud.ClusterInfo) error {
	var table [][]string
	for _, k := range clusters {
		conTmp := []string{k.ID, k.Name, k.RegionID, k.K8sVersion, k.ConsoleURL}
		table = append(table, conTmp)
	}

	if table != nil {
		tabulate := gotabulate.Create(table)
		tabulate.SetHeaders([]string{"ID", "NAME", "REGION ID", "VERSION", "CONSOLE URL"})
		// Turn On String Wrapping
		tabulate.SetWrapStrings(false)
		// Render the table
		tabulate.SetAlign("center")
		fmt.Println(tabulate.Render("grid", "left"))
	} else {
		return errors.New("context not found")
	}
	return nil
}

func cloudListExample() string {
	return `
# Supports Ali Cloud and Tencent Cloud
# The AK/AS of the cloud platform will be retrieved directly 
# if it exists in the environment variable, 
# otherwise a prompt box will appear asking for it.

# Set env AliCloud secret key
export ACCESS_KEY_ID=xxx
export ACCESS_KEY_SECRET=xxx
# Set env Tencent secret key
export TENCENTCLOUD_SECRET_ID=xxx
export TENCENTCLOUD_SECRET_KEY=xxx
# Set env Rancher secret key
export RANCHER_SERVER_URL=https://xxx
export RANCHER_API_KEY=xxx
# Interaction: list kubeconfig from cloud
kubecm cloud list
# Add kubeconfig from cloud
kubecm cloud list --provider alibabacloud --cluster_id=xxxxxx
`
}
