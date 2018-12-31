package miscOperations

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	awscommon "neuron/cloud/aws/operations/common"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type GetRegionsResponse struct {
	Regions         []string `json:"Regions,omitempty"`
	DefaultResponse string   `json:"DefaultResponse,omitempty"`
}

// being create_network my job is to create network and give back the response who called me
func (reg *GetRegionInput) GetRegions() (GetRegionsResponse, error) {

	switch strings.ToLower(reg.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: reg.Profile, Cloud: reg.Cloud})
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
		response, reg_err := regionin.GetRegions(authinpt)
		if reg_err != nil {
			return GetRegionsResponse{}, reg_err
		}
		return GetRegionsResponse{Regions: response}, nil

	case "azure":
		return GetRegionsResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return GetRegionsResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return GetRegionsResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:
		log.Info("")
		log.Error("I feel we are lost in getting list of REGIONS, was unable to find the appropriate cloud :S")
		log.Error("You have entered wrong cloud name else misspelled the name of it, please check it before passing")
		log.Info("")
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetRegions")
	}
}