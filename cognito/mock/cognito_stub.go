package mock

import (
	"context"
	"errors"

	"github.com/aws/smithy-go"

	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

const (
	internalError           = "internal-error"
	errCodeNotAuthorized    = "NotAuthorizedException"
	errCodeInternalError    = "InternalErrorException"
	errCodeInvalidParam     = "InvalidParameterException"
	errCodeInvalidPassword  = "InvalidPasswordException"
	errCodeUserNotFound     = "UserNotFoundException"
	errCodeCodeMismatch     = "CodeMismatchException"
	errCodeExpiredCode      = "ExpiredCodeException"
	errCodeTooManyRequests  = "TooManyRequestsException"
	errCodeGroupExists      = "GroupExistsException"
	errCodeResourceNotFound = "ResourceNotFoundException"
)

type CognitoIdentityProviderClientStub struct {
	cognitoidentityprovider.Client
	UserPools  []string
	Users      []*User
	Sessions   []Session
	Groups     []*Group
	GroupsList []cognitoidentityprovider.ListGroupsOutput
}

func (m *CognitoIdentityProviderClientStub) DescribeUserPool(_ context.Context, poolInputData *cognitoidentityprovider.DescribeUserPoolInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	for _, v := range m.UserPools {
		if v == *poolInputData.UserPoolId {
			return nil, nil
		}
	}
	return nil, errors.New("failed to load user pool data")
}

func (m *CognitoIdentityProviderClientStub) AdminCreateUser(_ context.Context, input *cognitoidentityprovider.AdminCreateUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	var (
		subjectAttrName, forenameAttrName, surnameAttrName, emailAttrName, username, subUUID, forename, surname, email = "sub", "given_name", "family_name", "email", "123e4567-e89b-12d3-a456-426614174000", "f0cf8dd9-755c-4caf-884d-b0c56e7d0704", "smileons", "bobbings", "emailx@ons.gov.uk"
		status                                                                                                         = types.UserStatusTypeForceChangePassword
	)

	if *input.UserAttributes[0].Value == "smileons" { // 201 - created successfully
		user := &models.CreateUserOutput{
			UserOutput: &cognitoidentityprovider.AdminCreateUserOutput{
				User: &types.UserType{
					Attributes: []types.AttributeType{
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
					UserStatus: status,
				},
			},
		}
		return user.UserOutput, nil
	}
	return nil, &smithy.GenericAPIError{
		Code:    errCodeInternalError,
		Message: "Failed to create new user in user pool",
	}
}

func (m *CognitoIdentityProviderClientStub) InitiateAuth(_ context.Context, input *cognitoidentityprovider.InitiateAuthInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	var expiration int32 = 123

	if input.AuthFlow == types.AuthFlowTypeUserPasswordAuth {
		// non-verified response - ChallengName = "NEW_PASSWORD_REQUIRED"
		var (
			challengeName = types.ChallengeNameTypeNewPasswordRequired
			sessionID     = "AYABeBBsY5be-this-is-a-test-session-id-string-123456789iuerhcfdisieo-end"
		)
		initiateAuthOutputChallenge := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: nil,
			ChallengeName:        challengeName,
			Session:              &sessionID,
		}

		// verified response - ChallengName = ""
		accessToken := "accessToken"
		idToken := "idToken"
		refreshToken := "refreshToken"
		initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &types.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    expiration,
				IdToken:      &idToken,
				RefreshToken: &refreshToken,
			},
		}

		for _, user := range m.Users {
			if (user.Email == input.AuthParameters["USERNAME"]) && (user.Password == input.AuthParameters["PASSWORD"]) {
				// non-challenge response
				if user.Status == "CONFIRMED" {
					return initiateAuthOutput, nil
				}
				return initiateAuthOutputChallenge, nil
			} else if user.Email != input.AuthParameters["USERNAME"] {
				return nil, &smithy.GenericAPIError{
					Code:    errCodeNotAuthorized,
					Message: "Incorrect username or password.",
				}
			} else if user.Password != input.AuthParameters["PASSWORD"] {
				return nil, &smithy.GenericAPIError{
					Code:    errCodeNotAuthorized,
					Message: "Password attempts exceeded",
				}
			}
		}

		if input.AuthParameters["PASSWORD"] == "internalerrorException" {
			return nil, &smithy.GenericAPIError{
				Code:    errCodeInternalError,
				Message: "Something went wrong",
			}
		}
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInvalidParam,
			Message: "A parameter was invalid",
		}
	} else if input.AuthFlow == "REFRESH_TOKEN_AUTH" {
		if input.AuthParameters["REFRESH_TOKEN"] == "InternalError" {
			return nil, &smithy.GenericAPIError{
				Code:    errCodeInternalError,
				Message: "Something went wrong",
			}
		} else if input.AuthParameters["REFRESH_TOKEN"] == "ExpiredToken" {
			return nil, &smithy.GenericAPIError{
				Code:    errCodeNotAuthorized,
				Message: "Refresh Token has expired",
			}
		}
		accessToken := "llll.mmmm.nnnn"
		idToken := "zzzz.yyyy.xxxx"
		initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &types.AuthenticationResultType{
				AccessToken: &accessToken,
				ExpiresIn:   expiration,
				IdToken:     &idToken,
			},
		}
		return initiateAuthOutput, nil
	}
	return nil, errors.New("InvalidParameterException: Unknown Auth Flow")
}

