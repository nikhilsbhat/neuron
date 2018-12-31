package loadbalancer

import (
	"fmt"
	aws "neuron/cloud/aws/interface"
	"strings"
	"time"
)

type GetLoadbalancerInput struct {

	//optional parameter; The names of the loadbalancers in array of which the information has to be fetched (both classic/network kind of loadbalancers).
	//one can omit this if he/she is passing ARN's of loadbalancers.
	//this parameter is mandatory if one wants to fetch the data of classic load balancers.
	LbNames []string `json:"LbNames,omitempty"`

	//optional parameter; The ARN's of the loadbalancers in array of which the information has to be fetched (only application kind of loadbalancers) one can omit this if he/she is passing names of loadbalancers.
	LbArns []string `json:"LbArns,omitempty"`

	//optional parameter if getallloadbalancer is used; Type of loadbalancers to fetch the appropriate data (classic/application).
	Type string `json:"LbArns,omitempty"`

	//optional parameter; Only when you need unfiltered result from cloud, enable this field by setting it to true. By default it is set to false.
	GetRaw bool `json:"GetRaw"`
}

func (lb *GetLoadbalancerInput) GetAllLoadbalancer(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	lb_chn := make(chan interface{}, 2)
	defer close(lb_chn)
	go func() {
		lbs, err := lb.GetAllClassicLb(con)
		time.Sleep(time.Second * 1)
		if err != nil {
			lb_chn <- err
		} else {
			lb_chn <- lbs
		}
	}()
	go func() {
		lbs, err := lb.GetAllApplicationLb(con)
		time.Sleep(time.Second * 2)
		if err != nil {
			lb_chn <- err
		} else {
			lb_chn <- lbs
		}
	}()

	// this is just a workaround and has to be fixed soon.
	classiclb := <-lb_chn
	applicationlb := <-lb_chn

	response := new(LoadBalanceResponse)
	switch calssic := classiclb.(type) {
	case []LoadBalanceResponse:
		response.ClassicLb = calssic
	case error:
		return nil, calssic
	default:
		return nil, fmt.Errorf("An unknown error occured while returning classiclb data")
	}

	switch app := applicationlb.(type) {
	case []LoadBalanceResponse:
		response.ApplicationLb = app
	case error:
		return nil, app
	default:
		return nil, fmt.Errorf("An unknown error occured while returning classiclb data")
	}

	resp := make([]LoadBalanceResponse, 0)
	resp = append(resp, *response)
	return resp, nil
}

func (load *GetLoadbalancerInput) GetAllClassicLb(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	// searching all the classic loadbalancer
	search_lb_result, search_err := elb.DescribeAllClassicLoadbalancer(
		&aws.DescribeLoadbalancersInput{},
	)
	if search_err != nil {
		return nil, search_err
	}

	lb_list := make([]LoadBalanceResponse, 0)

	if load.GetRaw == true {
		lb_list = append(lb_list, LoadBalanceResponse{GetClassicLbsRaw: search_lb_result})
		return lb_list, nil
	}
	for _, lb := range search_lb_result.LoadBalancerDescriptions {
		lb_list = append(lb_list, LoadBalanceResponse{Name: *lb.LoadBalancerName, LbDns: *lb.DNSName, Createdon: (*lb.CreatedTime).String(), Type: "classic", Scheme: *lb.Scheme, VpcId: *lb.VPCId})
	}
	return lb_list, nil
}

func (load *GetLoadbalancerInput) GetAllApplicationLb(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}
	// searching all the application loadbalancer
	search_lb_result, search_err := elb.DescribeAllApplicationLoadbalancer(
		&aws.DescribeLoadbalancersInput{},
	)
	if search_err != nil {
		return nil, search_err
	}

	lb_list := make([]LoadBalanceResponse, 0)

	for _, lb := range search_lb_result.LoadBalancers {

		// searching target group for the corresponding loadbalancer
		search_target, tarerr := elb.DescribeTargetgroups(
			&aws.DescribeLoadbalancersInput{
				LbArns: []string{*lb.LoadBalancerArn},
			},
		)
		if tarerr != nil {
			return nil, tarerr
		}

		// searching listners for the corresponding loadbalancer
		search_listners, lisserr := elb.DescribeListners(
			&aws.DescribeLoadbalancersInput{
				LbArns: []string{*lb.LoadBalancerArn},
			},
		)
		if lisserr != nil {
			return nil, lisserr
		}

		response := new(LoadBalanceResponse)
		if load.GetRaw == true {
			response.GetApplicationLbRaw.GetApplicationLbRaw = lb
			response.GetApplicationLbRaw.GetTargetGroupRaw = search_target
			response.GetApplicationLbRaw.GetListnersRaw = search_listners
			lb_list = append(lb_list, *response)
		} else {

			tar_arn := make([]string, 0)
			for _, tar := range search_target.TargetGroups {
				tar_arn = append(tar_arn, *tar.TargetGroupArn)
			}

			lis_arn := make([]string, 0)
			for _, lis := range search_listners.Listeners {
				lis_arn = append(lis_arn, *lis.ListenerArn)
			}

			response.Name = *lb.LoadBalancerName
			response.LbDns = *lb.DNSName
			response.LbArn = *lb.LoadBalancerArn
			response.Createdon = (*lb.CreatedTime).String()
			response.Type = *lb.Type
			response.Scheme = *lb.Scheme
			response.VpcId = *lb.VpcId
			response.TargetArn = tar_arn
			response.ListnerArn = lis_arn
			lb_list = append(lb_list, *response)
		}
	}
	return lb_list, nil
}

