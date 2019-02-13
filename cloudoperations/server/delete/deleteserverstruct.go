package deleteServer

import (
	cmn "github.com/nikhilsbhat/neuron/cloudoperations"
)

type DeleteServersInput struct {
	InstanceIds []string `json:"instanceids"`
	VpcId       string   `json:"vpcid"`
	Cloud       cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/delete
