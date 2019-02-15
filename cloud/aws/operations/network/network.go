package awsnetwork

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/nikhilsbhat/neuron/cloud/aws/interface"
	common "github.com/nikhilsbhat/neuron/cloud/aws/operations/common"
	"strconv"
	"strings"
)

type NetworkCreateInput struct {
	VpcCidr  string   `json:"vpccidr"`
	SubCidrs []string `json:"subcidrs"`
	SubCidr  string   `json:"subcidr"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Ports    []string `json:"ports"`
	Zone     string   `json:"zone"`
	VpcId    string   `json:"vpcid"`
	IgwId    string   `json:"igwid"`
	GetRaw   bool     `json:"getraw"`
}

type NetworkResponse struct {
	Name                  string                              `json:"name,omitempty"`
	VpcId                 string                              `json:"vpcid,omitempty"`
	Subnets               []SubnetReponse                     `json:"subnets,omitempty"`
	Vpcs                  []VpcResponse                       `json:"vpcs,omitempty"`
	Type                  string                              `json:"type,omitempty"`
	State                 string                              `json:"state,omitempty"`
	IgwId                 string                              `json:"igw,omitempty"`
	IsDefault             bool                                `json:"isdefault,omitempty"`
	SecGroupIds           []string                            `json:"secgroupid,omitempty"`
	Region                string                              `json:"region,omitempty"`
	GetVpcsRaw            *ec2.DescribeVpcsOutput             `json:"getvpcsraw,omitempty"`
	GetVpcRaw             *ec2.Vpc                            `json:"getvpcraw,omitempty"`
	GetSubnetRaw          *ec2.DescribeSubnetsOutput          `json:"getsubnetraw,omitempty"`
	CreateVpcRaw          VpcResponse                         `json:"createvpcraw,omitempty"`
	CreateSubnetRaw       []SubnetReponse                     `json:"createsubnetraw,omitempty"`
	CreateIgwRaw          *ec2.CreateInternetGatewayOutput    `json:"createigwraw,omitempty"`
	CreateSecRaw          *ec2.CreateSecurityGroupOutput      `json:"createsecraw,omitempty"`
	DescribeRouteTableRaw *ec2.DescribeRouteTablesOutput      `json:"describeroutetableraw,omitempty"`
	DescribeSecurityRaw   *ec2.DescribeSecurityGroupsOutput   `json:"describesecurityraw,omitempty"`
	DescribeIgwRaw        *ec2.DescribeInternetGatewaysOutput `json:"describeigwraw,omitempty"`
}

type DeleteNetworkInput struct {
	VpcIds        []string `json:"region"`
	SubnetIds     []string `json:"vpcids"`
	SecIds        []string `json:"secids"`
	IgwIds        []string `json:"igwid"`
	RouteTableIds []string `json:"routetableids"`
	GetRaw        bool     `json:"getraw"`
}

type GetNetworksInput struct {
	VpcIds    []string `json:"vpcids`
	SubnetIds []string `json:"subnetids"`
	Filters   Filters  `json:"filters"`
	Region    string   `json:"region"`
	GetRaw    bool     `json:"getraw"`
}

type DeleteNetworkResponse struct {
	Subnets         string `json:"subnets,omitempty"`
	SecurityGroups  string `json:"securities,omitempty"`
	Routetables     string `json:"routetables,omitempty"`
	Gateways        string `json:"gateways,omitempty"`
	Vpcs            string `json:"vpcs,omitempty"`
	DefaultResponse string `json:"defaultresponse,omitempty"`
	Status          string `json:"status,omitempty"`
}

type UpdateNetworkInput struct {
	Resource string             `json:"resource"`
	Network  NetworkCreateInput `json:"network"`
	Action   string             `json:"action"`
	GetRaw   bool               `json:"getRaw"`
}

type Filters struct {
	Name  string
	Value interface{}
}

