package steps

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-authorisation/v2/authorisationtest"

	"github.com/ONSdigital/dp-identity-api/api"
	"github.com/ONSdigital/dp-identity-api/models"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func (c *IdentityComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`I have a valid ID header for user "([^"]*)"$`, c.iHaveAValidIDHeaderForUser)
	ctx.Step(`^(\d+) groups exist in the database that username "([^"]*)" is a member$`, c.groupsExistInTheDatabaseThatUsernameIsAMember)
	ctx.Step(`^I GET the JSON web key set for cognito user pool$`, c.aResponseToAJWKSSetRequest)
	ctx.Step(`^I am a publisher user$`, c.publisherJWTToken)
	ctx.Step(`^I am an admin user$`, c.adminJWTToken)
	ctx.Step(`^I have an active session with access token "([^"]*)"$`, c.iHaveAnActiveSessionWithAccessToken)
	ctx.Step(`^a user with email "([^"]*)" and password "([^"]*)" exists in the database$`, c.aUserWithEmailAndPasswordExistsInTheDatabase)
	ctx.Step(`^a user with non-verified email "([^"]*)" and password "([^"]*)"$`, c.aUserWithNonverifiedEmailAndPassword)
	ctx.Step(`^a user with username "([^"]*)" and email "([^"]*)" exists in the database$`, c.aUserWithUsernameAndEmailExistsInTheDatabase)
	ctx.Step(`^a user with username "([^"]*)" and email "([^"]*)" and forename "([^"]*)" exists in the database$`, c.aUserWithUsernameAndEmailAndForenameExistsInTheDatabase)
	ctx.Step(`^an error is returned from Cognito$`, c.anErrorIsReturnedFromCognito)
	ctx.Step(`^an internal server error is returned from Cognito$`, c.anInternalServerErrorIsReturnedFromCognito)
	ctx.Step(`^group "([^"]*)" and description "([^"]*)" exists in the database$`, c.groupAndDescriptionExistsInTheDatabase)
	ctx.Step(`^group "([^"]*)" exists in the database$`, c.groupExistsInTheDatabase)
	ctx.Step(`^the AdminUserGlobalSignOut endpoint in cognito returns an internal server error$`, c.theAdminUserGlobalSignOutEndpointInCognitoReturnsAnInternalServerError)
	ctx.Step(`^the list response should contain "([^"]*)" entries$`, c.listResponseShouldContainCorrectNumberOfEntries)
	ctx.Step(`^the response code should be (\d+)$`, c.theResponseCodeShouldBe)
	ctx.Step(`^the response should match the following json for listgroups$`, c.theResponseShouldMatchTheFollowingJsonForListgroups)
	ctx.Step(`^there are "([^"]*)" active users and "([^"]*)" inactive users in the database$`, c.thereAreRequiredNumberOfActiveUsers)
	ctx.Step(`^there are "([^"]*)" users in the database$`, c.thereAreRequiredNumberOfUsers)
	ctx.Step(`^there are (\d+) groups in the database$`, c.thereAreGroupsInTheDatabase)
	ctx.Step(`^there are (\d+) users in group "([^"]*)"$`, c.thereAreUsersInGroup)
	ctx.Step(`^user "([^"]*)" active is "([^"]*)"$`, c.userSetState)
	ctx.Step(`^user "([^"]*)" is a member of group "([^"]*)"$`, c.userIsAMemberOfGroup)
	ctx.Step(`^request header Accept is "([^"]*)"$`, c.requestHeaderAcceptIs)
	ctx.Step(`^the response should match the following csv:$`, c.theResponseShouldMatchTheFollowingCsv)
	ctx.Step(`^the response header "([^"]*)" should contain "([^"]*)"$`, c.theResponseHeaderShouldContain)
	ctx.Step(`^a user with forename "([^"]*)", lastname "([^"]*)", email "([^"]*)", id "([^"]*)" and password "([^"]*)" exists in the database$`, c.aUserWithAttributesExistsInTheDatabase)
}

func (c *IdentityComponent) aResponseToAJWKSSetRequest() error {
	_, err := c.JWKSHandler.JWKSGetKeysetFunc("eu-west-1234XYZ", "eu-west-1234")
	return err
}

