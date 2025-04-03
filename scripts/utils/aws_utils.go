package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Reader struct {
}

func (s S3Reader) GetS3Reader(ctx context.Context, awsProfile, region, s3Bucket, s3FilePath string) io.ReadCloser {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(awsProfile),
	)

	if err != nil {
		log.Fatal(ctx, "unable to load the SDK", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3FilePath),
	}
	objectOutput, err := s3Client.GetObject(ctx, objectInput)
	if err != nil {
		log.Error(ctx, "failed to read groups file from s3", err)
		os.Exit(1)
	}
	return objectOutput.Body
}

func GetCognitoClient(ctx context.Context, awsProfile, region string) *cognitoidentityprovider.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(awsProfile),
	)

	if err != nil {
		log.Fatal(ctx, "unable to load the SDK", err)
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)

	return client
}
