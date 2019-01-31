// Package that helps with common functionalities accross the package commands
package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (cm *cliMeta) getRegion(cmd *cobra.Command) string {
	reg, regrr := cmd.Flags().GetString("region")
	if regrr != nil {
		cm.NeuronSaysItsError("region not passed for the cloud selected")
	}
	return reg
}

func (cm *cliMeta) getCloud(cmd *cobra.Command) string {
	cld, clderr := cmd.Flags().GetString("cloud")
	if clderr != nil {
		cm.NeuronSaysItsError("flag cloud is empty")
	}
	return cld
}

func (cm *cliMeta) getProfile(cmd *cobra.Command) string {
	prf, prferr := cmd.Flags().GetString("profile")
	if prferr != nil {
		cm.NeuronSaysItsError("flag profile not passed")
	}
	return prf
}

func (cm *cliMeta) getGetRaw(cmd *cobra.Command) bool {
	raw, rwerr := cmd.Flags().GetBool("getraw")
	_ = rwerr
	/*if rwerr != nil {
		cm.NeuronSaysItsError("flag getraw not used")
	}*/
	return raw
}

func (cm *cliMeta) isAll(cmd *cobra.Command) bool {
	raw, rwerr := cmd.Flags().GetBool("all")
	if rwerr != nil {
		fmt.Println("flag all not used")
	}
	return raw
}
