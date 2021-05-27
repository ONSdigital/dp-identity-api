package models

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type AccessToken struct {
	AuthHeader  string
	TokenString string
}

func (t *AccessToken) Validate(ctx context.Context) error {
	if t.AuthHeader == "" {
		return NewValidationError(ctx, InvalidTokenError, MissingAuthorizationTokenDescription)
	}
	authComponents := strings.Split(t.AuthHeader, " ")
	if len(authComponents) != 2 {
		return NewValidationError(ctx, InvalidTokenError, MalformedAuthorizationTokenDescription)
	}
	t.TokenString = authComponents[1]
	return nil
}

func (t *AccessToken) GenerateSignOutRequest() *cognitoidentityprovider.GlobalSignOutInput {
	return &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: &t.TokenString}
}

type IdClaims struct {
	Sub           string `json:"sub"`
	Aud           string `json:"aud"`
	EmailVerified bool   `json:"email_verified"`
	TokenUse      string `json:"token_use"`
	AuthTime      int    `json:"auth_time"`
	Iss           string `json:"iss"`
	CognitoUser   string `json:"cognito:username"`
	Exp           int    `json:"exp"`
	GivenName     string `json:"given_name"`
	Iat           int    `json:"iat"`
	Email         string `json:"email"`
	jwt.StandardClaims
}

type IdToken struct {
	TokenString string
	Claims      IdClaims
}

//ParseWithoutValidating parses the claims in an ID token JWT in to a IdClaims struct without validating the token
func (t *IdToken) ParseWithoutValidating(tokenString string) error {
	parser := new(jwt.Parser)

	idToken, _, parsingErr := parser.ParseUnverified(tokenString, &IdClaims{})
	if parsingErr != nil {
		return parsingErr
	}

	idClaims, ok := idToken.Claims.(*IdClaims)
	if !ok {
		return errors.New("the ID token could not be parsed")
	}

	t.Claims = *idClaims
	return nil
}

//Validate validates the existence of a JWT string and that it is correctly formatting, storing the tokens claims in an IdClaims struct
func (t *IdToken) Validate(ctx context.Context, errorList []apierrorsdeprecated.Error) []apierrorsdeprecated.Error {
	if t.TokenString == "" {
		log.Event(ctx, apierrorsdeprecated.MissingRefreshTokenMessage, log.ERROR)
		invalidIDTokenErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.InvalidIDTokenError,
			apierrorsdeprecated.MissingIDTokenMessage)
		return append(errorList, invalidIDTokenErrorBody)
	}
	parsingErr := t.ParseWithoutValidating(t.TokenString)
	if parsingErr != nil {
		log.Event(ctx, parsingErr.Error(), log.ERROR)
		invalidIDTokenErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.InvalidIDTokenError,
			apierrorsdeprecated.MalformedIDTokenMessage)
		errorList = append(errorList, invalidIDTokenErrorBody)
	}
	return errorList
}

type RefreshToken struct {
	TokenString string
}

//GenerateRefreshRequest produces a Cognito InitiateAuthInput struct for refreshing a users current session
func (t *RefreshToken) GenerateRefreshRequest(clientSecret string, username string, clientId string) *cognitoidentityprovider.InitiateAuthInput {
	refreshAuthFlow := "REFRESH_TOKEN_AUTH"

	secretHash := utilities.ComputeSecretHash(clientSecret, username, clientId)

	authParams := map[string]*string{
		"REFRESH_TOKEN": &t.TokenString,
		"SECRET_HASH":   &secretHash,
	}

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       &refreshAuthFlow,
		AuthParameters: authParams,
		ClientId:       &clientId,
	}
	return authInput
}

//Validate validates the existence of a JWT string
func (t *RefreshToken) Validate(ctx context.Context, errorList []apierrorsdeprecated.Error) []apierrorsdeprecated.Error {
	if t.TokenString != "" {
		return errorList
	}
	log.Event(ctx, apierrorsdeprecated.MissingRefreshTokenMessage, log.ERROR)
	invalidRefreshTokenErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.InvalidRefreshTokenError,
		apierrorsdeprecated.MissingRefreshTokenMessage)
	return append(errorList, invalidRefreshTokenErrorBody)
}
