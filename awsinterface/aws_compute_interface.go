package DengineAwsInterface

import (
	b64 "encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "neuron/logger"
	"strings"
	"time"
)

type ImageCreateInput struct {
	InstanceId string
}

type CreateServerInput struct {
	InstanceName string
	ImageId      string
	InstanceType string
	KeyName      string
	MaxCount     int64 `json:"MaxCount,omitempty"`
	MinCount     int64 `json:"MinCount,omitempty"`
	SubnetId     string
	UserData     string `json:"Name,omitempty"`
	AssignPubIp  bool   `json:"Name,omitempty"`
}

type ImageResponse struct {
	Name            string          `json:"Name,omitempty"`
	ImageId         string          `json:"ImageId,omitempty"`
	State           string          `json:"State,omitempty"`
	IsPublic        bool            `json:"IsPublic,omitempty"`
	CreationDate    string          `json:"CreationDate,omitempty"`
	Description     string          `json:"Description,omitempty"`
	DefaultResponse string          `json:"DefaultResponse,omitempty"`
	DeleteResponse  string          `json:"ImageResponse,omitempty"`
	SnapShot        SnapshotDetails `json:"SnapShot,omitempty"`
}

type SnapshotDetails struct {
	SnapshotId string `json:"SnapshotId,omitempty"`
	VolumeType string `json:"VolumeType,omitempty"`
	VolumeSize int64  `json:"VolumeSize,omitempty"`
}

type ServerResponse struct {
	InstanceName        string      `json:"InstanceName,omitempty"`
	InstanceId          string      `json:"InstanceId,omitempty"`
	SubnetId            string      `json:"SubnetId,omitempty"`
	PrivateIpAddress    string      `json:"IpAddress,omitempty"`
	PublicIpAddress     string      `json:"PublicIpAddress,omitempty"`
	PrivateDnsName      string      `json:"PrivateDnsName,omitempty"`
	CreatedOn           string      `json:"CreatedOn,omitempty"`
	State               string      `json:"State,omitempty"`
	InstanceDeleteState string      `json:"InstanceDeleteState,omitempty"`
	InstanceType        string      `json:"InstanceType,omitempty"`
	Cloud               string      `json:"Cloud,omitempty"`
	Region              string      `json:"Region,omitempty"`
	PreviousState       string      `json:"PreviousState,omitempty"`
	CurrentState        string      `json:"CurrentState,omitempty"`
	DefaultResponse     interface{} `json:"DefaultResponse,omitempty"`
	Error               string      `json:"Error,omitempty"`
}

type DeleteServerFromVpcInput struct {
	VpcId string
}

type UpdateServerInput struct {
	InstanceIds        []string `json:"InstanceIds,omitempty"`
	StartStopInstances string   `json:"StartStopInstances,omitempty"`
}

type DeleteImageInput struct {
	ImageId string
}

type GetImageInput struct {
	ImageId string
}

