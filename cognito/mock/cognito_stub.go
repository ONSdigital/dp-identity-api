package mock

import (
	"errors"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type CognitoIdentityProviderClientStub struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	UserPools []string
	Users     []*User
	Sessions  []Session
	Groups    []*Group
}

func (m *CognitoIdentityProviderClientStub) DescribeUserPool(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	for _, v := range m.UserPools {
		if v == *poolInputData.UserPoolId {
			return nil, nil
		}
	}
	return nil, errors.New("failed to load user pool data")
}

func (m *CognitoIdentityProviderClientStub) AdminCreateUser(input *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	var (
		status, subjectAttrName, forenameAttrName, surnameAttrName, emailAttrName, username, subUUID, forename, surname, email string = "FORCE_CHANGE_PASSWORD", "sub", "given_name", "family_name", "email", "123e4567-e89b-12d3-a456-426614174000", "f0cf8dd9-755c-4caf-884d-b0c56e7d0704", "smileons", "bobbings", "emailx@ons.gov.uk"
	)

	if *input.UserAttributes[0].Value == "smileons" { // 201 - created successfully
		user := &models.CreateUserOutput{
			UserOutput: &cognitoidentityprovider.AdminCreateUserOutput{
				User: &cognitoidentityprovider.UserType{
					Attributes: []*cognitoidentityprovider.AttributeType{
						{
							Name:  &subjectAttrName,
							Value: &subUUID,
						},
						{
							Name:  &forenameAttrName,
							Value: &forename,
						},
						{
							Name:  &surnameAttrName,
							Value: &surname,
						},
						{
							Name:  &emailAttrName,
							Value: &email,
						},
					},
					Username:   &username,
					UserStatus: &status,
				},
			},
		}
		return user.UserOutput, nil
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Failed to create new user in user pool", nil) // 500 - internal exception error
}

func (m *CognitoIdentityProviderClientStub) InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	var expiration int64 = 123

	if *input.AuthFlow == "USER_PASSWORD_AUTH" {
		// non-verified response - ChallengName = "NEW_PASSWORD_REQUIRED"
		var (
			challengeName, sessionID string = "NEW_PASSWORD_REQUIRED", "AYABeBBsY5be-this-is-a-test-session-id-string-123456789iuerhcfdisieo-end"
		)
		initiateAuthOutputChallenge := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: nil,
			ChallengeName:        &challengeName,
			Session:              &sessionID,
		}

		// verified response - ChallengName = ""
		accessToken := "accessToken"
		idToken := "idToken"
		refreshToken := "refreshToken"
		initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    &expiration,
				IdToken:      &idToken,
				RefreshToken: &refreshToken,
			},
		}

		for _, user := range m.Users {
			if (user.Email == *input.AuthParameters["USERNAME"]) && (user.Password == *input.AuthParameters["PASSWORD"]) {
				// non-challenge response
				if user.Status == "CONFIRMED" {
					return initiateAuthOutput, nil
				} else {
					return initiateAuthOutputChallenge, nil
				}
			} else if user.Email != *input.AuthParameters["USERNAME"] {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Incorrect username or password.", nil)
			} else {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Password attempts exceeded", nil)
			}
		}

		if *input.AuthParameters["PASSWORD"] == "internalerrorException" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
		}
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInvalidParameterException, "A parameter was invalid", nil)
	} else if *input.AuthFlow == "REFRESH_TOKEN_AUTH" {
		if *input.AuthParameters["REFRESH_TOKEN"] == "InternalError" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
		} else if *input.AuthParameters["REFRESH_TOKEN"] == "ExpiredToken" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Refresh Token has expired", nil)
		} else {
			accessToken := "llll.mmmm.nnnn"
			idToken := "zzzz.yyyy.xxxx"
			initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{
				AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
					AccessToken: &accessToken,
					ExpiresIn:   &expiration,
					IdToken:     &idToken,
				},
			}
			return initiateAuthOutput, nil
		}
	} else {
		return nil, errors.New("InvalidParameterException: Unknown Auth Flow")
	}
}

