package miscOperations

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	awscommon "neuron/cloud/aws/operations/common"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	support "neuron/cloudoperations/support"
	//log "neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type GetRegionsResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse awscommon.CommonResponse `json:"Regions,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// being create_network my job is to create network and give back the response who called me
func (reg *GetRegionInput) GetRegions() (GetRegionsResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(reg.Cloud)); status != true {
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(reg.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: reg.Profile,
				Cloud:   reg.Cloud,
			},
		)
		if err != nil {
			return GetRegionsResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: reg.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: reg.Region, Resource: "ec2", Session: sess}

		// I will call create_vpc and get the things done
		regionin := awscommon.CommonInput{}
		regionin.GetRaw = reg.GetRaw
		response, reg_err := regionin.GetRegions(authinpt)
		if reg_err != nil {
			return GetRegionsResponse{}, reg_err
		}
		return GetRegionsResponse{AwsResponse: response}, nil

	case "azure":
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetRegions")
	}
}

func New() *GetRegionInput {
	net := &GetRegionInput{}
	return net
}