func (csrv *CreateServerInput) CreateServer() (ServerResponse, error) {

	var create_server_response ServerResponse
	// I will make a decision which security group to pick
	// fetching the VPC-ID from the subnet which was passed to create instance
	subnet_input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			aws.String(csrv.SubnetId),
		},
	}
	subnet_result, sub_err := Svc.DescribeSubnets(subnet_input)

	if sub_err == nil {

		// Fetching the valid security group
		security_input := &ec2.DescribeSecurityGroupsInput{}
		security_result, sec_err := Svc.DescribeSecurityGroups(security_input)
		var security_id string
		if sec_err != nil {
			return ServerResponse{}, nil
		} else {
			for _, sec := range security_result.SecurityGroups {
				if *sec.VpcId == *subnet_result.Subnets[0].VpcId {
					if *sec.GroupName != "default" {
						security_id = *sec.GroupId
					}
				}
			}

			// I will be the spoc for the instance creation with the userdata passed to me
			var min_count int64
			var max_count int64
			var user_data string

			if csrv.UserData == "" {
				user_data = b64.StdEncoding.EncodeToString([]byte("echo 'nothing'"))
			} else {
				user_data = b64.StdEncoding.EncodeToString([]byte(csrv.UserData))
			}
			if (csrv.MinCount == 0) || (csrv.MaxCount == 0) {
				if csrv.MinCount == 0 {
					min_count = 1
				} else {
					min_count = csrv.MinCount
				}
				if csrv.MaxCount == 0 {
					max_count = 1
				} else {
					max_count = csrv.MaxCount
				}
			} else {
				min_count = csrv.MinCount
				max_count = csrv.MinCount
			}

			// support for custom ebs mapping will be rolled out soon
			create_server_input := &ec2.RunInstancesInput{
				ImageId:      aws.String(csrv.ImageId),
				InstanceType: aws.String(csrv.InstanceType),
				KeyName:      aws.String(csrv.KeyName),
				MaxCount:     aws.Int64(max_count),
				MinCount:     aws.Int64(min_count),
				UserData:     aws.String(user_data),
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{{
					AssociatePublicIpAddress: aws.Bool(csrv.AssignPubIp),
					DeviceIndex:              aws.Int64(0),
					DeleteOnTermination:      aws.Bool(true),
					SubnetId:                 aws.String(csrv.SubnetId),
					Groups:                   []*string{aws.String(security_id)},
				}},
			}
			server_create_result, err := Svc.RunInstances(create_server_input)
			// handling the error if it throws while subnet is under creation process
			if err != nil {
				log.Info("")
				log.Error("This Error is thrown if it encounters error while SERVER creation process is under progress")
				log.Error("Oops....!! we ran into an error while creating server, check log for more info")
				log.Error(err)
				log.Info("")
				return ServerResponse{Error: "Oops....!! we ran into an error while creating server, check log for more info"}, err
			} else {

				var instance_id string
				for _, instance := range server_create_result.Instances {
					instance_id = *instance.InstanceId
				}

				// I will make program wait untill instance become running
				instance_status_input := &ec2.DescribeInstancesInput{
					Filters: []*ec2.Filter{
						&ec2.Filter{
							Name:   aws.String("instance-id"),
							Values: aws.StringSlice([]string{instance_id}),
						},
					},
				}

				Svc.WaitUntilInstanceRunning(instance_status_input)

				instance_input := &ec2.DescribeInstancesInput{
					Filters: []*ec2.Filter{
						&ec2.Filter{
							Name:   aws.String("instance-id"),
							Values: aws.StringSlice([]string{instance_id}),
						},
					},
				}
				instance_describe_result, _ := Svc.DescribeInstances(instance_input)

				type server_response struct {
					instance_id string `json:"instance_id,omitempty"`
					ipaddress   string `json:"ipaddress,omitempty"`
					privatedns  string `json:"privatedns,omitempty"`
					publicIp    string `json:"publicIp,omitempty"`
					createdon   string `json:"createdon,omitempty"`
				}

				var response server_response
				// fetching the instance details which is created in previos process
				for _, reservation := range instance_describe_result.Reservations {
					for _, instance := range reservation.Instances {
						if csrv.AssignPubIp == true {
							response = server_response{instance_id: *instance.InstanceId, ipaddress: *instance.PrivateIpAddress, privatedns: *instance.PrivateDnsName, publicIp: *instance.PublicIpAddress, createdon: (*instance.LaunchTime).String()}
						} else {
							response = server_response{instance_id: *instance.InstanceId, ipaddress: *instance.PrivateIpAddress, privatedns: *instance.PrivateDnsName, createdon: (*instance.LaunchTime).String()}
						}
					}
				}

				// I will pass name to create_tags to set a name to the vpc
				tags := Tag{response.instance_id, "Name", csrv.InstanceName}
				server_tags, _ := create_tags(tags)

				if response.publicIp == "" {
					create_server_response = ServerResponse{InstanceName: server_tags, InstanceId: response.instance_id, SubnetId: csrv.SubnetId, PrivateIpAddress: response.ipaddress, PrivateDnsName: response.privatedns, CreatedOn: response.createdon, Cloud: "Amazon"}
				} else {
					create_server_response = ServerResponse{InstanceName: server_tags, InstanceId: response.instance_id, SubnetId: csrv.SubnetId, PrivateIpAddress: response.ipaddress, PublicIpAddress: response.publicIp, PrivateDnsName: response.privatedns, CreatedOn: response.createdon, Cloud: "Amazon"}
				}
			}
			return create_server_response, nil
		}

	} else {
		return ServerResponse{DefaultResponse: "Oops....!! There is a problem with subnet/network id you eneterd. Either the network does not exists else you would have enetred it wrongly"}, nil
	}
}