func (m *CognitoIdentityProviderClientStub) GlobalSignOut(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
	if *signOutInput.AccessToken == "InternalError" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}
	for _, session := range m.Sessions {
		if session.AccessToken == *signOutInput.AccessToken {
			return &cognitoidentityprovider.GlobalSignOutOutput{}, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Access Token has been revoked", nil)
}

func (m *CognitoIdentityProviderClientStub) AdminUserGlobalSignOut(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
	if *adminUserGlobalSignOutInput.Username == "internalservererror@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	} else if *adminUserGlobalSignOutInput.Username == "clienterror@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Something went wrong", nil)
	}
	return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
}

func (m *CognitoIdentityProviderClientStub) ListUsers(input *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
	var (
		emailVerifiedAttr, emailVerifiedValue    string = "email_verified", "true"
		givenNameAttr, familyNameAttr, emailAttr string = "given_name", "family_name", "email"
		enabled                                  bool   = true
	)

	var usersList []*cognitoidentityprovider.UserType

	if len(m.Users) > 0 && m.Users[0].Email == "internal.error@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}

	if input.Filter != nil {
		getEmailFromFilter, _ := regexp.Compile(`^email\s\=\s(\D+.*)$`)
		email := getEmailFromFilter.ReplaceAllString(*input.Filter, `$1`)

		var emailRegex = regexp.MustCompile(`^\"email(\d)?@(ext\.)?ons.gov.uk\"`)
		if emailRegex.MatchString(email) {
			usersList = append(usersList, &cognitoidentityprovider.UserType{
				Attributes: []*cognitoidentityprovider.AttributeType{
					{
						Name:  &emailVerifiedAttr,
						Value: &emailVerifiedValue,
					},
				},
				Username: &email,
			})
		}
	} else {
		for _, user := range m.Users {
			userDetails := cognitoidentityprovider.UserType{
				Attributes: []*cognitoidentityprovider.AttributeType{
					{
						Name:  &emailVerifiedAttr,
						Value: &emailVerifiedValue,
					},
					{
						Name:  &givenNameAttr,
						Value: aws.String(user.GivenName),
					},
					{
						Name:  &familyNameAttr,
						Value: aws.String(user.FamilyName),
					},
					{
						Name:  &emailAttr,
						Value: aws.String(user.Email),
					},
				},
				Enabled:    &enabled,
				UserStatus: aws.String(user.Status),
				Username:   aws.String(user.ID),
			}
			usersList = append(usersList, &userDetails)
		}
	}
	users := &models.ListUsersOutput{
		ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
			Users: usersList,
		},
	}
	return users.ListUsersOutput, nil
}

func (m *CognitoIdentityProviderClientStub) RespondToAuthChallenge(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
	var expiration int64 = 123

	if *input.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		accessToken := "accessToken"
		idToken := "idToken"
		refreshToken := "refreshToken"
		challengeResponseOutput := &cognitoidentityprovider.RespondToAuthChallengeOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    &expiration,
				IdToken:      &idToken,
				RefreshToken: &refreshToken,
			},
		}

		if *input.ChallengeResponses["NEW_PASSWORD"] == "internalerrorException" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
		} else if *input.ChallengeResponses["NEW_PASSWORD"] == "invalidpassword" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeInvalidPasswordException, "password does not meet requirements", nil)
		}

		for _, user := range m.Users {
			if user.Email == *input.ChallengeResponses["USERNAME"] {
				return challengeResponseOutput, nil
			}
		}
		return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "user not found", nil)
	} else {
		return nil, errors.New("InvalidParameterException: Unknown Auth Flow")
	}
}

func (m *CognitoIdentityProviderClientStub) ConfirmForgotPassword(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {

	challengeResponseOutput := &cognitoidentityprovider.ConfirmForgotPasswordOutput{}

	if *input.Password == "internalerrorException" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	} else if *input.Password == "invalidpassword" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInvalidPasswordException, "password does not meet requirements", nil)
	} else if *input.ConfirmationCode == "invalidtoken" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeCodeMismatchException, "verification token does not meet requirements", nil)
	}

	for _, user := range m.Users {
		if user.Email == *input.Username {
			return challengeResponseOutput, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "user not found", nil)

}

