package awsinstance

import (
	"fmt"
	aws "github.com/nikhilsbhat/neuron/cloud/aws/interface"
)

type DeleteServerInput struct {
	VpcId       string   `json:"VpcId,omitempty"`
	InstanceIds []string `json:"InstanceIds,omitempty"`
	GetRaw      bool
}

func (d *DeleteServerInput) DeleteServer(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	instance_search_input := CommonComputeInput{InstanceIds: d.InstanceIds}
	search_instance, serverr := instance_search_input.SearchInstance(con)

	if serverr != nil {
		return nil, serverr
	}

	if search_instance != true {
		return nil, fmt.Errorf("Could not find the entered Instances, please enter valid/existing InstanceIds")
	}
	delete_result, ins_term_err := ec2.DeleteInstance(
		&aws.DeleteComputeInput{
			InstanceIds: d.InstanceIds,
		},
	)
	if ins_term_err != nil {
		return nil, ins_term_err
	}

	waiterr := ec2.WaitTillInstanceTerminated(
		&aws.DescribeComputeInput{
			InstanceIds: d.InstanceIds,
		},
	)
	if waiterr != nil {
		return nil, waiterr
	}

	result, err := ec2.DescribeInstance(
		&aws.DescribeComputeInput{
			InstanceIds: d.InstanceIds,
		},
	)
	if err != nil {
		return nil, err
	}

	delete_response := make([]ServerResponse, 0)

	if d.GetRaw == true {
		delete_response = append(delete_response, ServerResponse{DeleteInstRaw: delete_result, Cloud: "Amazon"})
		return delete_response, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			delete_response = append(delete_response, ServerResponse{InstanceId: *instance.InstanceId, CurrentState: *instance.State.Name})
		}
	}
	return delete_response, nil
}

func (d *DeleteServerInput) DeleteServerFromVpc(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	instance_search_input := DescribeInstanceInput{
		VpcIds: []string{d.VpcId},
	}
	search_instance, serverr := instance_search_input.GetServersFromNetwork(con)
	if serverr != nil {
		return nil, serverr
	}

	insatanceids := make([]string, 0)
	for _, instanceid := range search_instance {
		insatanceids = append(insatanceids, instanceid.InstanceId)
	}

	result, serv_del_err := ec2.DeleteInstance(
		&aws.DeleteComputeInput{
			InstanceIds: insatanceids,
		},
	)
	if serv_del_err != nil {
		return nil, serv_del_err
	}

	delete_response := make([]ServerResponse, 0)

	if d.GetRaw == true {
		delete_response = append(delete_response, ServerResponse{DeleteInstRaw: result, Cloud: "Amazon"})
		return delete_response, nil
	}

	for _, instance := range result.TerminatingInstances {
		delete_response = append(delete_response, ServerResponse{InstanceId: *instance.InstanceId, CurrentState: *instance.CurrentState.Name})
	}

	return delete_response, nil
}