func (m *CognitoIdentityProviderClientStub) GlobalSignOut(_ context.Context, signOutInput *cognitoidentityprovider.GlobalSignOutInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
	if *signOutInput.AccessToken == "InternalError" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}
	for _, session := range m.Sessions {
		if session.AccessToken == *signOutInput.AccessToken {
			return &cognitoidentityprovider.GlobalSignOutOutput{}, nil
		}
	}
	return nil, &smithy.GenericAPIError{
		Code:    errCodeNotAuthorized,
		Message: "Access Token has been revoked",
	}
}

func (m *CognitoIdentityProviderClientStub) AdminUserGlobalSignOut(_ context.Context, adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
	if *adminUserGlobalSignOutInput.Username == "internalservererror@ons.gov.uk" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "something went wrong",
		}
	} else if *adminUserGlobalSignOutInput.Username == "clienterror@ons.gov.uk" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeNotAuthorized,
			Message: "something went wrong",
		}
	}
	return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
}

func (m *CognitoIdentityProviderClientStub) ListUsers(_ context.Context, input *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
	var (
		emailVerifiedAttr, emailVerifiedValue    = "email_verified", "true"
		givenNameAttr, familyNameAttr, emailAttr = "given_name", "family_name", "email"
	)

	var usersList []types.UserType

	if len(m.Users) > 0 && m.Users[0].Email == "internal.error@ons.gov.uk" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}

	if input.Filter != nil {
		getEmailFromFilter := regexp.MustCompile(`^email\s=\s(\D+.*)$`)
		email := getEmailFromFilter.ReplaceAllString(*input.Filter, `$1`)

		var emailRegex = regexp.MustCompile(`^\"email(\d)?@(ext\.)?ons.gov.uk\"`)
		if emailRegex.MatchString(email) {
			usersList = append(usersList, types.UserType{
				Attributes: []types.AttributeType{
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
			userDetails := types.UserType{
				Attributes: []types.AttributeType{
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
				Enabled:    true,
				UserStatus: user.Status,
				Username:   aws.String(user.ID),
			}
			usersList = append(usersList, userDetails)
		}
	}
	users := &models.ListUsersOutput{
		ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
			Users: usersList,
		},
	}
	return users.ListUsersOutput, nil
}

func (m *CognitoIdentityProviderClientStub) RespondToAuthChallenge(_ context.Context, input *cognitoidentityprovider.RespondToAuthChallengeInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
	var expiration int32 = 123

	if input.ChallengeName == types.ChallengeNameTypeNewPasswordRequired {
		accessToken := "accessToken"
		idToken := "idToken"
		refreshToken := "refreshToken"
		challengeResponseOutput := &cognitoidentityprovider.RespondToAuthChallengeOutput{
			AuthenticationResult: &types.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    expiration,
				IdToken:      &idToken,
				RefreshToken: &refreshToken,
			},
		}

		if input.ChallengeResponses["NEW_PASSWORD"] == "internalerrorException" {
			return nil, &smithy.GenericAPIError{
				Code:    errCodeInternalError,
				Message: "Something went wrong",
			}
		} else if input.ChallengeResponses["NEW_PASSWORD"] == "invalidpassword" {
			return nil, &smithy.GenericAPIError{
				Code:    errCodeInvalidPassword,
				Message: "password does not meet requirements",
			}
		}

		for _, user := range m.Users {
			if user.Email == input.ChallengeResponses["USERNAME"] {
				return challengeResponseOutput, nil
			}
		}
		return nil, &smithy.GenericAPIError{
			Code:    errCodeUserNotFound,
			Message: "user not found",
		}
	}
	return nil, errors.New("InvalidParameterException: Unknown Auth Flow")
}

func (m *CognitoIdentityProviderClientStub) ConfirmForgotPassword(_ context.Context, input *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
	challengeResponseOutput := &cognitoidentityprovider.ConfirmForgotPasswordOutput{}

	if *input.Password == "internalerrorException" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	} else if *input.Password == "invalidpassword" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInvalidPassword,
			Message: "password does not meet requirements",
		}
	} else if *input.ConfirmationCode == "invalid-token" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeCodeMismatch,
			Message: "verification token does not meet requirements",
		}
	} else if *input.ConfirmationCode == "expired-token" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeExpiredCode,
			Message: "verification token has expired",
		}
	}

	for _, user := range m.Users {
		if user.Email == *input.Username {
			return challengeResponseOutput, nil
		}
	}
	return nil, &smithy.GenericAPIError{
		Code:    errCodeUserNotFound,
		Message: "user not found",
	}
}

