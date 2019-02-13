package getloadbalancer

import (
	cmn "github.com/nikhilsbhat/neuron/cloudoperations"
)

type GetLoadbalancerInput struct {
	LbNames []string `json:"lbnames"`
	LbArns  []string `json:"lbarns"`
	Type    string   `json:"type"`
	Cloud   cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for loadbalance/get
