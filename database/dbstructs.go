package database

import (
	//"github.com/globalsign/mgo/bson"
	//"github.com/globalsign/mgo"
	"time"
)

var (
	Db interface{}
)

type Storage struct {
	Db interface{} `json:"Db,omitempty"`
	Fs string      `json:"Fs,omitempty"`
}

type DataDetail struct {
	Database   string
	Collection string
}

type UserData struct {
	Id            int             `bson:"_id,omitempty" json:"id"`
	UserName      string          `json:"UserName" bson:"username"`
	Password      string          `json:"Password" bson:"password"`
	CloudProfiles []CloudProfiles `json:"CloudProfiles" bson:"cloudprofiles"`
}

type CloudProfiles struct {
	Name           string    `json:"Name" bson:"name,omitempty"`
	Cloud          string    `json:"Cloud" bson:"cloud,omitempty"`
	KeyId          string    `json:"KeyId" bson:"keyid,omitempty"`
	SecretAccess   string    `json:"SecretAccess" bson:"secretaccess,omitempty"`
	ClientId       string    `json:"ClientID" bson:"clientid,omitempty"`
	SubscriptionId string    `json:"SubscriptionID" bson:"subscriptionid,omitempty"`
	TenantId       string    `json:"TenantID" bson:"tenantid,omitempty"`
	ClientSecret   string    `json:"ClientSecret" bson:"clientsecret,omitempty"`
	CreationTime   time.Time `json:"CreationTime" bson:"creationtime,omitempty"`
}

type CiData struct {
	Id         int       `json:"id" bson:"_id,omitempty"`
	CiName     string    `json:"CiName" bson:"ciname"`
	CiURL      string    `json:"CiURL" bson:"ciurl"`
	CiUsername string    `json:"CiUsername" bson:"ciusername"`
	CiPassword string    `json:"CiPassword" bson:"cipassword"`
	Timestamp  time.Time `json:"Timestamp" bson:"timestamp"`
}

type GetCloudAccess struct {
	ProfileName string
	Cloud       string
}
