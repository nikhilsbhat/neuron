package awsinstance

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/nikhilsbhat/neuron/cloud/aws/interface"
	common "github.com/nikhilsbhat/neuron/cloud/aws/operations/common"
	network "github.com/nikhilsbhat/neuron/cloud/aws/operations/network"
	"strconv"
)

type CreateServerInput struct {
	InstanceName string
	ImageId      string
	InstanceType string
	KeyName      string
	MaxCount     int64
	MinCount     int64
	SubnetId     string
	SecGroupId   string
	UserData     string
	AssignPubIp  bool
	GetRaw       bool
}

type ServerResponse struct {
	InstanceName        string                        `json:"InstanceName,omitempty"`
	InstanceId          string                        `json:"InstanceId,omitempty"`
	SubnetId            string                        `json:"SubnetId,omitempty"`
	PrivateIpAddress    string                        `json:"IpAddress,omitempty"`
	PublicIpAddress     string                        `json:"PublicIpAddress,omitempty"`
	PrivateDnsName      string                        `json:"PrivateDnsName,omitempty"`
	CreatedOn           string                        `json:"CreatedOn,omitempty"`
	State               string                        `json:"State,omitempty"`
	InstanceDeleteState string                        `json:"InstanceDeleteState,omitempty"`
	InstanceType        string                        `json:"InstanceType,omitempty"`
	Cloud               string                        `json:"Cloud,omitempty"`
	Region              string                        `json:"Region,omitempty"`
	PreviousState       string                        `json:"PreviousState,omitempty"`
	CurrentState        string                        `json:"CurrentState,omitempty"`
	DefaultResponse     interface{}                   `json:"DefaultResponse,omitempty"`
	Error               error                         `json:"Error,omitempty"`
	CreateInstRaw       *ec2.DescribeInstancesOutput  `json:"CreateInstRaw,omitempty"`
	GetInstRaw          *ec2.DescribeInstancesOutput  `json:"DescribeInstRaw,omitempty"`
	DeleteInstRaw       *ec2.TerminateInstancesOutput `json:"DeleteInstRaw,omitempty"`
	StartInstRaw        *ec2.StartInstancesOutput     `json:"StartInstRaw,omitempty"`
	StopInstRaw         *ec2.StopInstancesOutput      `json:"StopInstRaw,omitempty"`
	CreateImgRaw        *ec2.CreateImageOutput        `json:"CreateImgRaw,omitempty"`
	DescribeImg         *ec2.DescribeImagesOutput     `json:"DescribeImg,omitempty"`
}

func (csrv *CreateServerInput) CreateServer(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	// I will make a decision which security group to pick
	sub_input := network.GetNetworksInput{SubnetIds: []string{csrv.SubnetId}}
	sub_result, suberr := sub_input.FindSubnet(con)
	if suberr != nil {
		return nil, suberr
	}

	if sub_result != true {
		return nil, fmt.Errorf("Could not find the entered SUBNET, please enter valid/existing SUBNET id")
	}

	inst := new(aws.CreateServerInput)

	switch csrv.SecGroupId {
	case "":
		vpc_res, vpcerr := sub_input.GetVpcFromSubnet(con)
		if vpcerr != nil {
			return nil, vpcerr
		}

		sec_input := network.NetworkComponentInput{VpcIds: []string{vpc_res.VpcId}}
		sec_res, secerr := sec_input.GetSecFromVpc(con)
		if secerr != nil {
			return nil, nil
		}
		inst.SecurityGroups = sec_res.SecGroupIds

	default:
		inst.SecurityGroups = []string{csrv.SecGroupId}
	}

	// I will be the spoc for the instance creation with the userdata passed to me
	switch csrv.UserData {
	case "":
		inst.UserData = b64.StdEncoding.EncodeToString([]byte("echo 'nothing'"))
	default:
		inst.UserData = b64.StdEncoding.EncodeToString([]byte(csrv.UserData))
	}

	switch csrv.MinCount {
	case 0:
		inst.MinCount = 1
	default:
		inst.MinCount = csrv.MinCount
	}

	switch csrv.MaxCount {
	case 0:
		inst.MaxCount = 1
	default:
		inst.MaxCount = csrv.MaxCount
	}

	inst.ImageId = csrv.ImageId
	inst.InstanceType = csrv.InstanceType
	inst.KeyName = csrv.KeyName
	inst.AssignPubIp = csrv.AssignPubIp
	inst.SubnetId = csrv.SubnetId
	// support for custom ebs mapping will be rolled out soon
	server_create_result, err := ec2.CreateInstance(inst)

	if err != nil {
		return nil, err
	}

	instance_ids := make([]string, 0)
	for _, instance := range server_create_result.Instances {
		instance_ids = append(instance_ids, *instance.InstanceId)
	}

	// I will make program wait untill instance become running
	wait_err := ec2.WaitTillInstanceAvailable(
		&aws.DescribeComputeInput{
			InstanceIds: instance_ids,
		},
	)
	if wait_err != nil {
		return nil, wait_err
	}

	// creating tags for the server
	for i, instance := range instance_ids {
		tags := common.Tag{instance, "Name", csrv.InstanceName + "-" + strconv.Itoa(i)}
		_, tag_err := tags.CreateTags(con)
		if tag_err != nil {
			return nil, tag_err
		}
	}

	//fetching the deatils of server
	result, serverr := ec2.DescribeInstance(
		&aws.DescribeComputeInput{
			InstanceIds: instance_ids,
		},
	)
	if serverr != nil {
		return nil, serverr
	}

	type server_response struct {
		name        string
		instance_id string
		ipaddress   string
		privatedns  string
		publicIp    string
		createdon   string
	}

	response := make([]server_response, 0)
	create_server_response := make([]ServerResponse, 0)

	if csrv.GetRaw == true {
		create_server_response = append(create_server_response, ServerResponse{CreateInstRaw: result, Cloud: "Amazon"})
		return create_server_response, nil
	}

	// fetching the instance details which is created in previos process
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if csrv.AssignPubIp == true {
				response = append(response, server_response{name: *instance.Tags[0].Value, instance_id: *instance.InstanceId, ipaddress: *instance.PrivateIpAddress, privatedns: *instance.PrivateDnsName, publicIp: *instance.PublicIpAddress, createdon: (*instance.LaunchTime).String()})
			} else {
				response = append(response, server_response{name: *instance.Tags[0].Value, instance_id: *instance.InstanceId, ipaddress: *instance.PrivateIpAddress, privatedns: *instance.PrivateDnsName, createdon: (*instance.LaunchTime).String()})
			}
		}
	}

	for _, server := range response {
		create_server_response = append(create_server_response, ServerResponse{InstanceName: server.name, InstanceId: server.instance_id, SubnetId: csrv.SubnetId, PrivateIpAddress: server.ipaddress, PublicIpAddress: server.publicIp, PrivateDnsName: server.privatedns, CreatedOn: server.createdon, Cloud: "Amazon"})
	}

	return create_server_response, nil
}
