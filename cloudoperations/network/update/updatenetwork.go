package networkUpdate

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	awsnetwork "neuron/cloud/aws/operations/network"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	support "neuron/cloudoperations/support"
	//log "neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type UpdateNetworkResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse awsnetwork.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// Being UpdateNetwork, job of him is to update the network and its components
// and give back the response who called him.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (net *NetworkUpdateInput) UpdateNetwork() (UpdateNetworkResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud)); status != true {
		return UpdateNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "UpdateNetwork")
	}

	switch strings.ToLower(net.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: net.Profile,
				Cloud:   net.Cloud,
			},
		)

		if err != nil {
			return UpdateNetworkResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: net.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Region, Resource: "ec2", Session: sess}

		// I will call UpdateNetwork of interface and get the things done
		serverin := awsnetwork.UpdateNetworkInput{
			Resource: net.Catageory.Resource,
			Action:   net.Catageory.Action,
			GetRaw:   net.GetRaw,
			Network: awsnetwork.NetworkCreateInput{
				Name:     net.Catageory.Name,
				VpcCidr:  net.Catageory.VpcCidr,
				VpcId:    net.Catageory.VpcId,
				SubCidrs: net.Catageory.SubCidrs,
				Type:     net.Catageory.Type,
				Ports:    net.Catageory.Ports,
				Zone:     net.Catageory.Zone,
			},
		}
		response, err := serverin.UpdateNetwork(authinpt)
		if err != nil {
			return UpdateNetworkResponse{}, err
		}
		return UpdateNetworkResponse{AwsResponse: response}, nil

	case "azure":
		return UpdateNetworkResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return UpdateNetworkResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return UpdateNetworkResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:
		return UpdateNetworkResponse{DefaultResponse: common.DefaultCloudResponse + "NetworkUpdate"}, nil
	}
}

func New() *NetworkUpdateInput {
	net := &NetworkUpdateInput{}
	return net
}
