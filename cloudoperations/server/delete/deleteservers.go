package deleteServer

import (
	"fmt"
	auth "github.com/nikhilsbhat/neuron/cloud/aws/interface"
	server "github.com/nikhilsbhat/neuron/cloud/aws/operations/server"
	awssess "github.com/nikhilsbhat/neuron/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron/cloudoperations/support"
	//log "neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type DeleteServerResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []server.ServerResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// Being DeleteServer, job of him is to delete servers as per the instructions passed to him
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (serv *DeleteServersInput) DeleteServer() (DeleteServerResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
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
			return DeleteServerResponse{}, crederr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorize
		authinpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call DeleteServer of interface and get the things done
		if serv.InstanceIds != nil {
			serverin := server.DeleteServerInput{}
			serverin.InstanceIds = serv.InstanceIds
			serverin.GetRaw = serv.Cloud.GetRaw
			server_response, serverr := serverin.DeleteServer(authinpt)
			if serverr != nil {
				return DeleteServerResponse{}, serverr
			}
			return DeleteServerResponse{AwsResponse: server_response}, nil
		} else if serv.VpcId != "" {
			serverin := server.DeleteServerInput{}
			serverin.VpcId = serv.VpcId
			serverin.GetRaw = serv.Cloud.GetRaw
			server_response, serverr := serverin.DeleteServerFromVpc(authinpt)
			if serverr != nil {
				return DeleteServerResponse{}, serverr
			}
			return DeleteServerResponse{AwsResponse: server_response}, nil
		} else {
			return DeleteServerResponse{}, fmt.Errorf("You have not passed valid input to get details of server, the input looks like empty")
		}

	case "azure":
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteServer")
	}
}

func New() *DeleteServersInput {
	net := &DeleteServersInput{}
	return net
}
