package cognito

import (
	"context"
	"errors"

	"github.com/ONSdigital/log.go/log"

	"github.com/ONSdigital/dp-identity-api/models"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// Client defines an interface for interaction with aws cognitoidentityprovider.
type Client interface {
	DescribeUserPool(*cognito.DescribeUserPoolInput) (*cognito.DescribeUserPoolOutput, error)
	AdminCreateUser(input *cognito.AdminCreateUserInput) (*cognito.AdminCreateUserOutput, error)
}

func AdminCreateUser(ctx context.Context, id string, input *models.CognitoUser) (*cognito.AdminCreateUserOutput, error) {

	log.Event(ctx, "creating user", log.Data{"id": id})

	// Return an error if empty id was passed.
	if id == "" {
		return nil, errors.New("id must not be an empty string")
	}
}
