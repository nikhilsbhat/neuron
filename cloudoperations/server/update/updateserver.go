package updateServers

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	server "neuron/cloud/aws/operations/server"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	support "neuron/cloudoperations/support"
	log "neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type UpdateServersResponse struct {

	// Contains filtered/unfiltered response of AWS.
	AwsResponse []server.ServerResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// Being GetServersDetails, job of him is to update servers (start/stop, change ebs etc)
//  with the instructions passed to him and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (serv *UpdateServersInput) UpdateServers() (UpdateServersResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return UpdateServersResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":

		creds, crederr := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: serv.Cloud.Profile,
				Cloud:   serv.Cloud.Name,
			},
		)

		if crederr != nil {
			return UpdateServersResponse{}, crederr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call UpdateServer of interface and get the things done
		serverin := server.UpdateServerInput{InstanceIds: serv.InstanceIds, Action: serv.Action, GetRaw: serv.Cloud.GetRaw}
		response, err := serverin.UpdateServer(authinpt)
		if err != nil {
			return UpdateServersResponse{}, err
		}
		return UpdateServersResponse{AwsResponse: response}, nil

	case "azure":
		return UpdateServersResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return UpdateServersResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return UpdateServersResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:
		log.Info("")
		log.Error("I feel we are lost in updating servers :S")
		log.Info("")
		return UpdateServersResponse{}, fmt.Errorf(common.DefaultCloudResponse + "UpdateServers")
	}
}

func New() *UpdateServersInput {
	net := &UpdateServersInput{}
	return net
}
