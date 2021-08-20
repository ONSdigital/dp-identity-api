package steps

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"

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
	ctx.Step(`^group "([^"]*)" exists in the database$`, c.groupExistsInTheDatabase)
	ctx.Step(`^there are "([^"]*)" users in group "([^"]*)"$`, c.thereAreUsersInGroup)
	ctx.Step(`^user "([^"]*)" is a member of group "([^"]*)"$`, c.userIsAMemberOfGroup)
	ctx.Step(`^there are "([^"]*)" users in the database$`, c.thereAreRequiredNumberOfUsers)
	ctx.Step(`^the list response should contain "([^"]*)" entries$`, c.listResponseShouldContainCorrectNumberOfEntries)

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

func (c *IdentityComponent) groupExistsInTheDatabase(groupName string) error {
	err := c.CognitoClient.AddGroupWithName(groupName)
	return err
}

func (c *IdentityComponent) userIsAMemberOfGroup(username, groupName string) error {
	err := c.CognitoClient.AddUserToGroup(username, groupName)
	return err
}

func (c *IdentityComponent) thereAreUsersInGroup(userCount, groupName string) error {
	group := c.CognitoClient.ReadGroup(groupName)
	if group == nil {
		return errors.New("group not found")
	}
	userCountInt, err := strconv.Atoi(userCount)
	if err != nil {
		return errors.New("could not convert user count to int")
	}
	assert.Equal(c.apiFeature, userCountInt, len(group.Members))
	return nil
}

// listResponseShouldContainCorrectNumberOfEntries asserts that the list response 'count' matches the expected value
func (c *IdentityComponent) listResponseShouldContainCorrectNumberOfEntries(expectedListLength string) error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)
	var bodyObject map[string]interface{}
	err := json.Unmarshal(body, &bodyObject)
	if err != nil {
		return err
	}
	expectedListLengthInt, err := strconv.Atoi(expectedListLength)
	if err != nil {
		return err
	}
	assert.Equal(c.apiFeature, bodyObject["count"], expectedListLengthInt)
	return nil
}

// thereAreRequiredNumberOfUsers asserts that the list response 'count' matches the expected value
func (c *IdentityComponent) thereAreRequiredNumberOfUsers(requiredNumberOfUsers string) error {
	requiredNumberOfUsersInt, err := strconv.Atoi(requiredNumberOfUsers)
	if err != nil {
		return err
	}
	c.CognitoClient.AddMultipleUsers(requiredNumberOfUsersInt)
	return nil
}
