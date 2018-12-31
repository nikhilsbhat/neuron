package deleteServer

type DeleteServersInput struct {
	InstanceIds []string `json:"InstanceIds"`
	VpcId       string   `json:"VpcId"`
	Cloud       string   `json:"Cloud"`
	Region      string   `json:"Region"`
	Profile     string   `json:"Profile"`
	GetRaw      bool     `json:"GetRaw"`
}

//Nothing much from this file. This file contains only the structs for server/delete