// being create_vpc my job is to create vpc and give back the response who called me
func (net *NetworkCreateInput) CreateNetwork(con aws.EstablishConnectionInput) (NetworkResponse, error) {

	if (net.VpcCidr == "") || (net.Name == "") {
		return NetworkResponse{}, fmt.Errorf("You have not provided either CIDR or name for VPC, cannot proceed further")
	}

	/*get the relative sessions before proceeding further
	  ec2, seserr := con.EstablishConnection()
	  if seserr != nil {
	          return NetworkResponse{}, seserr
	  }*/

	netin := new(NetworkCreateInput)
	netin.VpcCidr = net.VpcCidr
	netin.Name = net.Name
	netin.Type = net.Type
	netin.Ports = net.Ports
	netin.GetRaw = net.GetRaw

	vpc, err := netin.CreateVpc(con)
	if err != nil {
		return NetworkResponse{}, err
	}

	zonein := common.CommonInput{}
	zones, zon_err := zonein.GetAvailabilityZones(con)
	if zon_err != nil {
		return NetworkResponse{}, zon_err
	}

	// This takes care creation of required number of subnets.
	subnets := make([]SubnetReponse, 0)

	zonenum := len(zones) - 1
	for i, sub := range net.SubCidrs {

		if zonenum < 0 {
			zonenum = len(zones) - 1
		}

		// Creating subnet by calling appropriate object
		netin.SubCidr = sub
		netin.Name = net.Name + "_sub" + strconv.Itoa(i)
		netin.Zone = zones[zonenum]
		if net.GetRaw == true {
			netin.VpcId = *vpc.CreateVpcRaw.Vpc.VpcId
			netin.IgwId = *vpc.CreateIgwRaw.InternetGateway.InternetGatewayId
		} else {
			netin.VpcId = vpc.VpcId
			netin.IgwId = vpc.IgwId
		}

		subnet, sub_err := netin.CreateSubnet(con)
		if sub_err != nil {
			return NetworkResponse{}, sub_err
		}
		subnets = append(subnets, subnet)

		zonenum--
	}
	if net.GetRaw == true {
		return NetworkResponse{CreateVpcRaw: vpc, CreateSubnetRaw: subnets}, nil
	}
	return NetworkResponse{Name: vpc.Name, VpcId: vpc.VpcId, Subnets: subnets, Type: vpc.Type, IgwId: vpc.IgwId, SecGroupIds: vpc.SecGroupIds}, nil

}

func (d *DeleteNetworkInput) DeleteNetwork(con aws.EstablishConnectionInput) (DeleteNetworkResponse, error) {

	vpcin := GetNetworksInput{VpcIds: d.VpcIds}
	vpc, err := vpcin.FindVpcs(con)
	if err != nil {
		return DeleteNetworkResponse{}, err
	}

	if vpc != true {
		return DeleteNetworkResponse{}, fmt.Errorf("Could not find the entered VPC, please enter valid/existing VPC id")
	}

	networkdel, neterr := d.getNetworkDeletables(con)
	if neterr != nil {
		return DeleteNetworkResponse{}, neterr
	}

	deletestatus, netdelerr := networkdel.deleteNetworkDeletables(con)
	if netdelerr != nil {
		return DeleteNetworkResponse{}, netdelerr
	}
	return deletestatus, nil
}