func (lb *GetLoadbalancerInput) Getloadbalancers(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	switch strings.ToLower(lb.Type) {
	case "classic":
		getlb, err := lb.GetClassicloadbalancers(con)
		if err != nil {
			return nil, err
		}
		return getlb, nil
	case "application":
		getlb, err := lb.GetApplicationloadbalancers(con)
		if err != nil {
			return nil, err
		}
		return getlb, nil
	default:
		return nil, fmt.Errorf("You provided unknown loadbalancer type, enter a valid LB type")
	}
}

func (load *GetLoadbalancerInput) GetClassicloadbalancers(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	lbresponse, err := elb.DescribeClassicLoadbalancer(
		&aws.DescribeLoadbalancersInput{
			LbNames: load.LbNames,
		},
	)
	if err != nil {
		return nil, err
	}
	lb_response := make([]LoadBalanceResponse, 0)

	if load.GetRaw == true {
		lb_response = append(lb_response, LoadBalanceResponse{GetClassicLbsRaw: lbresponse})
	}

	for _, lb := range lbresponse.LoadBalancerDescriptions {
		response := new(LoadBalanceResponse)
		response.Name = *lb.LoadBalancerName
		response.LbDns = *lb.DNSName
		response.Createdon = (*lb.CreatedTime).String()
		response.Type = "classic"
		response.Scheme = *lb.Scheme
		response.VpcId = *lb.VPCId
		lb_response = append(lb_response, *response)
	}

	return lb_response, nil
}

func (load *GetLoadbalancerInput) GetApplicationloadbalancers(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	lbin := new(aws.DescribeLoadbalancersInput)
	lbin.LbNames = load.LbNames
	lbin.LbArns = load.LbArns
	get_loadbalancer, err := elb.DescribeApplicationLoadbalancer(lbin)
	if err != nil {
		return nil, err
	}
	lb_response := make([]LoadBalanceResponse, 0)
	for _, lb := range get_loadbalancer.LoadBalancers {

		// searching target group for the corresponding loadbalancer
		lbin.LbArns = []string{*lb.LoadBalancerArn}
		search_target, tarerr := elb.DescribeTargetgroups(lbin)
		if tarerr != nil {
			return nil, tarerr
		}

		// searching listners for the corresponding loadbalancer
		search_listners, lisserr := elb.DescribeListners(lbin)
		if lisserr != nil {
			return nil, lisserr
		}

		response := new(LoadBalanceResponse)
		if load.GetRaw == true {
			response.GetApplicationLbRaw.GetApplicationLbRaw = lb
			response.GetApplicationLbRaw.GetTargetGroupRaw = search_target
			response.GetApplicationLbRaw.GetListnersRaw = search_listners
			lb_response = append(lb_response, *response)
		} else {
			tar_arn := make([]string, 0)
			for _, tar := range search_target.TargetGroups {
				tar_arn = append(tar_arn, *tar.TargetGroupArn)
			}

			lis_arn := make([]string, 0)
			for _, lis := range search_listners.Listeners {
				lis_arn = append(lis_arn, *lis.ListenerArn)
			}

			response.Name = *lb.LoadBalancerName
			response.LbDns = *lb.DNSName
			response.LbArn = *lb.LoadBalancerArn
			response.Createdon = (*lb.CreatedTime).String()
			response.Type = *lb.Type
			response.Scheme = *lb.Scheme
			response.VpcId = *lb.VpcId
			response.TargetArn = tar_arn
			response.ListnerArn = lis_arn
			lb_response = append(lb_response, *response)
		}
	}
	return lb_response, nil
}

func (lb *GetLoadbalancerInput) FindClassicLoadbalancer(con aws.EstablishConnectionInput) (bool, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return false, sesserr
	}

	get_loadbalancer, err := elb.DescribeClassicLoadbalancer(
		&aws.DescribeLoadbalancersInput{
			LbNames: lb.LbNames,
		},
	)
	if err != nil {
		return false, err
	}
	if len(get_loadbalancer.LoadBalancerDescriptions) != 0 {
		return true, nil
	}
	return false, fmt.Errorf("Could not find the entered loadbalancer, please enter valid/existing loadbalancer Name")
}

func (lb *GetLoadbalancerInput) FindApplicationLoadbalancer(con aws.EstablishConnectionInput) (bool, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return false, sesserr
	}

	get_loadbalancer, err := elb.DescribeApplicationLoadbalancer(
		&aws.DescribeLoadbalancersInput{
			LbNames: lb.LbNames,
			LbArns:  lb.LbArns,
		},
	)
	if err != nil {
		return false, err
	}
	if len(get_loadbalancer.LoadBalancers) != 0 {
		return true, nil
	}
	return false, fmt.Errorf("Could not find the entered loadbalancer, please enter valid/existing loadbalancer ARN/Name")
}

func (lb *GetLoadbalancerInput) GetArnFromLoadbalancer(con aws.EstablishConnectionInput) (LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return LoadBalanceResponse{}, sesserr
	}

	get_loadbalancer, err := elb.DescribeApplicationLoadbalancer(
		&aws.DescribeLoadbalancersInput{
			LbNames: lb.LbNames,
		},
	)
	if err != nil {
		return LoadBalanceResponse{}, err
	}
	arns := make([]string, 0)
	for _, lb := range get_loadbalancer.LoadBalancers {
		arns = append(arns, *lb.LoadBalancerArn)
	}

	return LoadBalanceResponse{LbArns: arns}, nil
}
