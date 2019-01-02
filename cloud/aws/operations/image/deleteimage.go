package awsimage

import (
	aws "neuron/cloud/aws/interface"
	err "neuron/error"
)

type DeleteImageInput struct {
	ImageIds []string
}

func (img *DeleteImageInput) DeleteImage(con aws.EstablishConnectionInput) (ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return ImageResponse{}, seserr
	}

	// desribing image to check if image exists
	search_image := GetImageInput{ImageIds: img.ImageIds}
	imagexists, des_img_err := search_image.IsImageAvailable(con)

	if des_img_err != nil {
		return ImageResponse{}, des_img_err
	}

	if imagexists != true {
		return ImageResponse{}, err.ImageNotFound()
	}

	image_result, image_err := ec2.DescribeImages(
		&aws.DescribeComputeInput{
			ImageIds: img.ImageIds,
		},
	)
	if image_err != nil {
		return ImageResponse{}, image_err
	}

	for _, image := range image_result.Images {

		// Deregistering image will be done by following code
		der_err := ec2.DeregisterImage(
			&aws.DeleteComputeInput{
				ImageId: *image.ImageId,
			},
		)

		if der_err != nil {
			return ImageResponse{}, der_err
		}

		// Deletion of snapshot will addressed by following code
		snap_err := ec2.DeleteSnapshot(
			&aws.DeleteComputeInput{
				SnapshotId: *image.BlockDeviceMappings[0].Ebs.SnapshotId,
			},
		)

		if snap_err != nil {
			return ImageResponse{}, snap_err
		}
	}
	return ImageResponse{DeleteResponse: "Image is successfully deleted"}, nil
}