func (d *DeleteNetworkInput) deleteNetworkDeletables(con aws.EstablishConnectionInput) (DeleteNetworkResponse, error) {

	/*ec2, seserr := con.EstablishConnection()
	  if seserr != nil {
	          return DeleteNetworkResponse{}, seserr
	  }*/

	if len(d.SecIds) != 0 {
		//Deletion of security groups
		delsecin := NetworkComponentInput{SecGroupIds: d.SecIds}
		secdelerr := delsecin.DeleteSecutiryGroup(con)
		if secdelerr != nil {
			return DeleteNetworkResponse{}, secdelerr
		}
	}

	if len(d.RouteTableIds) != 0 {
		//dissassociating the routetable before deleteing it.
		route := NetworkComponentInput{RouteTableIds: d.RouteTableIds}
		dessroutetable, dessrouterr := route.DisassociateRouteTable(con)
		if dessrouterr != nil {
			return DeleteNetworkResponse{}, dessrouterr
		}
		if dessroutetable != true {
			return DeleteNetworkResponse{}, fmt.Errorf("An error occured while dettaching routetable from subnet")
		}

		//deletion of routetable is handled by below loop.
		delrouterr := route.DeleteRouteTable(con)
		if delrouterr != nil {
			return DeleteNetworkResponse{}, delrouterr
		}
	}

	if len(d.IgwIds) != 0 {
		//dettachment of igw is been done by below snippet.
		dettach_gateway := NetworkComponentInput{IgwIds: d.IgwIds, VpcIds: d.VpcIds}
		detacherr := dettach_gateway.DetachIgws(con)
		if detacherr != nil {
			return DeleteNetworkResponse{}, detacherr
		}

		//deletion of igw is been done by below snippet.
		delete_gateway := NetworkComponentInput{IgwIds: d.IgwIds}
		deleteigwerr := delete_gateway.DeleteIgws(con)
		if deleteigwerr != nil {
			return DeleteNetworkResponse{}, deleteigwerr
		}
	}

	if len(d.SubnetIds) != 0 {
		subdelin := DeleteNetworkInput{SubnetIds: d.SubnetIds}
		subdelerr := subdelin.DeleteSubnets(con)
		if subdelerr != nil {
			return DeleteNetworkResponse{}, subdelerr
		}
	}

	//deletion of vpc is handled by below snippet
	delete_vpc := DeleteNetworkInput{VpcIds: d.VpcIds}
	deletevpcerr := delete_vpc.DeleteVpc(con)
	if deletevpcerr != nil {
		return DeleteNetworkResponse{}, deletevpcerr
	}
	return DeleteNetworkResponse{Status: "Network and all its components has been deleted successfully"}, nil
}

func (d *DeleteNetworkInput) getNetworkDeletables(con aws.EstablishConnectionInput) (DeleteNetworkInput, error) {

	//creating a session to perform actions
	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return DeleteNetworkInput{}, seserr
	}

	//Getting list of all subnets available in the network
	subnetres, suberr := ec2.DescribeSubnet(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: d.VpcIds,
			},
		},
	)
	if suberr != nil {
		return DeleteNetworkInput{}, suberr
	}

	subnets := make([]string, 0)
	for _, subnet := range subnetres.Subnets {
		if *subnet.DefaultForAz != true {
			subnets = append(subnets, *subnet.SubnetId)
		}
	}

	//Getting list of all secutiry groups in the entered vpc.
	secres, secerr := ec2.DescribeSecurityGroup(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: d.VpcIds,
			},
		},
	)
	if secerr != nil {
		return DeleteNetworkInput{}, secerr
	}
	sec_ids := make([]string, 0)
	for _, sec := range secres.SecurityGroups {
		if *sec.GroupName != "default" {
			sec_ids = append(sec_ids, *sec.GroupId)
		}
	}

	//describing all the routetables to fetch the right ones.
	routres, routerr := ec2.DescribeRouteTable(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: d.VpcIds,
			},
		},
	)
	if routerr != nil {
		return DeleteNetworkInput{}, routerr
	}
	route_ids := make([]string, 0)
	for _, route := range routres.RouteTables {
		if route.Associations != nil {
			if *route.Associations[0].Main != true {
				route_ids = append(route_ids, *route.RouteTableId)
			}
		} else {
			route_ids = append(route_ids, *route.RouteTableId)
		}
	}

	//describing all internet-gateways to get right one.
	response, err := ec2.DescribeIgw(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "attachment.vpc-id",
				Value: d.VpcIds,
			},
		},
	)
	if err != nil {
		return DeleteNetworkInput{}, err
	}
	igw_ids := make([]string, 0)
	for _, igw := range response.InternetGateways {
		igw_ids = append(igw_ids, *igw.InternetGatewayId)
	}

	//collating the data of entire network which was collected.
	deleteResponse := new(DeleteNetworkInput)
	deleteResponse.SubnetIds = subnets
	deleteResponse.SecIds = sec_ids
	deleteResponse.RouteTableIds = route_ids
	deleteResponse.IgwIds = igw_ids
	deleteResponse.VpcIds = d.VpcIds

	return *deleteResponse, nil

}

