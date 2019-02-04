package getServers

import (
	cmn "neuron/cloudoperations"
)

type GetServersInput struct {
	InstanceIds []string `json:"instanceids"`
	VpcIds      []string `json:"vpcids"`
	SubnetIds   []string `json:"subnetids"`
	Cloud       cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/get