func GetServersDetails(sub string) ([]ServerResponse, error) {

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("subnet-id"),
				Values: aws.StringSlice([]string{sub}),
			},
		},
	}

	result, des_inst_err := Svc.DescribeInstances(input)

	if des_inst_err != nil {
		return nil, des_inst_err
	} else {
		server_response := []ServerResponse{}
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				var instance_name string
				for _, inst_name := range instance.Tags {
					instance_name = *inst_name.Value
				}

				t := *instance.LaunchTime
				if *instance.State.Name == "running" {
					server_response = append(server_response, ServerResponse{InstanceName: instance_name, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PublicIpAddress: *instance.PublicIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: t.String(), State: *instance.State.Name, Cloud: "Amazon"})
				} else if *instance.State.Name == "stopped" {
					server_response = append(server_response, ServerResponse{InstanceName: instance_name, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: t.String(), State: *instance.State.Name, Cloud: "Amazon"})
				} else if *instance.State.Name == "terminated" {
					server_response = append(server_response, ServerResponse{State: *instance.State.Name, Cloud: "Amazon"})
				}
			}
		}
		return server_response, nil
	}
}

func GetAllServers() ([]ServerResponse, error) {
	input := &ec2.DescribeInstancesInput{}

	des_serv_result, des_serv_err := Svc.DescribeInstances(input)

	if des_serv_err != nil {
		return nil, des_serv_err
	} else {
		server_response := []ServerResponse{}
		for _, reservation := range des_serv_result.Reservations {
			for _, instance := range reservation.Instances {
				if *instance.State.Name != "terminated" {
					var instance_name string
					for _, inst_name := range instance.Tags {
						if *inst_name.Key == "Name" {
							instance_name = *inst_name.Value
						}
					}
					server_response = append(server_response, ServerResponse{InstanceName: instance_name, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, InstanceType: *instance.InstanceType, Cloud: "Amazon", Region: *instance.Placement.AvailabilityZone})

				} else {
					// change has to be made here (introduction of omitempty is required)
					server_response = append(server_response, ServerResponse{State: "terminated", Cloud: "Amazon"})
				}
			}
		}
		return server_response, nil
	}
}

func DeleteServer(instanceId []string) ([]ServerResponse, error) {

	search_instance_input := &ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice(instanceId)}
	_, search_err := Svc.DescribeInstances(search_instance_input)

	if search_err != nil {
		return nil, search_err
	} else {
		terminate_instance_input := &ec2.TerminateInstancesInput{
			InstanceIds: aws.StringSlice(instanceId),
		}

		_, ins_term_err := Svc.TerminateInstances(terminate_instance_input)

		if ins_term_err != nil {
			return nil, ins_term_err
		} else {
			instance_status_input := &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name:   aws.String("instance-id"),
						Values: aws.StringSlice(instanceId),
					},
				},
			}

			Svc.WaitUntilInstanceTerminated(instance_status_input)

			input := &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name:   aws.String("instance-id"),
						Values: aws.StringSlice(instanceId),
					},
				},
			}

			result, _ := Svc.DescribeInstances(input)

			var delete_response []ServerResponse
			for _, reservation := range result.Reservations {
				for _, instance := range reservation.Instances {
					delete_response = append(delete_response, ServerResponse{InstanceId: *instance.InstanceId, InstanceDeleteState: *instance.State.Name})
				}
			}
			return delete_response, nil
		}
	}
}

