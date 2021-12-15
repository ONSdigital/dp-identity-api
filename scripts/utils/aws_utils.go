package utils

import (
	"context"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"os"
)

type S3Reader struct {
}

func (s S3Reader) GetS3Reader(ctx context.Context, region string, s3Bucket string, s3FilePath string) io.ReadCloser {
	s3Client := s3.New(session.Must(session.NewSession(&aws.Config{Region: aws.String(region)})))
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

func GetCognitoClient(region string) *cognito.CognitoIdentityProvider {
	client := cognito.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})), &aws.Config{Region: &region})
	return client
}
