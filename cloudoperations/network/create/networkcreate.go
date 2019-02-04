package networkCreate

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
type CreateNetworkResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse network.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// Being CreateNetwork, job of him is to create network
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (net *NetworkCreateInput) CreateNetwork() (CreateNetworkResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateNetwork")
	}

	switch strings.ToLower(net.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: net.Cloud.Profile,
				Cloud:   net.Cloud.Name,
			},
		)

		if err != nil {
			return CreateNetworkResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: net.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// Fetching all the networks accross cloud aws
		networkin := new(network.NetworkCreateInput)
		networkin.Name = net.Name
		networkin.VpcCidr = net.VpcCidr
		networkin.SubCidrs = net.SubCidr
		networkin.Type = net.Type
		networkin.Ports = net.Ports
		networkin.GetRaw = net.Cloud.GetRaw
		response, net_err := networkin.CreateNetwork(authinpt)
		if net_err != nil {
			return CreateNetworkResponse{}, net_err
		}
		return CreateNetworkResponse{AwsResponse: response}, nil

	case "azure":
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateNetwork")
	}
}

func New() *NetworkCreateInput {
	net := &NetworkCreateInput{}
	return net
}
