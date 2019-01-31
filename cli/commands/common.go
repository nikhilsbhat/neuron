// Package that helps with common functionalities accross the package commands
package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func getRegion(cmd *cobra.Command) string {
	reg, regrr := cmd.Flags().GetString("region")
	if regrr != nil {
		fmt.Println("region not passed for the cloud selected")
	}
	return reg
}

func getCloud(cmd *cobra.Command) string {
	cld, clderr := cmd.Flags().GetString("cloud")
	if clderr != nil {
		fmt.Println("flag cloud is empty")
	}
	return cld
}

func getProfile(cmd *cobra.Command) string {
	prf, prferr := cmd.Flags().GetString("profile")
	if prferr != nil {
		fmt.Println("flag profile not passed")
	}
	return prf
}

func getGetRaw(cmd *cobra.Command) bool {
	raw, rwerr := cmd.Flags().GetBool("filter")
	if rwerr != nil {
		fmt.Println("flag filter not used")
	}
	return !raw
}

func isAll(cmd *cobra.Command) bool {
	raw, rwerr := cmd.Flags().GetBool("all")
	if rwerr != nil {
		fmt.Println("flag all not used")
	}
	return raw
}
