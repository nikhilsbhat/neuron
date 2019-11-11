package handlers

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type servercreateinput struct {
	InstanceName string
	ImageId      string
	SubnetId     string
	ConsulConfig interface{}
	KeyName      string
	Flavor       string
	UserData     string
	AssignPubIp  bool
	Cloud        cloudConfig
}

type startimageinput struct {
	AppVersion   string
	AppName      string
	RepoEmail    string
	RepoUsername string
	RepoPasswd   string
	ArtDomain    string
	ArtUsername  string
	ArtPasswd    string
	InstanceName string
	SubnetId     string
	KeyName      string
	Flavor       string
	Cloud        cloudConfig
	AssignPubIp  bool
}

type cloudConfig struct {
	Cloud  string
	Region string
}

type createbuildmachine struct {
	AppVersion string
	UniqueId   string
}

type getserverdetails struct {
	SubnetId string
	Cloud    cloudConfig
}

// Package FillStructs helps in filling structs for neuron UI
type FillStructs struct {
	Data interface{} `json:"Data,omitempty"`
	Type interface{} `json:"Type,omitempty"`
}

type uiTemp struct {
	Title string        `json:"title,omitempty"`
	Cont  string        `json:"string,omitempty"`
	Pass  []interface{} `json:"pass,omitempty"`
}

// Package Error holds the error message for neuron UI
type Error struct {
	Error string
}

// Package CIData maps the struct to neuron UI
type CIData struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id"`
	//Id         int       `bson:"_id,omitempty" json:"id"`
	CiName     string    `json:"CiName" bson:"ciname"`
	CiURL      string    `json:"CiURL" bson:"ciurl"`
	CiUsername string    `json:"CiUsername" bson:"ciusername"`
	CiPassword string    `json:"CiPassword" bson:"cipassword"`
	Timestamp  time.Time `json:"Timestamp" bson:"timestamp"`
}

type datadetail struct {
	Database   string
	Collection string
}
