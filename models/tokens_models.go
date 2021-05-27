package models

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type AccessToken struct {
	AuthHeader  string
	TokenString string
}

func (t *AccessToken) Validate(ctx context.Context) *Error {
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
func (t *IdToken) ParseWithoutValidating(ctx context.Context, tokenString string) *Error {
	parser := new(jwt.Parser)

	idToken, _, parsingErr := parser.ParseUnverified(tokenString, &IdClaims{})
	if parsingErr != nil {
		return NewError(ctx, parsingErr, InvalidTokenError, MalformedIDTokenDescription)
	}

	idClaims, ok := idToken.Claims.(*IdClaims)
	if !ok {
		return NewValidationError(ctx, InvalidTokenError, MalformedIDTokenDescription)
	}

	t.Claims = *idClaims
	return nil
}

//Validate validates the existence of a JWT string and that it is correctly formatting, storing the tokens claims in an IdClaims struct
func (t *IdToken) Validate(ctx context.Context) *Error {
	if t.TokenString == "" {
		return NewValidationError(ctx, InvalidTokenError, MissingIDTokenDescription)
	}
	parsingErr := t.ParseWithoutValidating(ctx, t.TokenString)
	if parsingErr != nil {
		return parsingErr
	}
	return nil
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
func (t *RefreshToken) Validate(ctx context.Context) *Error {
	if t.TokenString != "" {
		return nil
	}
	return NewValidationError(ctx, InvalidTokenError, MissingRefreshTokenDescription)
}

func (t *RefreshToken) BuildSuccessfulJsonResponse(ctx context.Context, result *cognitoidentityprovider.InitiateAuthOutput) ([]byte, error) {
	if result.AuthenticationResult != nil {
		tokenDuration := time.Duration(*result.AuthenticationResult.ExpiresIn)
		expirationTime := time.Now().UTC().Add(time.Second * tokenDuration).String()

		postBody := map[string]interface{}{"expirationTime": expirationTime}

		jsonResponse, err := json.Marshal(postBody)
		if err != nil {
			responseErr := NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
			return nil, responseErr
		}
		return jsonResponse, nil
	} else {
		responseErr := NewValidationError(ctx, InternalError, UnrecognisedCognitoResponseDescription)
		return nil, responseErr
	}
}
