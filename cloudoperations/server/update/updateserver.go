package updateServers

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	server "neuron/cloud/aws/operations/server"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type UpdateServersResponse struct {
	AwsResponse     []server.ServerResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                  `json:"AzureResponse,omitempty"`
	DefaultResponse string                  `json:"Response,omitempty"`
}

// being update_servers my job is to update servers (start/stop, change ebs etc) and give back the response who called me
func (serv *UpdateServersInput) UpdateServers() (UpdateServersResponse, error) {

	switch strings.ToLower(serv.Cloud) {
	case "aws":

		creds, crederr := common.GetCredentials(&common.GetCredentialsInput{Profile: serv.Profile, Cloud: serv.Cloud})
		if crederr != nil {
			return UpdateServersResponse{}, crederr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: serv.Region, Resource: "ec2", Session: sess}

		// I will call UpdateServer of interface and get the things done
		serverin := server.UpdateServerInput{InstanceIds: serv.InstanceIds, Action: serv.Action, GetRaw: serv.GetRaw}
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
