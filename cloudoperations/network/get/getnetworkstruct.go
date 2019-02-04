// The package which makes the tool cloud agnostic for fetching network details.
// The decision will be made here to route the request to respective package based on input.
package networkGet

import (
	cmn "neuron/cloudoperations"
)

// The struct which impliments method GetNetworks, GetSubnets.
type GetNetworksInput struct {
	// Ids or names of VPC's of which the information has to be fetched.
	VpcIds []string `json:"vpcids"`

	// Ids or names of the SUBNET's of which the informaion has to be fetched.
	SubnetIds []string `json:"subnetids"`

	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for network/get.
