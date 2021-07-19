package steps

import (
	"errors"

	"github.com/ONSdigital/dp-identity-api/api"

	"github.com/cucumber/godog"
)

func (c *IdentityComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^a user with email "([^"]*)" and password "([^"]*)" exists in the database$`, c.aUserWithEmailAndPasswordExistsInTheDatabase)
	ctx.Step(`^a user with username "([^"]*)" and email "([^"]*)" exists in the database$`, c.aUserWithUsernameAndEmailExistsInTheDatabase)
	ctx.Step(`^an internal server error is returned from Cognito$`, c.anInternalServerErrorIsReturnedFromCognito)
	ctx.Step(`^an error is returned from Cognito$`, c.anErrorIsReturnedFromCognito)
	ctx.Step(`^I have an active session with access token "([^"]*)"$`, c.iHaveAnActiveSessionWithAccessToken)
	ctx.Step(`I have a valid ID header for user "([^"]*)"$`, c.iHaveAValidIDHeaderForUser)
	ctx.Step(`^the AdminUserGlobalSignOut endpoint in cognito returns an internal server error$`, c.theAdminUserGlobalSignOutEndpointInCognitoReturnsAnInternalServerError)
	ctx.Step(`^a user with non-verified email "([^"]*)" and password "([^"]*)"$`, c.aUserWithNonverifiedEmailAndPassword)
	ctx.Step(`^user "([^"]*)" active is "([^"]*)"$`, c.userSetState)
}

func (c *IdentityComponent) aUserWithEmailAndPasswordExistsInTheDatabase(email, password string) error {
	c.CognitoClient.AddUserWithEmail(email, password, true)
	return nil
}

func (c *IdentityComponent) aUserWithUsernameAndEmailExistsInTheDatabase(username, email string) error {
	c.CognitoClient.AddUserWithUsername(username, email, true)
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

func (c *IdentityComponent) theAdminUserGlobalSignOutEndpointInCognitoReturnsAnInternalServerError() error {
	return nil
}

func (c *IdentityComponent) aUserWithNonverifiedEmailAndPassword(email, password string) error {
	c.CognitoClient.AddUserWithEmail(email, password, false)
	return nil
}

func (c *IdentityComponent) userSetState(username, active string) error {
	c.CognitoClient.SetUserActiveState(username, active)
	return nil
}
