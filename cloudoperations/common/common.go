package commonOperations

import (
	db "neuron/database"
	dbcommon "neuron/database/common"
)

const (
	DefaultOpResponse    = "We have not reached to openstack yet"
	DefaultAzResponse    = "We have not reached to azure yet"
	DefaultGcpResponse   = "We have not reached to google cloud yet"
	DefaultCloudResponse = "I feel we are lost in performing the action, guess you have entered wrong cloud. The action was: "
)

type GetCredentialsInput struct {
	Profile string
	Cloud   string
}

func GetCredentials(gcred *GetCredentialsInput) (db.CloudProfiles, error) {

	//fetchinig credentials from loged-in user to establish the connection with appropriate cloud.
	creds, crderr := dbcommon.GetCloudCredentails(
		db.UserData{UserName: "nikhibt434@gmail", Password: "42bhat24"},
		db.GetCloudAccess{ProfileName: gcred.Profile, Cloud: gcred.Cloud},
		db.DataDetail{"neuron", "users"},
	)
	if crderr != nil {
		return db.CloudProfiles{}, crderr
	}

	return creds, nil
}