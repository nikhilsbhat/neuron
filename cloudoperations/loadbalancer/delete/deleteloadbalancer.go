package deleteLoadbalancer

import (
	"fmt"
	auth "neuron/cloud/aws/interface"
	loadbalance "neuron/cloud/aws/operations/loadbalancer"
	awssess "neuron/cloud/aws/sessions"
	common "neuron/cloudoperations/common"
	log "neuron/logger"
	"strings"
)

type LoadBalancerDeleteResponse struct {
	AwsResponse     []loadbalance.LoadBalanceDeleteResponse `json:"AwsResponse,omitempty"`
	AzureResponse   string                                  `json:"AzureResponse,omitempty"`
	DefaultResponse string                                  `json:"DefaultResponse,omitempty"`
}

// being create_loadbalancer my job is to create required loadbalancer and give back the response who called me
func (d *LbDeleteInput) DeleteLoadBalancer() (LoadBalancerDeleteResponse, error) {

	switch strings.ToLower(d.Cloud) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: d.Profile, Cloud: d.Cloud})
		if err != nil {
			return LoadBalancerDeleteResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: d.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := new(auth.EstablishConnectionInput)
		authinpt.Region = d.Region
		authinpt.Session = sess
		switch strings.ToLower(d.Type) {
		case "classic":
			authinpt.Resource = "elb"
		case "application":
			authinpt.Resource = "elb2"
		}

		lbin := new(loadbalance.DeleteLoadbalancerInput)
		lbin.LbNames = d.LbNames
		lbin.LbArns = d.LbArns
		lbin.Type = d.Type
		lbin.GetRaw = d.GetRaw
		response, lberr := lbin.DeleteLoadbalancer(*authinpt)
		if lberr != nil {
			return LoadBalancerDeleteResponse{}, lberr
		}
		return LoadBalancerDeleteResponse{AwsResponse: response}, nil

	case "azure":
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		log.Info("")
		log.Error(common.DefaultCloudResponse + "DeleteLoadBalancer")
		log.Info("")
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteLoadBalancer")
	}
}
