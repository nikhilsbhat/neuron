// The package which makes the tool cloud agnostic for deleting network and its components.
// The decision will be made here to route the request to respective package based on input.
package networkDelete

import (
	cmn "neuron/cloudoperations"
)

// The struct which impliments method DeleteNetwork.
type DeleteNetworkInput struct {

	// Ids or names of VPC's which has to be deleted.
	VpcIds []string `json:"vpcids"`

	// Ids or names of SUBNET's which has to be deleted
	SubnetIds []string `json:"subnetids"`

	// Ids or name of Internet Gateways which has to be deleted.
	IgwIds []string `json:"igwids"`

	// Ids or name of Security Groups which has to be deletd.
	SecurityIds []string `json:"securityids"`

	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for network/create
