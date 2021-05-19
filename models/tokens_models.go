package models

import "github.com/dgrijalva/jwt-go"

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

func (t *IdToken) SignToken(signingKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.Claims)
	return token.SignedString(signingKey)
}

func (t *IdToken) ParseWithoutValidating(tokenString string) error {
	parser := new(jwt.Parser)
	idToken, _, parsingErr := parser.ParseUnverified(tokenString, &IdClaims{})
	if parsingErr != nil {
		return parsingErr
	}
	idClaims := idToken.Claims.(*IdClaims)
	t.Claims = *idClaims
	return nil
}