func (u *UpdateServerInput) UpdateServer() ([]ServerResponse, error) {

	var server_response []ServerResponse

	if u.StartStopInstances != "" {

		search_instance_input := &ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice(u.InstanceIds)}
		_, search_err := Svc.DescribeInstances(search_instance_input)

		if search_err != nil {
			server_response = append(server_response, ServerResponse{DefaultResponse: "There is a problem with instance-ids you provided, either instance-ids is wrong or instance does not exists"})
		} else {
			switch strings.ToLower(u.StartStopInstances) {
			case "start":
				start_server_input := &ec2.StartInstancesInput{InstanceIds: aws.StringSlice(u.InstanceIds)}
				start_server_response, start_err := Svc.StartInstances(start_server_input)

				if start_err != nil {
					server_response = append(server_response, ServerResponse{DefaultResponse: "Oops...!!!!. We encountered error while starting instances"})
				} else {
					instance_status_input := &ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice(u.InstanceIds)}
					Svc.WaitUntilInstanceRunning(instance_status_input)

					for _, inst := range start_server_response.StartingInstances {
						server_response = append(server_response, ServerResponse{InstanceId: *inst.InstanceId, CurrentState: "running", PreviousState: *inst.PreviousState.Name})
					}
				}

			case "stop":
				stop_server_input := &ec2.StopInstancesInput{InstanceIds: aws.StringSlice(u.InstanceIds)}
				stop_server_response, stop_err := Svc.StopInstances(stop_server_input)

				if stop_err != nil {
					server_response = append(server_response, ServerResponse{DefaultResponse: "Oops...!!!!. We encountered error while stopping instances"})
				} else {
					instance_status_input := &ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice(u.InstanceIds)}
					Svc.WaitUntilInstanceStopped(instance_status_input)

					for _, inst := range stop_server_response.StoppingInstances {
						server_response = append(server_response, ServerResponse{InstanceId: *inst.InstanceId, CurrentState: "stopped", PreviousState: *inst.PreviousState.Name})
					}
				}

			default:
				server_response = append(server_response, ServerResponse{DefaultResponse: "Sorry...!!!!. I am not aware of the action you asked me to perform, please enter the action which I know"})
			}
		}
		return server_response, nil
	} else {
		return nil, nil
	}
}

func (d *DeleteServerFromVpcInput) DeleteServerFromVpc() ([]ServerResponse, error) {

	var instanceids []string
	server_delete_input := &ec2.DescribeInstancesInput{}
	server_delete_result, serv_des_err := Svc.DescribeInstances(server_delete_input)
	if serv_des_err != nil {
		return nil, serv_des_err
	} else {
		for _, reservation := range server_delete_result.Reservations {
			for _, instance := range reservation.Instances {
				if *instance.State.Name != "terminated" {
					if *instance.VpcId == d.VpcId {
						instanceids = append(instanceids, *instance.InstanceId)
					}
				}
			}
		}

		delete_server, serv_del_err := DeleteServer(instanceids)
		if serv_del_err != nil {
			return nil, serv_del_err
		} else {
			return delete_server, nil
		}
	}
}

func getRegionFromAvail(avai string) (string, error) {

	avai_input := &ec2.DescribeAvailabilityZonesInput{ZoneNames: aws.StringSlice([]string{avai})}
	avai_result, avai_err := Svc.DescribeAvailabilityZones(avai_input)

	if avai_err != nil {
		return "", avai_err
	} else {
		return *avai_result.AvailabilityZones[0].RegionName, nil
	}
}