func (net *GetNetworksInput) GetNetwork(con aws.EstablishConnectionInput) ([]NetworkResponse, error) {

	network_response := make([]NetworkResponse, 0)
	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return nil, seserr
	}

	find_vpc_result, vpc_err := ec2.DescribeVpc(
		&aws.DescribeNetworkInput{
			VpcIds: net.VpcIds,
		},
	)

	if vpc_err != nil {
		return nil, vpc_err
	}

	for _, vpc := range find_vpc_result.Vpcs {

		// getting list of subnets in network
		subnetin := GetNetworksInput{VpcIds: []string{*vpc.VpcId}, GetRaw: net.GetRaw}
		subnets, suberr := subnetin.GetSubnetsFromVpc(con)
		if suberr != nil {
			return nil, suberr
		}

		// getting list of igws in network
		igwin := NetworkComponentInput{VpcIds: []string{*vpc.VpcId}, GetRaw: net.GetRaw}
		igw, igwerr := igwin.GetIgwFromVpc(con)
		if igwerr != nil {
			return nil, igwerr
		}

		// getting list of security group in network
		sec, secerr := igwin.GetSecFromVpc(con)
		if secerr != nil {
			return nil, secerr
		}

		if net.GetRaw == true {
			subnets.GetVpcRaw = vpc
			subnets.DescribeSecurityRaw = sec.GetSecurityRaw
			subnets.DescribeIgwRaw = igw.GetIgwRaw
			network_response = append(network_response, subnets)
		} else {
			if vpc.Tags != nil {
				network_response = append(network_response, NetworkResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, Subnets: subnets.Subnets, State: *vpc.State, IgwId: igw.IgwIds[0], IsDefault: *vpc.IsDefault, SecGroupIds: sec.SecGroupIds})
			} else {
				network_response = append(network_response, NetworkResponse{VpcId: *vpc.VpcId, Subnets: subnets.Subnets, State: *vpc.State, IgwId: igw.IgwIds[0], IsDefault: *vpc.IsDefault, SecGroupIds: sec.SecGroupIds})
			}
		}
	}
	return network_response, nil
}

func (net *GetNetworksInput) GetAllNetworks(con aws.EstablishConnectionInput) ([]NetworkResponse, error) {

	network_response := make([]NetworkResponse, 0)
	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return nil, seserr
	}

	find_vpc_result, vpc_err := ec2.DescribeAllVpc(
		&aws.DescribeNetworkInput{},
	)
	if vpc_err != nil {
		return nil, vpc_err
	}
	for _, vpc := range find_vpc_result.Vpcs {

		// getting list of subnets in network
		subnetin := GetNetworksInput{VpcIds: []string{*vpc.VpcId}, GetRaw: net.GetRaw}
		subnets, suberr := subnetin.GetSubnetsFromVpc(con)
		if suberr != nil {
			return nil, suberr
		}

		// getting list of igws in network
		igwin := NetworkComponentInput{VpcIds: []string{*vpc.VpcId}, GetRaw: net.GetRaw}
		igw, igwerr := igwin.GetIgwFromVpc(con)
		if igwerr != nil {
			return nil, igwerr
		}

		// getting list of security group in network
		sec, secerr := igwin.GetSecFromVpc(con)
		if secerr != nil {
			return nil, secerr
		}

		if net.GetRaw == true {
			netres := new(NetworkResponse)
			netres.GetVpcRaw = vpc
			netres.GetSubnetRaw = subnets.GetSubnetRaw
			netres.DescribeSecurityRaw = sec.GetSecurityRaw
			netres.DescribeIgwRaw = igw.GetIgwRaw
			network_response = append(network_response, *netres)
		} else {
			if vpc.Tags != nil {
				network_response = append(network_response, NetworkResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, Subnets: subnets.Subnets, State: *vpc.State, IgwId: igw.IgwIds[0], SecGroupIds: sec.SecGroupIds, IsDefault: *vpc.IsDefault, Region: con.Region})
			} else {
				network_response = append(network_response, NetworkResponse{VpcId: *vpc.VpcId, Subnets: subnets.Subnets, State: *vpc.State, IgwId: igw.IgwIds[0], SecGroupIds: sec.SecGroupIds, IsDefault: *vpc.IsDefault, Region: con.Region})
			}
		}
	}
	return network_response, nil
}

