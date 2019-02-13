package neuronaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type CreateSessionInput struct {
	Region   string
	KeyId    string
	AcessKey string
}

func (con *CreateSessionInput) CreateAwsSession() *session.Session {

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentialsFromCreds(
			credentials.Value{
				AccessKeyID:     con.KeyId,
				SecretAccessKey: con.AcessKey,
			}),
		Region: aws.String(con.Region),
	}))
	return sess
}