func (m *CognitoIdentityProviderClientStub) ForgotPassword(_ context.Context, input *cognitoidentityprovider.ForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	forgotPasswordOutput := &cognitoidentityprovider.ForgotPasswordOutput{
		CodeDeliveryDetails: &types.CodeDeliveryDetailsType{},
	}

	if *input.Username == "internal.error@ons.gov.uk" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}
	if *input.Username == "too.many@ons.gov.uk" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeTooManyRequests,
			Message: "Slow down",
		}
	}

	for _, user := range m.Users {
		if user.Email == *input.Username {
			return forgotPasswordOutput, nil
		}
	}
	return nil, &smithy.GenericAPIError{
		Code:    errCodeUserNotFound,
		Message: "user not found",
	}
}

func (m *CognitoIdentityProviderClientStub) AdminGetUser(_ context.Context, input *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	var (
		emailVerifiedAttr, emailVerifiedValue    = "email_verified", "true"
		givenNameAttr, familyNameAttr, emailAttr = "given_name", "family_name", "email"
	)
	for _, user := range m.Users {
		if user.ID == *input.Username {
			if user.Email == "internal.error@ons.gov.uk" {
				return nil, &smithy.GenericAPIError{
					Code:    errCodeInternalError,
					Message: "Something went wrong",
				}
			}
			return &cognitoidentityprovider.AdminGetUserOutput{
				UserAttributes: []types.AttributeType{
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
				Enabled:    user.Active,
				UserStatus: user.Status,
				Username:   aws.String(user.ID),
			}, nil
		}
	}
	return nil, &smithy.GenericAPIError{
		Code:    errCodeUserNotFound,
		Message: "the user could not be found",
	}
}

func (m *CognitoIdentityProviderClientStub) CreateGroup(_ context.Context, input *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
	userPoolID := "aaaa-bbbb-ccc-dddd"
	// non feature test functionality - input.GroupName starts with `test-group-` pattern
	// feature test functionality by default
	nonFeatureTesting := false
	if input.GroupName != nil {
		nonFeatureTesting, _ = regexp.MatchString("^test-group-.*", *input.GroupName)
	}
	var createGroupOutput *cognitoidentityprovider.CreateGroupOutput

	if nonFeatureTesting { // non feature test functionality
		if *input.GroupName == "internalError" {
			return nil, &smithy.GenericAPIError{
				Code:    errCodeInternalError,
				Message: "Something went wrong",
			}
		}

		for _, group := range m.Groups {
			if group.Name == *input.GroupName {
				return nil, &smithy.GenericAPIError{
					Code:    errCodeGroupExists,
					Message: "this group already exists",
				}
			}
		}

		newGroup, err := m.GenerateGroup(*input.GroupName, *input.Description, *input.Precedence)
		if err != nil {
			return nil, &smithy.GenericAPIError{
				Code:    errCodeInternalError,
				Message: err.Error(),
			}
		}
		m.Groups = append(m.Groups, newGroup)

		createGroupOutput = &cognitoidentityprovider.CreateGroupOutput{
			Group: &types.GroupType{
				Description:  input.Description,
				GroupName:    input.GroupName,
				Precedence:   input.Precedence,
				CreationDate: &newGroup.Created,
				UserPoolId:   &userPoolID,
			},
		}
	} else { // feature test functionality
		if *input.Description != "Internal Server Error" {
			// 201 response - group created
			response201 := `Thi$s is a te||st des$%£@^c ription for  a n ew group  $`
			createdTime, _ := time.Parse("2006-Jan-1", "2010-Jan-1")
			if *input.Description == response201 {
				createGroupOutput = &cognitoidentityprovider.CreateGroupOutput{
					Group: &types.GroupType{
						Description:  input.Description,
						GroupName:    input.GroupName,
						Precedence:   input.Precedence,
						CreationDate: &createdTime,
						UserPoolId:   &userPoolID,
					},
				}
			}
		} else {
			// 500 response - internal server error
			return nil, &smithy.GenericAPIError{
				Code:    errCodeInternalError,
				Message: "Something went wrong",
			}
		}
	}
	return createGroupOutput, nil
}

func (m *CognitoIdentityProviderClientStub) AdminUpdateUserAttributes(_ context.Context, input *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
	for _, user := range m.Users {
		if user.ID == *input.Username {
			if user.Email == "update.internalerror@ons.gov.uk" {
				return nil, &smithy.GenericAPIError{
					Code:    errCodeInternalError,
					Message: "Something went wrong",
				}
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
	return nil, &smithy.GenericAPIError{
		Code:    errCodeUserNotFound,
		Message: "the user could not be found",
	}
}

func (m *CognitoIdentityProviderClientStub) AdminEnableUser(_ context.Context, input *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
	for _, user := range m.Users {
		if user.ID == *input.Username {
			if user.Email == "enable.internalerror@ons.gov.uk" {
				return nil, &smithy.GenericAPIError{
					Code:    errCodeInternalError,
					Message: "Something went wrong whilst enabling",
				}
			}
			user.Active = true
			return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
		}
	}
	return nil, &smithy.GenericAPIError{
		Code:    errCodeUserNotFound,
		Message: "the user could not be found",
	}
}

func (m *CognitoIdentityProviderClientStub) AdminDisableUser(_ context.Context, input *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
	for _, user := range m.Users {
		if user.ID == *input.Username {
			if user.Email == "disable.internalerror@ons.gov.uk" {
				return nil, &smithy.GenericAPIError{
					Code:    errCodeInternalError,
					Message: "Something went wrong whilst disabling",
				}
			}
			user.Active = false
			return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
		}
	}
	return nil, &smithy.GenericAPIError{
		Code:    errCodeUserNotFound,
		Message: "the user could not be found",
	}
}

func (m *CognitoIdentityProviderClientStub) AdminAddUserToGroup(_ context.Context, input *cognitoidentityprovider.AdminAddUserToGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
	if *input.GroupName == internalError {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}

	group := m.ReadGroup(*input.GroupName)
	if group == nil {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "the group could not be found",
		}
	}

	user := m.ReadUser(*input.Username)
	if user == nil {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeUserNotFound,
			Message: "the user could not be found",
		}
	}

	user.Groups = append(user.Groups, group)
	group.Members = append(group.Members, user)

	return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
}

func (m *CognitoIdentityProviderClientStub) GetGroup(_ context.Context, input *cognitoidentityprovider.GetGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.GetGroupOutput, error) {
	if *input.GroupName == internalError || *input.GroupName == "get-group-internal-error" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}
	if *input.GroupName == "get-group-not-found" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "get group - group not found",
		}
	}

	group := m.ReadGroup(*input.GroupName)
	if group == nil {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "the group could not be found",
		}
	}
	timestamp := time.Now()
	return &cognitoidentityprovider.GetGroupOutput{
		Group: &types.GroupType{
			CreationDate:     &group.Created,
			Description:      &group.Description,
			GroupName:        &group.Name,
			LastModifiedDate: &timestamp,
			Precedence:       &group.Precedence,
			UserPoolId:       input.UserPoolId,
		},
	}, nil
}