func (net *UpdateNetworkInput) UpdateNetwork(con aws.EstablishConnectionInput) (NetworkResponse, error) {

	/*ec2, seserr := con.EstablishConnection()
	  if seserr != nil {
	          return NetworkResponse{}, seserr
	  }*/

	switch strings.ToLower(net.Resource) {
	case "subnets":

		switch strings.ToLower(net.Action) {
		case "create":
			// Collects all the available availability zones
			zones := make([]string, 0)
			if net.Network.Zone == "" {
				zonein := common.CommonInput{}
				zone, zon_err := zonein.GetAvailabilityZones(con)
				if zon_err != nil {
					return NetworkResponse{}, zon_err
				}
				zones = zone
			} else {
				zones = []string{net.Network.Zone}
			}
			// I will be the spoc for subnets creation in the loop as per the request made
			subnet_response := make([]SubnetReponse, 0)
			zonenum := len(zones) - 1

			//Fetching unique number to give our subnet a unique name
			subnets := make([]string, 0)
			subnetin := GetNetworksInput{VpcIds: []string{net.Network.VpcId}}
			subnetlist, suberr := subnetin.GetSubnetsFromVpc(con)
			if suberr != nil {
				return NetworkResponse{}, suberr
			}
			for _, subnet := range subnetlist.Subnets {
				if subnet.Name != "" {
					subnets = append(subnets, subnet.Name)
				}
			}

			// Getting Unique digit to name subnet uniquly
			uqnin := common.CommonInput{SortInput: subnets}
			uqnchr, unerr := uqnin.GetUniqueNumberFromTags()
			if unerr != nil {
				return NetworkResponse{}, unerr
			}
			for _, sub := range net.Network.SubCidrs {
				if zonenum < 0 {
					zonenum = len(zones) - 1
				}

				// Creating subnet by calling appropriate object
				subin := NetworkCreateInput{
					SubCidr: sub,
					Name:    net.Network.Name + "_sub" + strconv.Itoa(uqnchr),
					Zone:    zones[zonenum],
					VpcId:   net.Network.VpcId,
					GetRaw:  net.GetRaw,
				}
				subnet, sub_err := subin.CreateSubnet(con)
				if sub_err != nil {
					return NetworkResponse{}, sub_err
				}
				subnet_response = append(subnet_response, subnet)

				zonenum--
				uqnchr++
			}
			if net.GetRaw == true {
				return NetworkResponse{CreateSubnetRaw: subnet_response}, nil
			}
			return NetworkResponse{Subnets: subnet_response}, nil
		case "delete":
			return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting the action %s of the resource %s or you entered wrong name. The action you selected was: %s", net.Action, net.Resource, net.Action))
		default:
			return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting the action %s of the resource %s or you entered wrong name. The action you selected was: %s", net.Action, net.Resource, net.Action))
		}

	case "vpc":
		return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting updation of the resource you entered or you entered wrong name. The resource you enetered was: %s", net.Resource))
	case "igw":
		return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting updation of the resource you entered or you entered wrong name. The resource you enetered was: %s", net.Resource))
	default:
		return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting updation of the resource you entered or you entered wrong name. The resource you enetered was: %s", net.Resource))
	}
}
