package serverCreate

import (
	cmn "github.com/nikhilsbhat/neuron/cloudoperations"
)

type ServerCreateInput struct {
	InstanceName string `json:"instancename"`
	Count        int64  `json:"count"`
	ImageId      string `json:"imageid"`
	SubnetId     string `json:"subnetid"`
	KeyName      string `json:"keyname"`
	Flavor       string `json:"flavor"`
	UserData     string `json:"userdata"`
	AssignPubIp  bool   `json:"assignpubip"`
	Cloud        cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/create
