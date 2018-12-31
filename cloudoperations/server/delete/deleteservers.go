package deleteServer

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	server "neuron/cloud/aws/operations/server"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type DeleteServerResponse struct {
	AwsResponse     []server.ServerResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                  `json:"AzureResponse,omitempty"`
	DefaultResponse string                  `json:"Response,omitempty"`
}

// being delete server, my job is to delete servers/vms and give back the response who called me
func (serv *DeleteServersInput) DeleteServer() (DeleteServerResponse, error) {

	switch strings.ToLower(serv.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: serv.Profile, Cloud: serv.Cloud})
		if err != nil {
			return DeleteServerResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorize
		authinpt := auth.EstablishConnectionInput{Region: serv.Region, Resource: "ec2", Session: sess}

		// I will call DeleteServer of interface and get the things done
		if serv.InstanceIds != nil {
			serverin := server.DeleteServerInput{InstanceIds: serv.InstanceIds, GetRaw: serv.GetRaw}
			server_response, serverr := serverin.DeleteServer(authinpt)
			if serverr != nil {
				return DeleteServerResponse{}, serverr
			}
			return DeleteServerResponse{AwsResponse: server_response}, nil
		} else if serv.VpcId != "" {
			serverin := server.DeleteServerInput{VpcId: serv.VpcId, GetRaw: serv.GetRaw}
			server_response, serverr := serverin.DeleteServerFromVpc(authinpt)
			if serverr != nil {
				return DeleteServerResponse{}, serverr
			}
			return DeleteServerResponse{AwsResponse: server_response}, nil
		} else {
			return DeleteServerResponse{}, fmt.Errorf("You have not passed valid input to get details of server, the input struct looks like empty")
		}

	case "azure":
		return DeleteServerResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return DeleteServerResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return DeleteServerResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:

		log.Info("")
		log.Error("I feel we are lost in getting details of all the server :S, guess you have entered wrong cloud")
		log.Info("")

		return DeleteServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteServer")
	}
}
