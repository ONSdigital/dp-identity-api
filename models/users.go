package models

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

type UserParams struct {
	Forename string `json:"forename"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
}
type CreateUserInput struct {
	UserInput *cognitoidentityprovider.AdminCreateUserInput
}
type CreateUserOutput struct {
	UserOutput *cognitoidentityprovider.AdminCreateUserOutput
}
type ListUsersInput struct {
	ListUsersInput *cognitoidentityprovider.ListUsersInput
}
type ListUsersOutput struct {
	ListUsersOutput *cognitoidentityprovider.ListUsersOutput
}
