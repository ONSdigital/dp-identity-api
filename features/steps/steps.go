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
	ctx.Step(`^I have an active session with access token "([^"]*)"$`, c.iHaveAnActiveSessionWithAccessToken)
}

func (c *IdentityComponent) iShouldReceiveAHelloworldResponse() error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)

	assert.Equal(&c.ErrorFeature, `{"message":"Hello, World!"}`, strings.TrimSpace(string(body)))

	return c.ErrorFeature.StepError()
}

func (c *IdentityComponent) iHaveAnActiveSessionWithAccessToken(accessToken string) error {
	c.CognitoClient.CreateSessionWithAccessToken(accessToken)
	return nil
}
