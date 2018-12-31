package awsimage

import (
	aws "neuron/cloud/aws/interface"
)

type GetImageInput struct {
	Kind     string   `json:"Kind"`
	ImageIds []string `json:"ImageIds"`
	GetRaw   bool     `json:"GetRaw"`
}

//This function is tuned to get the details of the images, who's Id will be passed to it.
func (i *GetImageInput) GetImage(con aws.EstablishConnectionInput) ([]ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return nil, seserr
	}

	// desribing image to check if image exists
	image_result, image_err := ec2.DescribeImages(
		&aws.DescribeComputeInput{
			ImageIds: i.ImageIds,
		},
	)
	if image_err != nil {
		return nil, image_err
	}

	image_response := make([]ImageResponse, 0)
	for _, img := range image_result.Images {
		if i.GetRaw == true {
			image_response = append(image_response, ImageResponse{GetImageRaw: img})
		} else {
			resp := new(ImageResponse)
			resp.Name = *img.Name
			resp.ImageId = *img.ImageId
			resp.CreationDate = *img.CreationDate
			resp.State = *img.State
			resp.IsPublic = *img.Public

			snap := new(SnapshotDetails)
			snap.SnapshotId = *img.BlockDeviceMappings[0].Ebs.SnapshotId
			snap.VolumeType = *img.BlockDeviceMappings[0].Ebs.VolumeType
			snap.VolumeSize = *img.BlockDeviceMappings[0].Ebs.VolumeSize
			resp.SnapShot = *snap
			image_response = append(image_response, *resp)
		}
	}
	return image_response, nil
}

// This function is tuned to get the details of all images present under this account in the entered region.
func (i *GetImageInput) GetAllImage(con aws.EstablishConnectionInput) ([]ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return nil, seserr
	}

	// desribing image to check if image exists
	result, err := ec2.DescribeAllImages(
		&aws.DescribeComputeInput{},
	)

	if err != nil {
		return nil, err
	}

	image_response := make([]ImageResponse, 0)
	for _, img := range result.Images {
		if i.GetRaw == true {
			image_response = append(image_response, ImageResponse{GetImageRaw: img})
		} else {
			resp := new(ImageResponse)
			resp.Name = *img.Name
			resp.ImageId = *img.ImageId
			resp.CreationDate = *img.CreationDate
			resp.State = *img.State
			resp.IsPublic = *img.Public

			snap := new(SnapshotDetails)
			snap.SnapshotId = *img.BlockDeviceMappings[0].Ebs.SnapshotId
			snap.VolumeType = *img.BlockDeviceMappings[0].Ebs.VolumeType
			snap.VolumeSize = *img.BlockDeviceMappings[0].Ebs.VolumeSize
			resp.SnapShot = *snap
			image_response = append(image_response, *resp)
		}
	}
	return image_response, nil
}