//.......... Territory of Images starts..........
func FindImageId(kind string) ([]string, error) {

	var image_id []string
	input := &ec2.DescribeImagesInput{Filters: []*ec2.Filter{&ec2.Filter{Name: aws.String("is-public"), Values: aws.StringSlice([]string{"false"})}}}

	result, des_img_err := Svc.DescribeImages(input)

	if des_img_err != nil {
		return nil, des_img_err
	} else {
		for _, image := range result.Images {
			if strings.Contains(*image.Name, kind) {
				image_id = append(image_id, *image.ImageId)
			}
		}
		return image_id, nil
	}
}

// I will be responsible for capturing the image of the server when I am called
func (img *ImageCreateInput) CreateImage() (ImageResponse, error) {

	// fetching instance details as I need to pass this while taking server backup
	var ServerName string
	var inst_search string
	instace_details_input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("instance-id"),
				Values: aws.StringSlice([]string{img.InstanceId}),
			}}}

	instance_details_result, inst_des_err := Svc.DescribeInstances(instace_details_input)
	if inst_des_err != nil {
		return ImageResponse{}, inst_des_err
	} else {
		for _, reservation := range instance_details_result.Reservations {
			for _, instance := range reservation.Instances {
				if img.InstanceId == *instance.InstanceId {
					tags := instance.Tags
					for _, tag := range tags {
						ServerName = *tag.Value
					}
					inst_search = "Found"
				}
			}
		}

		if inst_search == "Found" {
			// Here where do stuff to take server backup
			now_time := time.Now().Local().Format("2006-01-02 09:10:31")
			image_create_input := &ec2.CreateImageInput{
				Description: aws.String("This image is captured by D-Engine api for " + ServerName + " @ " + now_time),
				InstanceId:  aws.String(img.InstanceId),
				Name:        aws.String(ServerName + "-snapshot"),
			}

			image_create_result, err := Svc.CreateImage(image_create_input)
			// handling the error if it throws while subnet is under creation process
			if err != nil {
				return ImageResponse{}, err
			}

			// I will pass name to create_tags to set a name to the Image which I captured just now
			tags := Tag{*image_create_result.ImageId, "Name", ServerName + "-snapshot"}
			image_tag, tag_err := create_tags(tags)
			ver_tags := Tag{*image_create_result.ImageId, "Version", "1"}
			_, ver_tg_err := create_tags(ver_tags)
			if (ver_tg_err != nil) || (tag_err != nil) {
				if ver_tg_err != nil {
					return ImageResponse{}, ver_tg_err
				} else if tag_err != nil {
					return ImageResponse{}, tag_err
				}
			}
			return ImageResponse{Name: image_tag, ImageId: *image_create_result.ImageId, Description: "This image is captured by Neuron api for " + ServerName + " @" + now_time}, nil
		} else {
			return ImageResponse{DefaultResponse: "Could not capture image for entered image, because no instance was found with provided ID: " + img.InstanceId}, nil
		}
	}
}