func (m *CognitoIdentityProviderClientStub) ForgotPassword(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	forgotPasswordOutput := &cognitoidentityprovider.ForgotPasswordOutput{
		CodeDeliveryDetails: &cognitoidentityprovider.CodeDeliveryDetailsType{},
	}

	if *input.Username == "internal.error@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}
	if *input.Username == "too.many@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeTooManyRequestsException, "Slow down", nil)
	}

	for _, user := range m.Users {
		if user.Email == *input.Username {
			return forgotPasswordOutput, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "user not found", nil)
}

func (m *CognitoIdentityProviderClientStub) AdminGetUser(input *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	var (
		emailVerifiedAttr, emailVerifiedValue    string = "email_verified", "true"
		givenNameAttr, familyNameAttr, emailAttr string = "given_name", "family_name", "email"
	)
	for _, user := range m.Users {
		if user.ID == *input.Username {
			if user.Email == "internal.error@ons.gov.uk" {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
			}
			return &cognitoidentityprovider.AdminGetUserOutput{
				UserAttributes: []*cognitoidentityprovider.AttributeType{
					{
						Name:  &emailVerifiedAttr,
						Value: &emailVerifiedValue,
					},
					{
						Name:  &givenNameAttr,
						Value: aws.String(user.GivenName),
					},
					{
						Name:  &familyNameAttr,
						Value: aws.String(user.FamilyName),
					},
					{
						Name:  &emailAttr,
						Value: aws.String(user.Email),
					},
				},
				Enabled:    aws.Bool(user.Active),
				UserStatus: aws.String(user.Status),
				Username:   aws.String(user.ID),
			}, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "the user could not be found", nil)
}

func (m *CognitoIdentityProviderClientStub) CreateGroup(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error) {
	userPoolId := "aaaa-bbbb-ccc-dddd"

	if *input.GroupName == "internalError" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "something went wrong", nil)
	}

	for _, group := range m.Groups {
		if group.Name == *input.GroupName {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeGroupExistsException, "this group already exists", nil)
		}
	}

	newGroup, err := m.GenerateGroup(*input.GroupName, *input.Description, *input.Precedence)
	if err != nil {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, err.Error(), nil)
	}
	m.Groups = append(m.Groups, newGroup)

	return &cognitoidentityprovider.CreateGroupOutput{
		Group: &cognitoidentityprovider.GroupType{
			Description:  input.Description,
			GroupName:    input.GroupName,
			Precedence:   input.Precedence,
			CreationDate: &newGroup.Created,
			UserPoolId:   &userPoolId,
		},
	}, nil
}

func (m *CognitoIdentityProviderClientStub) AdminUpdateUserAttributes(input *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
	for _, user := range m.Users {
		if user.ID == *input.Username {
			if user.Email == "update.internalerror@ons.gov.uk" {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
			}
			for _, attr := range input.UserAttributes {
				if *attr.Name == "given_name" {
					user.GivenName = *attr.Value
				} else if *attr.Name == "family_name" {
					user.FamilyName = *attr.Value
				} else if *attr.Name == "custom:status_notes" {
					user.StatusNotes = *attr.Value
				}
			}
			return &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "the user could not be found", nil)
}

func (m *CognitoIdentityProviderClientStub) AdminEnableUser(input *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
	for _, user := range m.Users {
		if user.ID == *input.Username {
			if user.Email == "enable.internalerror@ons.gov.uk" {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong whilst enabling", nil)
			}
			user.Active = true
			return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "the user could not be found", nil)
}

func (m *CognitoIdentityProviderClientStub) AdminDisableUser(input *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
	for _, user := range m.Users {
		if user.ID == *input.Username {
			if user.Email == "disable.internalerror@ons.gov.uk" {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong whilst disabling", nil)
			}
			user.Active = false
			return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "the user could not be found", nil)
}

func (m *CognitoIdentityProviderClientStub) AdminAddUserToGroup(input *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
	if *input.GroupName == "internal-error" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}

	group := m.ReadGroup(*input.GroupName)
	if group == nil {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeResourceNotFoundException, "the group could not be found", nil)
	}

	user := m.ReadUser(*input.Username)
	if user == nil {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "the user could not be found", nil)
	}

	user.Groups = append(user.Groups, group)
	group.Members = append(group.Members, user)

	return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
}

