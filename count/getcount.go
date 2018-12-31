package GetCount

import (
	"neuron/awsinterface"
	log "neuron/logger"
	"strings"
)

type GetCountInput struct {
	Cloud  string
	Region string
}

type GetCountResponse struct {
	AwsResponse     DengineAwsInterface.GetCountResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                               `json:"AzureResponse,omitempty"`
	DefaultResponse string                               `json:"DefaultResponse,omitempty"`
}

// being create_network my job is to create network and give back the response who called me
func (c GetCountInput) GetCount() GetCountResponse {

	var get_count_response GetCountResponse
	switch strings.ToLower(c.Cloud) {
	case "aws":

		// I will establish session so that we can carry out the process in cloud
		session_input := DengineAwsInterface.EstablishConnectionInput{c.Region, "ec2"}
		session_input.EstablishConnection()

		get_count_response = GetCountResponse{AwsResponse: DengineAwsInterface.GetCount()}

	case "azure":
		get_count_response = GetCountResponse{DefaultResponse: "We have not reached to azure yet"}
	case "gcp":
		get_count_response = GetCountResponse{DefaultResponse: "We have not reached to gcp yet"}
	case "openstack":
		get_count_response = GetCountResponse{DefaultResponse: "We have not reached to openstack yet"}
	default:
		log.Info("")
		log.Error("I feel we are lost in getting details of all the networks :S")
		log.Info("")
		get_count_response = GetCountResponse{DefaultResponse: "I feel we are lost in getting details of all the networks :S. Not found valid response?, check applog for more info"}
	}
	return get_count_response
}