func (m *CognitoIdentityProviderClientStub) ListUsersInGroup(_ context.Context, input *cognitoidentityprovider.ListUsersInGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
	var (
		emailVerifiedAttr, emailVerifiedValue    = "email_verified", "true"
		givenNameAttr, familyNameAttr, emailAttr = "given_name", "family_name", "email"
	)

	if *input.GroupName == internalError || *input.GroupName == "list-group-users-internal-error" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}
	if *input.GroupName == "list-group-users-not-found" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "list members - group not found",
		}
	}

	group := m.ReadGroup(*input.GroupName)
	if group == nil {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "the group could not be found",
		}
	}
	userList := make([]types.UserType, 0, len(group.Members))

	for _, user := range group.Members {
		userDetails := types.UserType{
			Attributes: []types.AttributeType{
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
			Enabled:    user.Active,
			UserStatus: user.Status,
			Username:   aws.String(user.ID),
		}
		userList = append(userList, userDetails)
	}
	return &cognitoidentityprovider.ListUsersInGroupOutput{
		Users: userList,
	}, nil
}

func (m *CognitoIdentityProviderClientStub) AdminRemoveUserFromGroup(_ context.Context, input *cognitoidentityprovider.AdminRemoveUserFromGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
	if *input.GroupName == internalError {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}

	group := m.ReadGroup(*input.GroupName)
	if group == nil {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "the group could not be found",
		}
	}

	user := m.ReadUser(*input.Username)
	if user == nil {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeUserNotFound,
			Message: "the user could not be found",
		}
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

