package mock

import (
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/golang-jwt/jwt/v4"
)

type Session struct {
	AccessToken  string
	IdToken      string
	RefreshToken string
}

func (m *CognitoIdentityProviderClientStub) CreateSessionWithAccessToken(accessToken string) {
	idToken, refreshToken := "aaaa.bbbb.cccc", "dddd.eeee.ffff"
	m.Sessions = append(m.Sessions, m.GenerateSession(accessToken, idToken, refreshToken))
}

func (m *CognitoIdentityProviderClientStub) GenerateSession(accessToken string, idToken string, refreshToken string) Session {
	return Session{
		AccessToken:  accessToken,
		IdToken:      idToken,
		RefreshToken: refreshToken,
	}
}

func (m *CognitoIdentityProviderClientStub) CreateIdTokenForEmail(email string) string {
	return GenerateMockIDToken(email)
}

func GenerateMockIDToken(email string) string {
	testSigningKey := []byte("TestSigningKey")
	idClaims := models.IdClaims{
		Sub:           "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		Aud:           "xxxxxxxxxxxxexample",
		EmailVerified: true,
		TokenUse:      "id",
		AuthTime:      1500009400,
		Iss:           "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_example",
		CognitoUser:   "TestONS",
		Exp:           1500013000,
		GivenName:     "Test",
		Iat:           1500009400,
		Email:         email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, idClaims)
	signedString, err := token.SignedString(testSigningKey)
	if err != nil {
		return ""
	}
	return signedString
}
