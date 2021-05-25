package models

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

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

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *AuthParams) ValidateCredentials(ctx context.Context) *[]Error {
	var validationErrors []Error
	if validation.IsPasswordValid(p.Password) {
		validationErrors = append(validationErrors, *NewValidationError(ctx, InvalidPasswordError, InvalidPasswordDescription))
	}
	if !validation.IsEmailValid(p.Email) {
		validationErrors = append(validationErrors, *NewValidationError(ctx, InvalidPasswordError, InvalidPasswordDescription))
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return &validationErrors
}
