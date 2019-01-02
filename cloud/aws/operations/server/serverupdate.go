package awsinstance

import (
	"fmt"
	aws "neuron/cloud/aws/interface"
	"strings"
)

type UpdateServerInput struct {
	InstanceIds []string `json:"InstanceIds,omitempty"`
	Action      string
	GetRaw      bool
}

func (u *UpdateServerInput) UpdateServer(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	search_input := CommonComputeInput{InstanceIds: u.InstanceIds}
	search, serverr := search_input.SearchInstance(con)
	if serverr != nil {
		return nil, serverr
	}

	if search != true {
		return nil, fmt.Errorf("Could not find the entered Instances, please enter valid/existing InstanceIds")
	}
	server_response := make([]ServerResponse, 0)

	switch strings.ToLower(u.Action) {
	case "start":
		result, start_err := ec2.StartInstances(
			&aws.UpdateComputeInput{
				InstanceIds: u.InstanceIds,
			},
		)

		if start_err != nil {
			return nil, start_err
		}

		wait_err := ec2.WaitTillInstanceRunning(
			&aws.DescribeComputeInput{
				InstanceIds: u.InstanceIds,
			},
		)
		if wait_err != nil {
			return nil, wait_err
		}

		if u.GetRaw == true {
			server_response = append(server_response, ServerResponse{StartInstRaw: result, Cloud: "Amazon"})
			return server_response, nil
		}

		for _, inst := range result.StartingInstances {
			server_response = append(server_response, ServerResponse{InstanceId: *inst.InstanceId, CurrentState: "running", PreviousState: *inst.PreviousState.Name})
		}
		return server_response, nil

	case "stop":
		result, stop_err := ec2.StopInstances(
			&aws.UpdateComputeInput{
				InstanceIds: u.InstanceIds,
			},
		)

		if stop_err != nil {
			return nil, stop_err
		}
		wait_err := ec2.WaitTillInstanceStopped(
			&aws.DescribeComputeInput{
				InstanceIds: u.InstanceIds,
			},
		)

		if wait_err != nil {
			return nil, wait_err
		}

		if u.GetRaw == true {
			server_response = append(server_response, ServerResponse{StopInstRaw: result, Cloud: "Amazon"})
			return server_response, nil
		}

		for _, inst := range result.StoppingInstances {
			server_response = append(server_response, ServerResponse{InstanceId: *inst.InstanceId, CurrentState: "stopped", PreviousState: *inst.PreviousState.Name})
		}
		return server_response, nil

	default:
		return nil, fmt.Errorf("Sorry...!!!!. I am not aware of the action you asked me to perform, please enter the action which we support. The available actions are: start/stop")
	}
}