// AdminConfirmSignUp - Added to fully implement interface but only used in the local dummy data builder
func (m *CognitoIdentityProviderClientStub) AdminConfirmSignUp(_ *cognitoidentityprovider.AdminConfirmSignUpInput) (
	*cognitoidentityprovider.AdminConfirmSignUpOutput, error) {
	return nil, nil
}

// AdminDeleteUser - Added to fully implement interface but only used in the local dummy data builder
func (m *CognitoIdentityProviderClientStub) AdminDeleteUser(_ context.Context, _ *cognitoidentityprovider.AdminDeleteUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDeleteUserOutput, error) {
	return nil, nil
}

// DeleteGroup was added to fully implement interface but is only used in the local dummy data builder
func (m *CognitoIdentityProviderClientStub) DeleteGroup(_ context.Context, input *cognitoidentityprovider.DeleteGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DeleteGroupOutput, error) {
	if *input.GroupName == internalError || *input.GroupName == "get-group-internal-error" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}
	if *input.GroupName == "delete-group-not-found" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "get group - group not found",
		}
	}
	return nil, nil
}

// AdminSetUserPassword - Added to fully implement interface but only used in the local dummy data builder
func (m *CognitoIdentityProviderClientStub) AdminSetUserPassword(_ context.Context, _ *cognitoidentityprovider.AdminSetUserPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminSetUserPasswordOutput, error) {
	return nil, nil
}

