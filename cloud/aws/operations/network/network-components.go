package awsnetwork

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	aws "neuron/cloud/aws/interface"
	common "neuron/cloud/aws/operations/common"
	"strconv"
	"strings"
)

type NetworkComponentInput struct {
	Name            string   `json:"Name"`
	VpcIds          []string `json:"VpcId"`
	SubId           string   `json:"SubId"`
	IgwId           string   `json:"IgwId"`
	IgwIds          []string `json:"IgwIds"`
	SubType         string   `json:"SubType"`
	Ports           []string `json:"Ports"`
	Filters         Filters  `json:"Filters"`
	SecGroupIds     []string `json:"SecGroupIds"`
	RouteTableIds   []string `json:"RouteTableIds"`
	DestinationCidr string   `json:"DestinationCidr"`
	GetRaw          bool     `json:"GetRaw"`
}

type NetworkComponentResponse struct {
	IgwIds            []string                            `json:"IgwId,omitempty"`
	SecGroupIds       []string                            `json:"SecGroupIds,omitempty"`
	RouteTableIds     []string                            `json:"RouteTableIds,omitempty"`
	CreateIgwRaw      *ec2.CreateInternetGatewayOutput    `json:"CreateIgwRaw,omitempty"`
	GetIgwRaw         *ec2.DescribeInternetGatewaysOutput `json:"GetIgwRaw,omitempty"`
	CreateSecurityRaw *ec2.CreateSecurityGroupOutput      `json:"CreateSecRaw,omitempty"`
	GetRouteTableRaw  *ec2.DescribeRouteTablesOutput      `json:"DescribeRouteTableRaw,omitempty"`
	GetSecurityRaw    *ec2.DescribeSecurityGroupsOutput   `json:"DescribeSecurityRaw,omitempty"`
}

//This is customized internet-gateway creation, if one needs plain internet-gateway creation he/she has call interface the GOD which talks to cloud.
func (igw *NetworkComponentInput) CreateIgw(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	ig, ig_err := ec2.CreateIgw()
	if ig_err != nil {
		return NetworkComponentResponse{}, ig_err
	}

	if igw.VpcIds != nil {
		at_err := ec2.AttachIgw(
			&aws.DescribeNetworkInput{
				IgwIds: []string{*ig.InternetGateway.InternetGatewayId},
				VpcIds: igw.VpcIds,
			},
		)
		if at_err != nil {
			return NetworkComponentResponse{}, at_err
		}
	}

	igtags := common.Tag{*ig.InternetGateway.InternetGatewayId, "Name", igw.Name + "_igw"}
	_, igtag_err := igtags.CreateTags(con)
	if igtag_err != nil {
		return NetworkComponentResponse{}, igtag_err
	}

	if igw.GetRaw == true {
		return NetworkComponentResponse{CreateIgwRaw: ig}, nil
	}
	return NetworkComponentResponse{IgwIds: []string{*ig.InternetGateway.InternetGatewayId}}, nil
}

func (i *NetworkComponentInput) GetIgwFromVpc(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	response, err := ec2.DescribeIgw(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "attachment.vpc-id",
				Value: i.VpcIds,
			},
		},
	)
	if err != nil {
		return NetworkComponentResponse{}, err
	}

	if i.GetRaw == true {
		return NetworkComponentResponse{GetIgwRaw: response}, nil
	}

	igw_ids := make([]string, 0)
	for _, igw := range response.InternetGateways {
		igw_ids = append(igw_ids, *igw.InternetGatewayId)
	}

	return NetworkComponentResponse{IgwIds: igw_ids}, nil
}

//This is customized security-group creation, if one needs plain security-group creation he/she has call interface the GOD which talks to cloud.
func (sec *NetworkComponentInput) CreateSecurityGroup(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	security, sec_err := ec2.CreateSecurityGroup(
		&aws.CreateNetworkInput{
			VpcId: sec.VpcIds[0],
			Name:  sec.Name + "_sec",
		},
	)
	if sec_err != nil {
		return NetworkComponentResponse{}, sec_err
	}

	sctags := common.Tag{*security.GroupId, "Name", sec.Name + "_sec"}
	_, sctag_err := sctags.CreateTags(con)
	if sctag_err != nil {
		return NetworkComponentResponse{}, sctag_err
	}

	//creating egree and ingres rules for the security group which I created just now
	for _, port := range sec.Ports {
		int_port, _ := strconv.ParseInt(port, 10, 64)
		ingres_err := ec2.CreateIngressRule(
			&aws.IngressEgressInput{
				Port:  int_port,
				SecId: *security.GroupId,
			},
		)
		if ingres_err != nil {
			return NetworkComponentResponse{}, ingres_err
		}
	}
	egres_err := ec2.CreateEgressRule(
		&aws.IngressEgressInput{
			SecId: *security.GroupId,
		},
	)
	if egres_err != nil {
		return NetworkComponentResponse{}, egres_err
	}

	if sec.GetRaw == true {
		return NetworkComponentResponse{CreateSecurityRaw: security}, nil
	}
	return NetworkComponentResponse{SecGroupIds: []string{*security.GroupId}}, nil
}

func (s *NetworkComponentInput) GetSecFromVpc(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	response, err := ec2.DescribeSecurityGroup(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: s.VpcIds,
			},
		},
	)
	if err != nil {
		return NetworkComponentResponse{}, err
	}
	sec_ids := make([]string, 0)
	if s.GetRaw == true {
		return NetworkComponentResponse{GetSecurityRaw: response}, nil
	}
	for _, sec := range response.SecurityGroups {
		sec_ids = append(sec_ids, *sec.GroupId)
	}
	return NetworkComponentResponse{SecGroupIds: sec_ids}, nil
}

