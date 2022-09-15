package utils

import (
	"context"
	"io"
	"os"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Reader struct {
}

func (s S3Reader) GetS3Reader(ctx context.Context, awsProfile, region, s3Bucket, s3FilePath string) io.ReadCloser {
	s3Client := s3.New(session.Must(session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
		Config: aws.Config{
			Region: aws.String(region),
		},
		SharedConfigState: session.SharedConfigEnable,
	})))

	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3FilePath),
	}
	objectOutput, err := s3Client.GetObject(objectInput)
	if err != nil {
		log.Error(ctx, "failed to read groups file from s3", err)
		os.Exit(1)
	}
	return objectOutput.Body
}

func GetCognitoClient(awsProfile, region string) *cognito.CognitoIdentityProvider {
	client := cognito.New(session.Must(session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
		Config: aws.Config{
			Region: aws.String(region),
		},
		SharedConfigState: session.SharedConfigEnable,
	})))
	return client
}
