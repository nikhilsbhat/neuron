package commands

import (
	"github.com/spf13/cobra"
)

func init() {

	nwRegister("networkCreate", &cobra.Command{
		Use:          "create [flags]",
		Short:        "command to create a complete network",
		Long:         `This will help user to create network in cloud which he wants.`,
		RunE:         cc.createNetwork,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
	nwRegister("networkDelete", &cobra.Command{
		Use:          "delete [flags]",
		Short:        "command to delete a complete network or its components",
		Long:         `This will help user to delete network in cloud which he wants.`,
		RunE:         cc.deleteNetwork,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
	nwRegister("networkGet", &cobra.Command{
		Use:          "get [flags]",
		Short:        "command to get the details of network and its components",
		Long:         `This will help user to get network in cloud which he wants.`,
		RunE:         cc.getNetwork,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
	nwRegister("networkUpdate", &cobra.Command{
		Use:   "update [flags]",
		Short: "command to update the network and its components",
		Long: `This will help user to update network in cloud which he wants,
               by letting one to create various components in it.`,
		RunE:         cc.updateNetwork,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
	svRegister("serverCreate", &cobra.Command{
		Use:          "create [flags]",
		Short:        "command to create the instances",
		Long:         `This will help you to create/provision instance/server in cloud in the cloud you wish.`,
		RunE:         cc.createServer,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
	svRegister("serverDelete", &cobra.Command{
		Use:          "delete [flags]",
		Short:        "command to delete the instances",
		Long:         `This will help you to delete the servers from a particular network in cloud you whish.`,
		RunE:         cc.deleteServer,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
	svRegister("serverUpdate", &cobra.Command{
		Use:   "update [flags]",
		Short: "command to update the server",
		Long: `This will help you to update server by letting you to perform actions,
               such as start/stop etc, on server.`,
		RunE:         cc.updateServer,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
	svRegister("serverGet", &cobra.Command{
		Use:   "get [flags]",
		Short: "command to get the details of server",
		Long: `This will help you to get the details of server from the cloud you wish,
               make use of filtering capability here.`,
		RunE:         cc.getServer,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	})
	msRegister("regionsGet", &cobra.Command{
		Use:          "regions [flags]",
		Short:        "command to list regions",
		Long:         `This will help you to list available regions from the cloud you wish.`,
		RunE:         cc.getRegions,
		SilenceUsage: true,
	})
}
