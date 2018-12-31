package serverCreate

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	server "neuron/cloud/aws/operations/server"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"reflect"
	"strings"
)

type ServerCreateResponse struct {
	AwsResponse     []server.ServerResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                  `json:"AzureResponse,omitempty"`
	DefaultResponse interface{}             `json:"Response,omitempty"`
}

// being create_network my job is to create network and give back the response who called me
func (serv ServerCreateInput) CreateServer() (ServerCreateResponse, error) {

	switch strings.ToLower(serv.Cloud) {
	case "aws":

		creds, crederr := common.GetCredentials(&common.GetCredentialsInput{Profile: serv.Profile, Cloud: serv.Cloud})
		if crederr != nil {
			return ServerCreateResponse{}, crederr
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: serv.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: serv.Region, Resource: "ec2", Session: sess}

		// I will call CreateServer of interface and get the things done

		serverin := server.CreateServerInput{InstanceName: serv.InstanceName, ImageId: serv.ImageId, InstanceType: serv.Flavor, KeyName: serv.KeyName, MaxCount: serv.Count, SubnetId: serv.SubnetId, UserData: serv.UserData, AssignPubIp: serv.AssignPubIp, GetRaw: serv.GetRaw}
		response, err := serverin.CreateServer(authinpt)
		if err != nil {
			return ServerCreateResponse{}, err
		}
		return ServerCreateResponse{AwsResponse: response}, nil

	case "azure":
		return ServerCreateResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return ServerCreateResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return ServerCreateResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:

		log.Info("")
		log.Error("I feel we are lost in getting details of all the server :S, guess you have entered wrong cloud")
		log.Info("")
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateServer")
	}
}

func CreateServerMock() (ServerCreateResponse, error) {

	input := new(ServerCreateInput)
	defaultval := reflect.Indirect(reflect.ValueOf(input))

	defaults := make(map[string]interface{})
	for i := 0; i < defaultval.NumField(); i++ {
		defaults[defaultval.Type().Field(i).Name] = defaultval.Type().Field(i).Type
	}

	return ServerCreateResponse{DefaultResponse: defaults},nil
}
