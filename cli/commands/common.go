// Package that helps with common functionalities accross the package commands
package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (cm *cliMeta) getRegion(cmd *cobra.Command) string {
	reg, regrr := cmd.Flags().GetString("region")
	_ = regrr
	/*if regrr != nil {
	        cm.NeuronSaysItsError("region not passed for the cloud selected")
	}*/
	if reg != "" {
		return reg
	}
	cld, _ := cmd.Flags().GetString("cloud")
	return cm.getRegionFormConfig(cld)
}

func (cm *cliMeta) getCloud(cmd *cobra.Command) string {
	cld, clderr := cmd.Flags().GetString("cloud")
	_ = clderr
	/*if clderr != nil {
	        cm.NeuronSaysItsError("flag cloud is empty")
	}*/
	if cld != "" {
		return cld
	}
	return cm.getCloudFormConfig()
}

func (cm *cliMeta) getProfile(cmd *cobra.Command) string {
	prf, prferr := cmd.Flags().GetString("profile")
	if prferr != nil {
		cm.NeuronSaysItsError("flag profile not passed")
	}
	return prf
}

func (cm *cliMeta) getGetRaw(cmd *cobra.Command) bool {
	raw, rwerr := cmd.Flags().GetBool("filter")
	_ = rwerr
	/*if rwerr != nil {
	        cm.NeuronSaysItsError("flag getraw not used")
	}*/
	return !raw
}

func (cm *cliMeta) isAll(cmd *cobra.Command) bool {
	raw, rwerr := cmd.Flags().GetBool("all")
	if rwerr != nil {
		fmt.Println("flag all not used")
	}
	return raw
}

func (cm *cliMeta) getRegionFormConfig(cld string) string {
	for _, cloud := range cm.Cloud {
		if cld != "" {
			if cld == cloud.Name {
				return cloud.Region
			}
		}
		if cloud.Default == true {
			return cloud.Region
		}
	}
	return ""
}

func (cm *cliMeta) getCloudFormConfig() string {
	for _, cloud := range cm.Cloud {
		if cloud.Default == true {
			return cloud.Name
		}
	}
	cm.NeuronSaysItsError("you neither passed cloud, nor set a default cloud")
	return ""
}
