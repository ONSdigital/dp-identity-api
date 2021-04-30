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
	ctx.Step(`^user pool with id "([^"]*)" exists$`, c.userPoolWithIdExists)
	ctx.Step(`^a user with username "([^"]*)" exists$`, c.userWithUsernameExists)
}

func (c *IdentityComponent) iShouldReceiveAHelloworldResponse() error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)

	assert.Equal(&c.ErrorFeature, `{"message":"Hello, World!"}`, strings.TrimSpace(string(body)))

	return c.ErrorFeature.StepError()
}

func (c *IdentityComponent) userPoolWithIdExists(userPoolId string) error {
	c.CognitoClient.AddUserPool(userPoolId)
	return nil
}

func (c *IdentityComponent) userWithUsernameExists(username string) error {
	c.CognitoClient.AddUserWithUsername(username)
	return nil
}
