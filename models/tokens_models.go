package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/apierrors"
	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
)

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

func (t *IdToken) ParseWithoutValidating(tokenString string) error {
	parser := new(jwt.Parser)

	idToken, _, parsingErr := parser.ParseUnverified(tokenString, &IdClaims{})
	if parsingErr != nil {
		fmt.Println(parsingErr)
		return parsingErr
	}

	idClaims, ok := idToken.Claims.(*IdClaims)
	if !ok {
		return errors.New("the ID token could not be parsed")
	}

	t.Claims = *idClaims
	return nil
}

func (t *IdToken) Validate(ctx context.Context, errorList *[]apierrors.IndividualError) {
	field := ""
	param := ""

	if t.TokenString == "" {
		invalidIDTokenErrorBody := apierrors.IndividualErrorBuilder(apierrors.InvalidIDTokenError,
			apierrors.MissingIDTokenMessage, field, param)
		*errorList = append(*errorList, invalidIDTokenErrorBody)
		log.Event(ctx, apierrors.MissingRefreshTokenMessage, log.ERROR)
	} else {
		parsingErr := t.ParseWithoutValidating(t.TokenString)
		if parsingErr != nil {
			invalidIDTokenErrorBody := apierrors.IndividualErrorBuilder(apierrors.InvalidIDTokenError,
				apierrors.MalformedIDTokenMessage, field, param)
			*errorList = append(*errorList, invalidIDTokenErrorBody)
			log.Event(ctx, parsingErr.Error(), log.ERROR)
		}
	}
}

type RefreshToken struct {
	TokenString string
}

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

func (t *RefreshToken) Validate(ctx context.Context, errorList *[]apierrors.IndividualError) {
	field := ""
	param := ""

	if t.TokenString == "" {
		invalidRefreshTokenErrorBody := apierrors.IndividualErrorBuilder(apierrors.InvalidRefreshTokenError,
			apierrors.MissingRefreshTokenMessage, field, param)
		*errorList = append(*errorList, invalidRefreshTokenErrorBody)
		log.Event(ctx, apierrors.MissingRefreshTokenMessage, log.ERROR)
	}
}
