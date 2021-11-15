package steps

import (
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-authorisation/v2/authorisationtest"
	"io/ioutil"
	"strconv"

	"github.com/stretchr/testify/assert"

	"github.com/ONSdigital/dp-identity-api/api"
	"github.com/ONSdigital/dp-identity-api/models"

	"github.com/cucumber/godog"
)

func (c *IdentityComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^a user with email "([^"]*)" and password "([^"]*)" exists in the database$`, c.aUserWithEmailAndPasswordExistsInTheDatabase)
	ctx.Step(`^a user with username "([^"]*)" and email "([^"]*)" exists in the database$`, c.aUserWithUsernameAndEmailExistsInTheDatabase)
	ctx.Step(`^an internal server error is returned from Cognito$`, c.anInternalServerErrorIsReturnedFromCognito)
	ctx.Step(`^an error is returned from Cognito$`, c.anErrorIsReturnedFromCognito)
	ctx.Step(`^I am an admin user$`, c.adminJWTToken)
	ctx.Step(`^I am a publisher user$`, c.publisherJWTToken)
	ctx.Step(`^I have an active session with access token "([^"]*)"$`, c.iHaveAnActiveSessionWithAccessToken)
	ctx.Step(`I have a valid ID header for user "([^"]*)"$`, c.iHaveAValidIDHeaderForUser)
	ctx.Step(`^the AdminUserGlobalSignOut endpoint in cognito returns an internal server error$`, c.theAdminUserGlobalSignOutEndpointInCognitoReturnsAnInternalServerError)
	ctx.Step(`^a user with non-verified email "([^"]*)" and password "([^"]*)"$`, c.aUserWithNonverifiedEmailAndPassword)
	ctx.Step(`^user "([^"]*)" active is "([^"]*)"$`, c.userSetState)
	ctx.Step(`^group "([^"]*)" exists in the database$`, c.groupExistsInTheDatabase)
	ctx.Step(`^there are "([^"]*)" users in group "([^"]*)"$`, c.thereAreUsersInGroup)
	ctx.Step(`^user "([^"]*)" is a member of group "([^"]*)"$`, c.userIsAMemberOfGroup)
	ctx.Step(`^there are "([^"]*)" users in the database$`, c.thereAreRequiredNumberOfUsers)
	ctx.Step(`^there are "([^"]*)" active users and "([^"]*)" inactive users in the database$`, c.thereAreRequiredNumberOfActiveUsers)
	ctx.Step(`^the list response should contain "([^"]*)" entries$`, c.listResponseShouldContainCorrectNumberOfEntries)
	ctx.Step(`^there (\d+) groups exists in the database that username "([^"]*)" is a member$`, c.thereGroupsExistsInTheDatabaseThatUsernameIsAMember)
	ctx.Step(`^there "([^"]*)" groups exists in the database$`, c.thereGroupsExistsInTheDatabase)
	ctx.Step(`^the response code should be (\d+)$`, c.theResponseCodeShouldBe)
	ctx.Step(`^the response should match the following json for listgroups$`, c.theResponseShouldMatchTheFollowingJsonForListgroups)

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

func (c *IdentityComponent) adminJWTToken() error {
	err := c.apiFeature.ISetTheHeaderTo(api.AccessTokenHeaderName, authorisationtest.AdminJWTToken)
	return err
}

func (c *IdentityComponent) publisherJWTToken() error {
	err := c.apiFeature.ISetTheHeaderTo(api.AccessTokenHeaderName, authorisationtest.PublisherJWTToken)
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

func (c *IdentityComponent) thereGroupsExistsInTheDatabaseThatUsernameIsAMember(groupCount, userName string) error {
	user := c.CognitoClient.ReadUser(userName)
	if user == nil {
		return errors.New("user not found")
	}
	groupCountInt, err := strconv.Atoi(groupCount)
	if err != nil {
		return errors.New("could not convert " + groupCount + " to int")
	}
	c.CognitoClient.BulkGenerateGroups(groupCountInt)
	c.CognitoClient.MakeUserMember(user.ID)

	assert.Equal(c.apiFeature, groupCountInt, len(user.Groups))
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
func (c *IdentityComponent) thereAreRequiredNumberOfActiveUsers(requiredNumberOfActiveUsers, requiredNumberOfInActiveUsers string) error {
	requiredNumberOfActiveUsersInt, err := strconv.Atoi(requiredNumberOfActiveUsers)
	if err != nil {
		return err
	}
	requiredNumberOfInActiveUsersInt, err := strconv.Atoi(requiredNumberOfInActiveUsers)
	if err != nil {
		return err
	}
	c.CognitoClient.AddMultipleActiveUsers(requiredNumberOfActiveUsersInt, requiredNumberOfInActiveUsersInt)
	return nil
}

func (c *IdentityComponent) thereGroupsExistsInTheDatabase(groupCount string) error {

	groupCountInt, err := strconv.Atoi(groupCount)
	if err != nil {
		return errors.New("could not convert" + groupCount + "to int")
	}
	c.CognitoClient.BulkGenerateGroupsList(groupCountInt)
	tmpGroupsList := c.CognitoClient.GroupsList
	lengroups := 0
	for _, x := range tmpGroupsList {
		lengroups = lengroups + len(x.Groups)
	}

	assert.Equal(c.apiFeature, groupCountInt, lengroups)
	return nil
}

func (c *IdentityComponent) theResponseCodeShouldBe(code int) error {
	expectedStatusString := strconv.Itoa(code)
	actualStatusString := strconv.Itoa(c.apiFeature.HttpResponse.StatusCode)
	if code != c.apiFeature.HttpResponse.StatusCode {
		return errors.New("expected response status code to be: " + expectedStatusString + ", but actual is: " + actualStatusString)
	}
	return nil
}

func (c *IdentityComponent) theResponseShouldMatchTheFollowingJsonForListgroups(body *godog.DocString) (err error) {
	var expected, actual models.ListUserGroups

	// re-encode expected response
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	responseBody := c.apiFeature.HttpResponse.Body
	resBody, _ := ioutil.ReadAll(responseBody)
	if err = json.Unmarshal(resBody, &actual); err != nil {
		return
	}

	// the matching may be adapted per different requirements.

	assert.Equal(c.apiFeature, expected.NextToken, actual.NextToken)
	assert.Equal(c.apiFeature, expected.Count, actual.Count)
	assert.Equal(c.apiFeature, len(expected.Groups), actual.Count)
	// if actual.Count > 0 && expected.Count > 0 {
	if actual.Count > 0 {

		assert.Equal(c.apiFeature, *expected.Groups[0].Description, *actual.Groups[0].Description)
		assert.Equal(c.apiFeature, *expected.Groups[0].GroupName, *actual.Groups[0].GroupName)
		tmpPrecedence := int(*expected.Groups[0].Precedence)
		assert.GreaterOrEqual(c.apiFeature, 13, tmpPrecedence)
		assert.LessOrEqual(c.apiFeature, 100, tmpPrecedence)
		tmpPrecedence = int(*actual.Groups[0].Precedence)
		assert.GreaterOrEqual(c.apiFeature, 13, tmpPrecedence)
		assert.LessOrEqual(c.apiFeature, 100, tmpPrecedence)

	}
	return nil
}
