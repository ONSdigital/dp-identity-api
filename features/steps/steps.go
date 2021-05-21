package steps

import (
	"errors"
	"github.com/ONSdigital/dp-identity-api/api"
	"io/ioutil"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func (c *IdentityComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^I should receive a hello-world response$`, c.iShouldReceiveAHelloworldResponse)

	ctx.Step(`^a user with email "([^"]*)" and password "([^"]*)" exists in the database$`, c.aUserWithEmailAndPasswordExistsInTheDatabase)
	ctx.Step(`^an internal server error is returned from Cognito$`, c.anInternalServerErrorIsReturnedFromCognito)
	ctx.Step(`^an error is returned from Cognito$`, c.anErrorIsReturnedFromCognito)
	ctx.Step(`^I have an active session with access token "([^"]*)"$`, c.iHaveAnActiveSessionWithAccessToken)
	ctx.Step(`I have a valid ID header for user "([^"]*)"$`, c.iHaveAValidIDHeaderForUser)
}

func (c *IdentityComponent) iShouldReceiveAHelloworldResponse() error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)

	assert.Equal(&c.ErrorFeature, `{"message":"Hello, World!"}`, strings.TrimSpace(string(body)))

	return c.ErrorFeature.StepError()
}

func (c *IdentityComponent) aUserWithEmailAndPasswordExistsInTheDatabase(username, password string) error {
	c.CognitoClient.AddUserWithUsername(username, password)
	return nil
}

func (c *IdentityComponent) anInternalServerErrorIsReturnedFromCognito() error {
	return nil
}

func (c *IdentityComponent) anErrorIsReturnedFromCognito() error {
	return nil
}

func (c *IdentityComponent) iHaveAnActiveSessionWithAccessToken(accessToken string) error {
	c.CognitoClient.CreateSessionWithAccessToken(accessToken)
	return nil
}

func (c *IdentityComponent) iHaveAValidIDHeaderForUser(email string) error {
	idToken := c.CognitoClient.CreateIdTokenForEmail(email)
	if idToken == "" {
		return errors.New("id token generation failed")
	}
	err := c.apiFeature.ISetTheHeaderTo(api.IdTokenHeaderName, idToken)
	return err
}
