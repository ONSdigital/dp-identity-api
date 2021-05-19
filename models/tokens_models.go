package models

import (
	"errors"
	"fmt"
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
	Claims IdClaims
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
