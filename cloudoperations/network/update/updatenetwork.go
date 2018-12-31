package networkUpdate

import (
	auth "neuron/cloud/aws/interface"
	awsnetwork "neuron/cloud/aws/operations/network"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type UpdateNetworkResponse struct {
	AwsResponse     awsnetwork.NetworkResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                     `json:"AzureResponse,omitempty"`
	DefaultResponse string                     `json:"DefaultResponse,omitempty"`
}

// being create_network my job is to create network and give back the response who called me
func (net *NetworkUpdateInput) UpdateNetwork() (UpdateNetworkResponse, error) {

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
		log.Info("")
		log.Error("I feel we are lost in updating NETWORK, was unable to find the appropriate cloud ")
		log.Error("You might have entered wrong cloud name else misspelled it, please check it before passing")
		log.Info("")
		return UpdateNetworkResponse{DefaultResponse: common.DefaultCloudResponse + "NetworkUpdate"}, nil
	}
}
