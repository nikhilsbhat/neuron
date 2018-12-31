package serverCreate

type ServerCreateInput struct {
	InstanceName string `json:"InstanceName"`
	Count        int64  `json:"Count"`
	ImageId      string `json:"ImageId"`
	SubnetId     string `json:"SubnetId"`
	KeyName      string `json:"KeyName"`
	Flavor       string `json:"Flavor"`
	UserData     string `json:"UserData"`
	Cloud        string `json:"Cloud"`
	Region       string `json:"Region"`
	Profile      string `json:"Profile"`
	AssignPubIp  bool   `json:"AssignPubIp"`
	GetRaw       bool   `json:"GetRaw"`
}

//Nothing much from this file. This file contains only the structs for server/create
