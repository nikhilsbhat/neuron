package networkDelete

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	network "neuron/cloud/aws/operations/network"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	support "neuron/cloudoperations/support"
	//log "neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type DeleteNetworkResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse network.DeleteNetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// Being DeleteNetwork, job of him is to delete the network and its components
// and give back the response who called him.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (net *DeleteNetworkInput) DeleteNetwork() (DeleteNetworkResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud)); status != true {
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "UpdateNetwork")
	}

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
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteNetwork")
	}
}

func New() *DeleteNetworkInput {
	net := &DeleteNetworkInput{}
	return net
}
