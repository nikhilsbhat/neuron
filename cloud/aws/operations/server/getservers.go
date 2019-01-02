package awsinstance

import (
	"neuron/cloud/aws/interface"
	//"neuron/logger"
	"strings"
	//"time"
	"fmt"
)

type DescribeInstanceInput struct {
	InstanceIds []string `json:"InstanceIds,omitempty"`
	VpcIds      []string `json:"VpcIds,omitempty"`
	SubnetIds   []string `json:"SubnetIds,omitempty"`
	Filters     Filters  `json:"Filters,omitempty"`
	GetRaw      bool
}

type Filters struct {
	Name  string
	Value []string
}

// This function is tailored to fectch the servers from network, to fetch the data one has to pass either subnet-id else vpc-id to filters to make the life easy.
func (d *DescribeInstanceInput) GetServersFromNetwork(con neuronaws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, des_inst_err := ec2.DescribeInstance(
		&neuronaws.DescribeComputeInput{
			Filters: neuronaws.Filters{
				Name:  "vpc-id",
				Value: d.VpcIds,
			},
		},
	)

	if des_inst_err != nil {
		return nil, des_inst_err
	}

	server_response := make([]ServerResponse, 0)

	if d.GetRaw == true {
		server_response = append(server_response, ServerResponse{GetInstRaw: result, Cloud: "Amazon"})
		return server_response, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			switch strings.ToLower(*instance.State.Name) {
			case "running":
				server_response = append(server_response, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PublicIpAddress: *instance.PublicIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "stopped":
				server_response = append(server_response, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "terminated":
				server_response = append(server_response, ServerResponse{State: *instance.State.Name, Cloud: "Amazon"})
			default:
				return nil, fmt.Errorf("Oops...!!!!. few instances are not in a state of fetching its details, check back after few seconds")
			}
		}
	}
	return server_response, nil
}

func (d *DescribeInstanceInput) GetServersFromSubnet(con neuronaws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, des_inst_err := ec2.DescribeInstance(
		&neuronaws.DescribeComputeInput{
			Filters: neuronaws.Filters{
				Name:  "subnet-id",
				Value: d.SubnetIds,
			},
		},
	)

	if des_inst_err != nil {
		return nil, des_inst_err
	}

	server_response := make([]ServerResponse, 0)

	if d.GetRaw == true {
		server_response = append(server_response, ServerResponse{GetInstRaw: result, Cloud: "Amazon"})
		return server_response, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {

			switch strings.ToLower(*instance.State.Name) {
			case "running":
				server_response = append(server_response, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PublicIpAddress: *instance.PublicIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "stopped":
				server_response = append(server_response, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "terminated":
				server_response = append(server_response, ServerResponse{State: *instance.State.Name, Cloud: "Amazon"})
			default:
				return nil, fmt.Errorf("Oops...!!!!. instances are not in a state of fetching the details of it, check back after few minutes")
			}
		}
	}
	return server_response, nil
}

// This function is meant to get all the servers from a particular region.
func (d *DescribeInstanceInput) GetAllServers(con neuronaws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, err := ec2.DescribeAllInstances(
		&neuronaws.DescribeComputeInput{},
	)
	if err != nil {
		return nil, err
	}

	server_response := make([]ServerResponse, 0)

	if d.GetRaw == true {
		server_response = append(server_response, ServerResponse{GetInstRaw: result, Cloud: "Amazon"})
		return server_response, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if (*instance.State.Name == "running") || (*instance.State.Name == "stopped") {
				server_response = append(server_response, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, InstanceType: *instance.InstanceType, Cloud: "Amazon", Region: *instance.Placement.AvailabilityZone})

			} else {
				// change has to be made here (introduction of omitempty is required)
				server_response = append(server_response, ServerResponse{State: "terminated", Cloud: "Amazon"})
			}
		}
	}
	return server_response, nil
}

//This function is tailored to get the details of the random servers you enter
func (d *DescribeInstanceInput) GetServersDetails(con neuronaws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, des_inst_err := ec2.DescribeInstance(
		&neuronaws.DescribeComputeInput{
			InstanceIds: d.InstanceIds,
		},
	)
	if des_inst_err != nil {
		return nil, des_inst_err
	}

	server_response := make([]ServerResponse, 0)

	if d.GetRaw == true {
		server_response = append(server_response, ServerResponse{GetInstRaw: result, Cloud: "Amazon"})
		return server_response, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {

			switch strings.ToLower(*instance.State.Name) {
			case "running":
				server_response = append(server_response, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PublicIpAddress: *instance.PublicIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "stopped":
				server_response = append(server_response, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "terminated":
				server_response = append(server_response, ServerResponse{State: *instance.State.Name, Cloud: "Amazon"})
			default:
				return nil, fmt.Errorf("Oops...!!!!. instances are not in a state of fetching the details of it, check back after few minutes")
			}
		}
	}
	return server_response, nil
}
