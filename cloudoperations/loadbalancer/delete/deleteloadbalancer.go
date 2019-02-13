package deleteLoadbalancer

import (
	"fmt"
	auth "github.com/nikhilsbhat/neuron/cloud/aws/interface"
	loadbalance "github.com/nikhilsbhat/neuron/cloud/aws/operations/loadbalancer"
	awssess "github.com/nikhilsbhat/neuron/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron/cloudoperations/support"
	log "github.com/nikhilsbhat/neuron/logger"
	"strings"
)

// The struct that will return the filtered/unfiltered responses of variuos clouds.
type LoadBalancerDeleteResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []loadbalance.LoadBalanceDeleteResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// Being GetLoadbalancers, job of him is to create loadbalancer asked for
// and give back the response who called him.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (d *LbDeleteInput) DeleteLoadBalancer() (LoadBalancerDeleteResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(d.Cloud.Name)); status != true {
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(d.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: d.Cloud.Profile,
				Cloud:   d.Cloud.Name,
			},
		)
		if err != nil {
			return LoadBalancerDeleteResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		session_input := awssess.CreateSessionInput{Region: d.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := session_input.CreateAwsSession()

		//authorizing to request further
		authinpt := new(auth.EstablishConnectionInput)
		authinpt.Region = d.Cloud.Region
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
		lbin.GetRaw = d.Cloud.GetRaw
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

func New() *LbDeleteInput {
	net := &LbDeleteInput{}
	return net
}