func (m *CognitoIdentityProviderClientStub) GetGroup(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
	if *input.GroupName == "internal-error" || *input.GroupName == "get-group-internal-error" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}
	if *input.GroupName == "get-group-not-found" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeResourceNotFoundException, "get group - group not found", nil)
	}

	group := m.ReadGroup(*input.GroupName)
	if group == nil {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeResourceNotFoundException, "the group could not be found", nil)
	}
	timestamp := time.Now()
	return &cognitoidentityprovider.GetGroupOutput{
		Group: &cognitoidentityprovider.GroupType{
			CreationDate:     &group.Created,
			Description:      &group.Description,
			GroupName:        &group.Name,
			LastModifiedDate: &timestamp,
			Precedence:       &group.Precedence,
			UserPoolId:       input.UserPoolId,
		},
	}, nil
}

func (m *CognitoIdentityProviderClientStub) ListUsersInGroup(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
	var (
		emailVerifiedAttr, emailVerifiedValue    string = "email_verified", "true"
		givenNameAttr, familyNameAttr, emailAttr string = "given_name", "family_name", "email"
	)

	if *input.GroupName == "internal-error" || *input.GroupName == "list-group-users-internal-error" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}
	if *input.GroupName == "list-group-users-not-found" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeResourceNotFoundException, "list members - group not found", nil)
	}

	group := m.ReadGroup(*input.GroupName)
	if group == nil {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeResourceNotFoundException, "the group could not be found", nil)
	}
	var userList []*cognitoidentityprovider.UserType

	for _, user := range group.Members {
		userDetails := cognitoidentityprovider.UserType{
			Attributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  &emailVerifiedAttr,
					Value: &emailVerifiedValue,
				},
				{
					Name:  &givenNameAttr,
					Value: aws.String(user.GivenName),
				},
				{
					Name:  &familyNameAttr,
					Value: aws.String(user.FamilyName),
				},
				{
					Name:  &emailAttr,
					Value: aws.String(user.Email),
				},
			},
			Enabled:    &user.Active,
			UserStatus: aws.String(user.Status),
			Username:   aws.String(user.ID),
		}
		userList = append(userList, &userDetails)
	}
	return &cognitoidentityprovider.ListUsersInGroupOutput{
		Users: userList,
	}, nil
}

func (m *CognitoIdentityProviderClientStub) AdminRemoveUserFromGroup(input *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
	if *input.GroupName == "internal-error" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}

	group := m.ReadGroup(*input.GroupName)
	if group == nil {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeResourceNotFoundException, "the group could not be found", nil)
	}

	user := m.ReadUser(*input.Username)
	if user == nil {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "the user could not be found", nil)
	}

	var newGroupMembersList []*User
	for _, member := range group.Members {
		if member.ID != user.ID {
			newGroupMembersList = append(newGroupMembersList, member)
		}
	}
	group.Members = newGroupMembersList

	var newUserGroupList []*Group
	for _, memberGroup := range user.Groups {
		if memberGroup.Name != group.Name {
			newUserGroupList = append(newUserGroupList, group)
		}
	}
	user.Groups = newUserGroupList

	return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
}

func (m *CognitoIdentityProviderClientStub) AdminListGroupsForUser(
	input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
	if *input.Username == "internal-error" || *input.Username == "get-group-internal-error" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}
	if *input.Username == "get-user-not-found" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeResourceNotFoundException, "get user - user not found", nil)
	}

	Description := "some Group Desciption"
	Name := "test-group"
	var Precedence int64 = 97
	nextToken := ""
	usergroupList := cognitoidentityprovider.AdminListGroupsForUserOutput{}
	group := cognitoidentityprovider.GroupType{
		// CreationDate:     &timestamp,
		Description: &Description,
		GroupName:   &Name,
		// LastModifiedDate: &timestamp,
		Precedence: &Precedence,
		UserPoolId: input.UserPoolId,
	}
	usergroupList.Groups = append(usergroupList.Groups, &group)

	return &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups:    usergroupList.Groups,
		NextToken: &nextToken,
	}, nil
}
