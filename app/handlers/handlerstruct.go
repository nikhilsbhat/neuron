package handlers

import (
	"github.com/globalsign/mgo/bson"
	"time"
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

type FillStructs struct {
	Data interface{} `json:"Data,omitempty"`
	Type interface{} `json:"Type,omitempty"`
}

type uiTemp struct {
	Title string        `json:"title,omitempty"`
	Cont  string        `json:"string,omitempty"`
	Pass  []interface{} `json:"pass,omitempty"`
}

type Error struct {
	Error string
}

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
