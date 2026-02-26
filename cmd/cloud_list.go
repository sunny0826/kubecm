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
		Short: "List Cluster from cloud",
		Long:  "List Cluster from cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runCloudList(cmd, args)
		},
		Example: cloudListExample(),
	}
}

func (cl *CloudListCommand) runCloudList(cmd *cobra.Command, args []string) error {
	provider, _ := cl.command.Flags().GetString("provider")
	regionID, _ := cl.command.Flags().GetString("region_id")
	awsProfile, _ := cl.command.Flags().GetString("aws_profile")
	var num int
	if provider == "" {
		num = selectCloud(Clouds, "Select Cloud")
	} else {
		num = checkFlags(provider)
	}
	clusters, err := getClusters(provider, regionID, awsProfile, num)
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
		conTmp := []string{k.ID, k.Account, k.Name, k.RegionID, k.K8sVersion, k.ConsoleURL}
		table = append(table, conTmp)
	}

	if table != nil {
		tabulate := gotabulate.Create(table)
		tabulate.SetHeaders([]string{"ID", "ACCOUNT", "NAME", "REGION ID", "VERSION", "CONSOLE URL"})
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
# Supports AlibabaCloud, Tencent Cloud, Rancher, AWS and Azure

# Interaction: list clusters from cloud
kubecm cloud list

# AlibabaCloud
kubecm cloud list --provider alibabacloud

# AWS with named profile
kubecm cloud list --provider aws --aws_profile my-profile --region_id eu-west-3

# AWS with default credential chain
kubecm cloud list --provider aws --region_id us-east-1

# Azure
kubecm cloud list --provider azure
`
}