func (i *DeleteImageInput) DeleteImage() (ImageResponse, error) {

	// desribing image to check if image exists
	search_image_input := &ec2.DescribeImagesInput{Filters: []*ec2.Filter{&ec2.Filter{Name: aws.String("is-public"), Values: aws.StringSlice([]string{"false"})}}}
	search_image_result, des_img_err := Svc.DescribeImages(search_image_input)

	var image_response ImageResponse
	if des_img_err != nil {
		return ImageResponse{}, des_img_err
	} else {

		var search_input string
		var image_status string
		for _, img := range search_image_result.Images {
			if (i.ImageId == *img.ImageId) && (*img.State == "available") {
				search_input = "Found"
				image_status = *img.State
			}
		}

		if search_input == "Found" {

			for _, img := range search_image_result.Images {

				// deregistering image will be done by below code
				deregister_image_input := &ec2.DeregisterImageInput{ImageId: aws.String(i.ImageId)}
				_, der_err := Svc.DeregisterImage(deregister_image_input)

				// Deletion of snapshot will addressed by below code
				snap_delete_input := &ec2.DeleteSnapshotInput{SnapshotId: aws.String(*img.BlockDeviceMappings[0].Ebs.SnapshotId)}
				_, snap_err := Svc.DeleteSnapshot(snap_delete_input)

				if (der_err != nil) || (snap_err != nil) {

					//logging the errors into applog
					if der_err != nil {
						return ImageResponse{DefaultResponse: "Unable to deregister the Image with the Id: " + i.ImageId}, der_err
					} else if snap_err != nil {
						return ImageResponse{DefaultResponse: "Unable to delete the snapshot assosiated the Image with the Id: " + i.ImageId + ". Guess Image is already deleted, if not check log for more errors."}, snap_err
					}
				} else {
					image_response = ImageResponse{DeleteResponse: "Image is successfully deleted"}
				}
			}
			return image_response, nil

		} else {
			return ImageResponse{DefaultResponse: "Feels like image with Id: " + i.ImageId + " does not exists else the current state of image is :" + image_status + " please wait till it reach 'available' state. And also please check the entered ImageId, make sure you enter existing ImangeId . For more detail look into app log"}, nil
		}
	}
}

func (i *GetImageInput) GetImage() (ImageResponse, error) {

	// desribing image to check if image exists
	search_image_input := &ec2.DescribeImagesInput{Filters: []*ec2.Filter{&ec2.Filter{Name: aws.String("is-public"), Values: aws.StringSlice([]string{"false"})}}}
	search_image_result, des_img_err := Svc.DescribeImages(search_image_input)

	var image_response ImageResponse
	if des_img_err != nil {
		return ImageResponse{}, des_img_err
	} else {

		var search_input string
		var image_status string
		for _, img := range search_image_result.Images {
			if (i.ImageId == *img.ImageId) && (*img.State == "available") {
				search_input = "Found"
				image_status = *img.State
			}
		}

		if search_input == "Found" {

			for _, img := range search_image_result.Images {

				var snap_details SnapshotDetails
				for _, snap := range img.BlockDeviceMappings {
					snap_details = SnapshotDetails{*snap.Ebs.SnapshotId, *snap.Ebs.VolumeType, *snap.Ebs.VolumeSize}
				}

				image_response = ImageResponse{Name: *img.Name, ImageId: *img.ImageId, CreationDate: *img.CreationDate, State: *img.State, IsPublic: *img.Public, SnapShot: snap_details}
			}
			return image_response,nil

		} else {
			return ImageResponse{DefaultResponse: "Feels like image with Id: " + i.ImageId + " does not exists else the current state of image is :" + image_status + " please wait till it reach 'available' state. And please check the entered ImageId, please make sure you enter existing ImangeId . For more detail look into app log"}, nil
		}
	}
}

func GetAllImage() ([]ImageResponse, error) {

	var image_response []ImageResponse
	// desribing image to check if image exists
	search_image_input := &ec2.DescribeImagesInput{Filters: []*ec2.Filter{&ec2.Filter{Name: aws.String("is-public"), Values: aws.StringSlice([]string{"false"})}}}
	search_image_result, des_img_err := Svc.DescribeImages(search_image_input)

	if des_img_err != nil {
		return nil, des_img_err
	} else {
		for _, img := range search_image_result.Images {
			snap_details := SnapshotDetails{SnapshotId: *img.BlockDeviceMappings[0].Ebs.SnapshotId, VolumeType: *img.BlockDeviceMappings[0].Ebs.VolumeType, VolumeSize: *img.BlockDeviceMappings[0].Ebs.VolumeSize}

			image_response = append(image_response, ImageResponse{Name: *img.Name, ImageId: *img.ImageId, CreationDate: *img.CreationDate, State: *img.State, IsPublic: *img.Public, SnapShot: snap_details})
		}
		return image_response, nil
	}
}
