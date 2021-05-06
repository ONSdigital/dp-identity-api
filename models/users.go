package models

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

type UserParams struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}
type CreateUserOutput struct {
	*cognitoidentityprovider.AdminCreateUserOutput
}
