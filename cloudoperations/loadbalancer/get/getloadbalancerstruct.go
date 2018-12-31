package getloadbalancer

type GetLoadbalancerInput struct {
	LbNames []string
	LbArns  []string
	Type    string
	Cloud   string
	Region  string
	Profile string
	GetRaw  bool
}

//Nothing much from this file. This file contains only the structs for loadbalance/get
