package networkDelete

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	network "neuron/cloud/aws/operations/network"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type DeleteNetworkResponse struct {
	AwsResponse     network.DeleteNetworkResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                        `json:"AzureResponse,omitempty"`
	DefaultResponse string                        `json:"DefaultResponse,omitempty"`
}

// being create_network my job is to create network and give back the response who called me
func (net *DeleteNetworkInput) DeleteNetwork() (DeleteNetworkResponse, error) {

	switch strings.ToLower(net.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: net.Profile, Cloud: net.Cloud})
		if err != nil {
			return DeleteNetworkResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: net.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Region, Resource: "ec2", Session: sess}

		// deleteting network from aws
		networkin := new(network.DeleteNetworkInput)
		networkin.VpcIds = net.VpcIds
		networkin.GetRaw = net.GetRaw
		response, net_err := networkin.DeleteNetwork(authinpt)
		if net_err != nil {
			return DeleteNetworkResponse{}, net_err
		}
		return DeleteNetworkResponse{AwsResponse: response}, nil

	case "azure":
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error(common.DefaultCloudResponse + "DeleteNetwork")
		log.Info("")
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteNetwork")
	}
}
