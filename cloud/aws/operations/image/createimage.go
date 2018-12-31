package awsimage

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	aws "neuron/cloud/aws/interface"
	common "neuron/cloud/aws/operations/common"
	server "neuron/cloud/aws/operations/server"
	err "neuron/error"
	"strconv"
	"time"
)

type ImageCreateInput struct {
	InstanceId string
	GetRaw     bool
}

type ImageResponse struct {
	Name            string                    `json:"Name,omitempty"`
	ImageId         string                    `json:"ImageId,omitempty"`
	ImageIds        []string                  `json:"ImageIds,omitempty"`
	State           string                    `json:"State,omitempty"`
	IsPublic        bool                      `json:"IsPublic,omitempty"`
	CreationDate    string                    `json:"CreationDate,omitempty"`
	Description     string                    `json:"Description,omitempty"`
	DefaultResponse string                    `json:"DefaultResponse,omitempty"`
	DeleteResponse  string                    `json:"ImageResponse,omitempty"`
	SnapShot        SnapshotDetails           `json:"SnapShot,omitempty"`
	CreateImageRaw  *ec2.CreateImageOutput    `json:"CreateImageRaw,omitempty"`
	GetImagesRaw    *ec2.DescribeImagesOutput `json:"GetImagesRaw,omitempty"`
	GetImageRaw     *ec2.Image                `json:"GetImageRaw,omitempty"`
}

type SnapshotDetails struct {
	SnapshotId string `json:"SnapshotId,omitempty"`
	VolumeType string `json:"VolumeType,omitempty"`
	VolumeSize int64  `json:"VolumeSize,omitempty"`
}

// I will be responsible for capturing the image of the server when I am called
func (img *ImageCreateInput) CreateImage(con aws.EstablishConnectionInput) (ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return ImageResponse{}, seserr
	}

	// fetching instance details as I need to pass this while taking server backup
	search_instance := server.CommonComputeInput{InstanceIds: []string{img.InstanceId}}
	instance_result, insterr := search_instance.SearchInstance(con)
	if insterr != nil {
		return ImageResponse{}, insterr
	}
	if instance_result == false {
		return ImageResponse{}, err.ServerNotFound()
	}

	get_instname := server.DescribeInstanceInput{InstanceIds: []string{img.InstanceId}}
	instanceName, instgeterr := get_instname.GetServersDetails(con)
	if instgeterr != nil {
		return ImageResponse{}, instgeterr
	}

	// Here where do stuff to take server backup
	now_time := time.Now().Local().Format("2006-01-02 09:10:31")

	// fetching names from images so that we can name the new image uniquely
	result, deserr := ec2.DescribeAllImages(
		&aws.DescribeComputeInput{},
	)

	if deserr != nil {
		return ImageResponse{}, deserr
	}

	imagenames := make([]string, 0)
	for _, imgs := range result.Images {
		imagenames = append(imagenames, *imgs.Name)
	}

	// Getting Unique number to name image uniquely
	uqnin := common.CommonInput{SortInput: imagenames}
	uqnchr, unerr := uqnin.GetUniqueNumberFromTags()
	if unerr != nil {
		return ImageResponse{}, unerr
	}

	image_create_result, imgerr := ec2.CreateImage(
		&aws.ImageCreateInput{
			Description: "This image is captured by neuron api for " + instanceName[0].InstanceName + " @ " + now_time,
			InstanceId:  img.InstanceId,
			ServerName:  instanceName[0].InstanceName + "-snapshot-" + strconv.Itoa(uqnchr),
		},
	)

	// handling the error if it throws while subnet is under creation process
	if imgerr != nil {
		return ImageResponse{}, imgerr
	}

	// This will take care of creation of primary tags to the image
	tags := common.Tag{*image_create_result.ImageId, "Name", instanceName[0].InstanceName + "-snapshot" + strconv.Itoa(uqnchr)}
	_, tag_err := tags.CreateTags(con)
	if tag_err != nil {
		return ImageResponse{}, tag_err
	}

	/* This will be versioning the images, now this has no much impact but once neuron is built completely this will be helpful
	tags2 := common.Tag{*image_create_result.ImageId, "Version", "1"}
	_, tag2_err := tags2.CreateTags()
	if tag2_err != nil {
		return nil, tag2_err
	}*/

	if img.GetRaw == true {
		return ImageResponse{CreateImageRaw: image_create_result}, nil
	}

	return ImageResponse{Name: instanceName[0].InstanceName + "-snapshot", ImageId: *image_create_result.ImageId, Description: "This image is captured by Neuron api for " + instanceName[0].InstanceName + " @ " + now_time}, nil
}
