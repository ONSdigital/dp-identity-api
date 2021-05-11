package steps

import (
	"io/ioutil"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func (c *IdentityComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^I should receive a hello-world response$`, c.iShouldReceiveAHelloworldResponse)
	ctx.Step(`^a user with email "([^"]*)" and password "([^"]*)" exists in the database$`, c.aUserWithEmailAndPasswordExistsInTheDatabase)
	ctx.Step(`^a user does not exists in the database$`, c.aUserDoesNotExistsInTheDatabase)
	ctx.Step(`^Cognito has an internal server error$`, c.cognitoHasAnInternalServerError)
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

func (c *IdentityComponent) aUserDoesNotExistsInTheDatabase() error {
	return nil
}

func (c *IdentityComponent) cognitoHasAnInternalServerError() error {
	return nil
}