func (c *IdentityComponent) aUserWithEmailAndPasswordExistsInTheDatabase(email, password string) error {
	c.CognitoClient.AddUserWithEmail(email, password, true)
	return nil
}

func (c *IdentityComponent) aUserWithAttributesExistsInTheDatabase(forename, lastname, email, id, password string) error {
	c.CognitoClient.AddUserWithAttributes(id, forename, lastname, email, password, true)
	return nil
}

func (c *IdentityComponent) aUserWithUsernameAndEmailExistsInTheDatabase(username, email string) error {
	c.CognitoClient.AddUserWithUsername(username, email, true)
	return nil
}

func (c *IdentityComponent) aUserWithUsernameAndEmailAndForenameExistsInTheDatabase(username, email, forename string) error {
	c.CognitoClient.AddUserWithForename(username, email, forename, true)
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

func (c *IdentityComponent) groupAndDescriptionExistsInTheDatabase(groupName, description string) error {
	err := c.CognitoClient.AddGroupWithNameAndDescription(groupName, description)
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

func (c *IdentityComponent) groupsExistInTheDatabaseThatUsernameIsAMember(groupCount, userName string) error {
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
	responseBody := c.apiFeature.HTTPResponse.Body
	body, _ := io.ReadAll(responseBody)
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

func (c *IdentityComponent) thereAreGroupsInTheDatabase(groupCount string) error {

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
	actualStatusString := strconv.Itoa(c.apiFeature.HTTPResponse.StatusCode)
	if code != c.apiFeature.HTTPResponse.StatusCode {
		return errors.New("expected response status code to be: " + expectedStatusString + ", but actual is: " + actualStatusString)
	}
	return nil
}

func (c *IdentityComponent) theResponseShouldMatchTheFollowingJsonForListgroups(body *godog.DocString) (err error) {
	var expected, actual models.ListUserGroups
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	responseBody := c.apiFeature.HTTPResponse.Body
	resBody, _ := io.ReadAll(responseBody)
	if err = json.Unmarshal(resBody, &actual); err != nil {
		return
	}
	assert.Equal(c.apiFeature, expected.NextToken, actual.NextToken)
	assert.Equal(c.apiFeature, expected.Count, actual.Count)
	assert.Equal(c.apiFeature, len(expected.Groups), actual.Count)

	if actual.Count > 0 {
		assert.Equal(c.apiFeature, *expected.Groups[0].Name, *actual.Groups[0].Name)
		assert.Equal(c.apiFeature, *expected.Groups[0].ID, *actual.Groups[0].ID)
		tmpPrecedence := int(*expected.Groups[0].Precedence)
		assert.GreaterOrEqual(c.apiFeature, 13, tmpPrecedence)
		assert.LessOrEqual(c.apiFeature, 100, tmpPrecedence)
		tmpPrecedence = int(*actual.Groups[0].Precedence)
		assert.GreaterOrEqual(c.apiFeature, 13, tmpPrecedence)
		assert.LessOrEqual(c.apiFeature, 100, tmpPrecedence)
	}
	return nil
}

func (c *IdentityComponent) requestHeaderAcceptIs() error {
	err := c.apiFeature.ISetTheHeaderTo("Accept", "text/csv")
	return err
}

func (c *IdentityComponent) theResponseShouldMatchTheFollowingCsv(body *godog.DocString) (err error) {
	tmpExpected, _ := io.ReadAll(c.apiFeature.HTTPResponse.Body)
	actual := strings.Replace(strings.TrimSpace(string(tmpExpected[:])), "\t", "", -1)
	expected := strings.Replace(strings.TrimSpace(body.Content), "\t", "", -1)

	if actual != expected {
		return errors.New("expected body to be: " + "\n" + expected + "\n\t but actual is: " + "\n" + actual)
	}
	return nil
}

func (c *IdentityComponent) theResponseHeaderShouldContain(key, value string) (err error) {
	responseHeader := c.apiFeature.HTTPResponse.Header
	actualValue, actualExist := responseHeader[key]
	if !actualExist {
		return errors.New("expected header key " + key + ", does not exist in the header ")
	}
	if actualValue[0] != value {
		return errors.New("expected header value " + value + ", but is actually is :" + actualValue[0])
	}

	return nil
}