func (s *NetworkComponentInput) DeleteSecutiryGroup(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}
	err := ec2.DeleteSecurityGroup(
		&aws.DescribeNetworkInput{
			SecIds: s.SecGroupIds,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

//This is customized route-table creation, if one needs plain route-table creation he/she has call interface the GOD which talks to cloud.
func (r *NetworkComponentInput) CreateRouteTable(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}
	route_table, routetable_err := ec2.CreateRouteTable(
		&aws.CreateNetworkInput{
			VpcId: r.VpcIds[0],
		},
	)

	if routetable_err != nil {
		return routetable_err
	}

	if r.IgwId != "" {
		if strings.ToLower(r.SubType) == "public" {
			route_err := ec2.WriteRoute(
				&aws.CreateNetworkInput{
					DestinationCidr: "0.0.0.0/0",
					IgwId:           r.IgwId,
					RouteTableId:    *route_table.RouteTable.RouteTableId,
				},
			)
			if route_err != nil {
				return route_err
			}

			route_attach_err := ec2.AttachRouteTable(
				&aws.CreateNetworkInput{
					RouteTableId: *route_table.RouteTable.RouteTableId,
					SubId:        r.SubId,
				},
			)
			if route_attach_err != nil {
				return route_attach_err
			}

			return nil
		} else {
			// Releasing Soon...!!. We are not supporting writing custom routes into route tables as of now.
			return nil
		}
	} else {
		if strings.ToLower(r.SubType) == "public" {
			igws, igw_err := ec2.DescribeAllIgw(
				&aws.DescribeNetworkInput{},
			)
			if igw_err != nil {
				return igw_err
			}
			for _, igw := range igws.InternetGateways {
				if *igw.Attachments[0].VpcId == r.VpcIds[0] {
					route_err := ec2.WriteRoute(
						&aws.CreateNetworkInput{
							DestinationCidr: "0.0.0.0/0",
							IgwId:           *igw.InternetGatewayId,
							RouteTableId:    *route_table.RouteTable.RouteTableId,
						},
					)
					if route_err != nil {
						return route_err
					}

					route_attach_err := ec2.AttachRouteTable(
						&aws.CreateNetworkInput{
							RouteTableId: *route_table.RouteTable.RouteTableId,
							SubId:        r.SubId,
						},
					)
					if route_attach_err != nil {
						return route_attach_err
					}

					return nil
				}
			}

		} else {
			// Releasing Soon...!!. We are not supporting writing custom routes into route tables as of now.
			return nil
		}
		return nil
	}
}

func (d *NetworkComponentInput) DisassociateRouteTable(con aws.EstablishConnectionInput) (bool, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return false, seserr
	}

	response, reserr := ec2.DescribeRouteTable(
		&aws.DescribeNetworkInput{
			RouteTableIds: d.RouteTableIds,
		},
	)
	if reserr != nil {
		return false, reserr
	}

	associationId := make([]string, 0)
	for _, routetable := range response.RouteTables {
		if routetable.Associations != nil {
			associationId = append(associationId, *routetable.Associations[0].RouteTableAssociationId)
		}
	}

	for _, id := range associationId {
		deterr := ec2.DettachRouteTable(
			&aws.DescribeNetworkInput{
				AssociationsId: id,
			},
		)
		if deterr != nil {
			return false, deterr
		}
	}
	return true, nil
}

func (s *NetworkComponentInput) GetRouteTableFromVpc(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	response, err := ec2.DescribeRouteTable(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: s.VpcIds,
			},
		},
	)
	if err != nil {
		return NetworkComponentResponse{}, err
	}

	route_ids := make([]string, 0)

	if s.GetRaw == true {
		return NetworkComponentResponse{GetRouteTableRaw: response}, nil
	}

	for _, route := range response.RouteTables {
		route_ids = append(route_ids, *route.RouteTableId)
	}
	return NetworkComponentResponse{RouteTableIds: route_ids}, nil
}

func (s *NetworkComponentInput) DeleteRouteTable(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}

	for _, route := range s.RouteTableIds {
		err := ec2.DeleteRouteTable(
			&aws.DescribeNetworkInput{
				RouteTableIds: []string{route},
			},
		)
		if err != nil {
			return err
		}

		//Waiting till reoutetables deletion is successfully completed
		routewait, routwaiterr := ec2.WaitUntilRoutTableDeleted(
			&aws.DescribeNetworkInput{
				RouteTableIds: []string{route},
			},
		)
		if routwaiterr != nil {
			return routwaiterr
		}
		if routewait == false {
			return fmt.Errorf("An error occured while deleting a routetable")
		}
	}
	return nil
}

func (s *NetworkComponentInput) DetachIgws(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}
	for _, igw := range s.IgwIds {
		err := ec2.DetachIgw(
			&aws.DescribeNetworkInput{
				IgwIds: []string{igw},
				VpcIds: s.VpcIds,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *NetworkComponentInput) DeleteIgws(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}

	for _, igw := range i.IgwIds {
		err := ec2.DeleteIgw(
			&aws.DescribeNetworkInput{
				IgwIds: []string{igw},
				VpcIds: i.VpcIds,
			},
		)
		if err != nil {
			return err
		}

		//Waiting till internetgateways deletion is successfully completed
		igwwait, igwwaiterr := ec2.WaitUntilIgwDeleted(
			&aws.DescribeNetworkInput{
				IgwIds: []string{igw},
			},
		)
		if igwwaiterr != nil {
			return igwwaiterr
		}
		if igwwait == false {
			return fmt.Errorf("An error occured while deleting a igws")
		}
	}
	return nil
}
