package updateServers

import (
	cmn "github.com/nikhilsbhat/neuron/cloudoperations"
)

type UpdateServersInput struct {
	InstanceIds []string `json:"instanceids"`
	Action      string   `json:"action"`
	Cloud       cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/update