func (m *CognitoIdentityProviderClientStub) AdminListGroupsForUser(_ context.Context, input *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
	nextToken := "nextToken"
	nextTokenNil := ""
	if *input.Username == internalError || *input.Username == "get-group-internal-error" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}
	if *input.Username == "get-user-not-found" {
		println(errCodeUserNotFound)
		return nil, &smithy.GenericAPIError{
			Code:    errCodeUserNotFound,
			Message: "get user - user not found",
		}
	}
	if *input.UserPoolId == "get-user-pool-not-found" {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "get userpool  - userpool not found",
		}
	}
	user := m.ReadUser(*input.Username)
	if user == nil {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeUserNotFound,
			Message: "the user could not be found",
		}
	}

	newGroups := make([]types.GroupType, 0, len(user.Groups))

	if user.Groups == nil {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}

	for _, group := range user.Groups {
		newGroups = append(newGroups, types.GroupType{
			Description: &group.Description,
			GroupName:   &group.Name,
			Precedence:  &group.Precedence,
		})
	}

	if newGroups == nil {
		return &cognitoidentityprovider.AdminListGroupsForUserOutput{
			Groups:    newGroups,
			NextToken: &nextTokenNil,
		}, nil
	}

	if input.NextToken != nil && *input.NextToken != "" {
		return &cognitoidentityprovider.AdminListGroupsForUserOutput{
			Groups:    newGroups,
			NextToken: &nextToken,
		}, nil
	}
	return &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups:    newGroups,
		NextToken: &nextTokenNil,
	}, nil
}

func (m *CognitoIdentityProviderClientStub) ListGroups(_ context.Context, input *cognitoidentityprovider.ListGroupsInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListGroupsOutput, error) {
	if *input.UserPoolId == internalError {
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	}

	output := cognitoidentityprovider.ListGroupsOutput{}
	for _, group := range m.GroupsList {
		output.Groups = append(output.Groups, group.Groups...)
		output.NextToken = group.NextToken
	}
	return &output, nil
}

func (m *CognitoIdentityProviderClientStub) DescribeUserPoolClient(_ context.Context, _ *cognitoidentityprovider.DescribeUserPoolClientInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DescribeUserPoolClientOutput, error) {
	tokenValidDays := int32(1)
	refreshTokenUnits := types.TimeUnitsTypeDays
	userPoolClient := &cognitoidentityprovider.DescribeUserPoolClientOutput{
		UserPoolClient: &types.UserPoolClientType{
			RefreshTokenValidity: tokenValidDays,
			TokenValidityUnits: &types.TokenValidityUnitsType{
				RefreshToken: refreshTokenUnits,
			},
		},
	}
	return userPoolClient, nil
}

func (m *CognitoIdentityProviderClientStub) UpdateGroup(_ context.Context, input *cognitoidentityprovider.UpdateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.UpdateGroupOutput, error) {
	var (
		updateGroupOutput *cognitoidentityprovider.UpdateGroupOutput
		response200       = `Thi$s is a te||st des$%£@^c ription for  existing group  $`
		response200Up     = `Thi$s is a te||st des$%£@^c ription for  updated group  $`
		response500       = `Internal Server Error`
		userPoolID        = `aaaa-bbbb-ccc-dddd`
	)
	if *input.Description == response200 {
		// 200 response - group updated
		createdTime, _ := time.Parse("2006-Jan-1", "2010-Jan-1")
		groupName := "123e4567-e89b-12d3-a456-426614174000"
		updateGroupOutput = &cognitoidentityprovider.UpdateGroupOutput{
			Group: &types.GroupType{
				Description:  &response200,
				GroupName:    &groupName,
				CreationDate: &createdTime,
				UserPoolId:   &userPoolID,
			},
		}
		if input.Precedence != nil {
			updateGroupOutput.Group.Precedence = input.Precedence
		}
	} else if *input.Description == response200Up {
		// 200 response - group updated
		createdTime, _ := time.Parse("2006-Jan-1", "2010-Jan-1")
		groupName := "123e4567-e89b-12d3-a456-426614174000"
		updateGroupOutput = &cognitoidentityprovider.UpdateGroupOutput{
			Group: &types.GroupType{
				Description:  &response200Up,
				GroupName:    &groupName,
				CreationDate: &createdTime,
				UserPoolId:   &userPoolID,
			},
		}
		if input.Precedence != nil {
			updateGroupOutput.Group.Precedence = input.Precedence
		}
	} else if *input.Description == response500 {
		// 500 response - internal server error
		return nil, &smithy.GenericAPIError{
			Code:    errCodeInternalError,
			Message: "Something went wrong",
		}
	} else {
		// 404 response - resource not found error
		return nil, &smithy.GenericAPIError{
			Code:    errCodeResourceNotFound,
			Message: "Resource not found",
		}
	}
	return updateGroupOutput, nil
}
