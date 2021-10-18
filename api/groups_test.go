package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	addUserToGroupEndPoint      = "http://localhost:25600/v1/groups/efgh5678/memebers"
	removeUserFromGroupEndPoint = "http://localhost:25600/v1/groups/efgh5678/memebers/abcd1234"
	getUsersInGroupEndPoint     = "http://localhost:25600/v1/groups/efgh5678/members"
	createGroupEndPoint         = "http://localhost:25600/v1/groups"
	getListGroupsEndPoint       = "http://localhost:25600/v1/groups"
	updateGroupEndPoint         = "http://localhost:25600/v1/groups/123e4567-e89b-12d3-a456-426614174000"
)

var groupNotFoundDescription, internalErrorDescription string = "group not found", "internal error"

func TestAddUserToGroupHandler(t *testing.T) {

	var (
		ctx                                                                                = context.Background()
		userId                                                                      string = "abcd1234"
		userNotFoundDescription, groupNotFoundDescription, internalErrorDescription string = "user not found", "group not found", "internal error"
	)

	api, w, m := apiSetup()

	timeStamp := time.Now()
	getGroupData := &cognitoidentityprovider.GroupType{
		Description:  aws.String("a test group"),
		GroupName:    aws.String("test-group"),
		Precedence:   aws.Int64(100),
		CreationDate: &timeStamp,
	}

	Convey("Add a user to a group - check expected responses", t, func() {
		adminCreateUsersTests := []struct {
			addUserToGroupFunction    func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
			getGroupFunction          func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
			listUsersForGroupFunction func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			assertions                func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			// 200 response - user added to group
			{
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(errorResponse, ShouldBeNil)
				},
			},
			// Cognito 404 response - user not found
			{
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var userNotFoundException cognitoidentityprovider.UserNotFoundException
					userNotFoundException.Message_ = &userNotFoundDescription
					return nil, &userNotFoundException
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.UserNotFoundError)
					So(castErr.Description, ShouldEqual, userNotFoundDescription)
				},
			},
			// Cognito 404 response - group not found
			{
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, groupNotFoundDescription)
				},
			},
			// Cognito 500 response - internal error
			{
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)
				},
			},
			// Cognito GetGroup 404 response - internal error
			{
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{},
					}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, groupNotFoundDescription)
				},
			},
			// Cognito ListUsersInGroup 404 response - internal error
			{
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, groupNotFoundDescription)
				},
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.AdminAddUserToGroupFunc = tt.addUserToGroupFunction
			m.GetGroupFunc = tt.getGroupFunction
			m.ListUsersInGroupFunc = tt.listUsersForGroupFunction

			postBody := map[string]interface{}{"user_id": userId}
			body, _ := json.Marshal(postBody)
			r := httptest.NewRequest(http.MethodPost, addUserToGroupEndPoint, bytes.NewReader(body))

			urlVars := map[string]string{
				"id": "efgh5678",
			}
			r = mux.SetURLVars(r, urlVars)

			successResponse, errorResponse := api.AddUserToGroupHandler(ctx, w, r)

			tt.assertions(successResponse, errorResponse)
		}
	})

	Convey("Add a user to a group - returns 500 error unmarshalling invalid request body", t, func() {
		r := httptest.NewRequest(http.MethodPost, addUserToGroupEndPoint, bytes.NewReader(nil))

		successResponse, errorResponse := api.AddUserToGroupHandler(ctx, w, r)

		So(successResponse, ShouldBeNil)
		castErr := errorResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.JSONUnmarshalError)
		So(castErr.Description, ShouldEqual, models.ErrorUnmarshalFailedDescription)
	})

	Convey("Validation fails 400: validating user id throws validation errors", t, func() {
		userValidationTests := []struct {
			userDetails  map[string]interface{}
			errorCodes   []string
			httpResponse int
		}{
			// missing user id
			{
				map[string]interface{}{"user_id": ""},
				[]string{
					models.InvalidUserIdError,
				},
				http.StatusBadRequest,
			},
		}

		for _, tt := range userValidationTests {
			body, _ := json.Marshal(tt.userDetails)
			r := httptest.NewRequest(http.MethodPost, addUserToGroupEndPoint, bytes.NewReader(body))

			urlVars := map[string]string{
				"id": "efgh5678",
			}
			r = mux.SetURLVars(r, urlVars)

			successResponse, errorResponse := api.AddUserToGroupHandler(ctx, w, r)

			So(successResponse, ShouldBeNil)
			So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			castErr := errorResponse.Errors[0].(*models.Error)
			So(castErr.Code, ShouldEqual, tt.errorCodes[0])
			if len(errorResponse.Errors) > 1 {
				castErr = errorResponse.Errors[1].(*models.Error)
				So(castErr.Code, ShouldEqual, tt.errorCodes[1])
			}
		}
	})
}

func TestRemoveUserFromGroupHandler(t *testing.T) {
	var (
		ctx                                                                                = context.Background()
		userNotFoundDescription, groupNotFoundDescription, internalErrorDescription string = "user not found", "group not found", "internal error"
	)

	api, w, m := apiSetup()

	timeStamp := time.Now()
	getGroupData := &cognitoidentityprovider.GroupType{
		Description:  aws.String("a test group"),
		GroupName:    aws.String("test-group"),
		Precedence:   aws.Int64(100),
		CreationDate: &timeStamp,
	}

	Convey("Remove a user from a group - check expected responses", t, func() {
		adminCreateUsersTests := []struct {
			removeUserFromGroupFunction func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
			getGroupFunction            func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
			listUsersForGroupFunction   func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			assertions                  func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			// 202 response - user removed to group
			{
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(errorResponse, ShouldBeNil)
				},
			},
			// Cognito 404 response - user not found
			{
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					var userNotFoundException cognitoidentityprovider.UserNotFoundException
					userNotFoundException.Message_ = &userNotFoundDescription
					return nil, &userNotFoundException
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.UserNotFoundError)
					So(castErr.Description, ShouldEqual, userNotFoundDescription)
				},
			},
			// Cognito 404 response - group not found
			{
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, groupNotFoundDescription)
				},
			},
			// Cognito 500 response - internal error
			{
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)
				},
			},
			// Cognito GetGroup 404 response - internal error
			{
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{},
					}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, groupNotFoundDescription)
				},
			},
			// Cognito ListUsersInGroup 404 response - internal error
			{
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, groupNotFoundDescription)
				},
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.AdminRemoveUserFromGroupFunc = tt.removeUserFromGroupFunction
			m.GetGroupFunc = tt.getGroupFunction
			m.ListUsersInGroupFunc = tt.listUsersForGroupFunction

			r := httptest.NewRequest(http.MethodDelete, removeUserFromGroupEndPoint, bytes.NewReader(nil))

			urlVars := map[string]string{
				"id":      "efzgh5678",
				"user_id": "abcd1234",
			}
			r = mux.SetURLVars(r, urlVars)

			successResponse, errorResponse := api.RemoveUserFromGroupHandler(ctx, w, r)

			tt.assertions(successResponse, errorResponse)
		}
	})

	Convey("Validation fails 400: validating user id throws validation errors", t, func() {
		userValidationTests := []struct {
			urlVars      map[string]string
			errorCodes   []string
			httpResponse int
		}{
			// missing user id
			{
				map[string]string{"user_id": "", "id": "efgh5678"},
				[]string{
					models.InvalidUserIdError,
				},
				http.StatusBadRequest,
			},
			// missing group id
			{
				map[string]string{"user_id": "abcd1234", "id": ""},
				[]string{
					models.InvalidGroupNameError,
				},
				http.StatusBadRequest,
			},
			// missing group id and user id
			{
				map[string]string{"user_id": "", "id": ""},
				[]string{
					models.InvalidGroupNameError,
					models.InvalidUserIdError,
				},
				http.StatusBadRequest,
			},
		}

		for _, tt := range userValidationTests {
			r := httptest.NewRequest(http.MethodPost, removeUserFromGroupEndPoint, bytes.NewReader(nil))

			r = mux.SetURLVars(r, tt.urlVars)

			successResponse, errorResponse := api.RemoveUserFromGroupHandler(ctx, w, r)

			So(successResponse, ShouldBeNil)
			So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			castErr := errorResponse.Errors[0].(*models.Error)
			So(castErr.Code, ShouldEqual, tt.errorCodes[0])
			if len(errorResponse.Errors) > 1 {
				castErr = errorResponse.Errors[1].(*models.Error)
				So(castErr.Code, ShouldEqual, tt.errorCodes[1])
			}
		}
	})
}

func TestGetUsersFromGroupHandler(t *testing.T) {

	var (
		ctx                                                       = context.Background()
		groupNotFoundDescription, internalErrorDescription string = "group not found", "internal error"
	)

	api, w, m := apiSetup()

	Convey("adds the returned users to the user list and sets the count", t, func() {
		cognitoResponse := cognitoidentityprovider.ListUsersInGroupOutput{
			Users: []*cognitoidentityprovider.UserType{
				{
					Enabled:    aws.Bool(true),
					UserStatus: aws.String("CONFIRMED"),
					Username:   aws.String("user-1"),
				},
				{
					Enabled:    aws.Bool(true),
					UserStatus: aws.String("CONFIRMED"),
					Username:   aws.String("user-2"),
				},
			},
		}
		listOfUsers := models.UsersList{}
		listOfUsers.MapCognitoUsers(&cognitoResponse.Users)

		So(len(listOfUsers.Users), ShouldEqual, len(cognitoResponse.Users))
		So(listOfUsers.Count, ShouldEqual, len(cognitoResponse.Users))
	})

	Convey("and the expected responses", t, func() {
		listUsersInGroupTests := []struct {
			listUsersForGroupFunction func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			assertions                func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			// 200 response - user added to group
			{
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(errorResponse, ShouldBeNil)
				},
			},
			// Cognito 404 response - group not found
			{
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, groupNotFoundDescription)
				},
			},
			// Cognito 500 response - internal error
			{
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)
				},
			},
		}

		for _, tt := range listUsersInGroupTests {

			m.ListUsersInGroupFunc = tt.listUsersForGroupFunction

			r := httptest.NewRequest(http.MethodGet, getUsersInGroupEndPoint, nil)

			urlVars := map[string]string{
				"id": "efgh5678",
			}
			r = mux.SetURLVars(r, urlVars)

			successResponse, errorResponse := api.ListUsersInGroupHandler(ctx, w, r)

			tt.assertions(successResponse, errorResponse)
		}
	})
}

func TestGetUsersInAGroup(t *testing.T) {

	var (
		groupNotFoundDescription string = "group not found"
		name                     string = "name"
	)

	getGroupData := models.Group{
		Name: "test-group",
	}

	api, _, m := apiSetup()
	Convey("error is returned when list users in group returns an error", t, func() {
		m.ListUsersInGroupFunc = func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
			var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
			groupNotFoundException.Message_ = &groupNotFoundDescription
			return nil, &groupNotFoundException
		}

		listOfUsersResponse, errorResponse := api.getUsersInAGroup(nil, getGroupData)

		So(listOfUsersResponse, ShouldBeNil)
		So(errorResponse.Error(), ShouldResemble, "ResourceNotFoundException: group not found")
	})

	Convey("When there is no next token cognito is called once and the list of users in returned", t, func() {
		listOfUsers := []*cognitoidentityprovider.UserType{
			{
				Username: &name,
			},
		}

		m.ListUsersInGroupFunc = func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
			listUsersInGroup := &cognitoidentityprovider.ListUsersInGroupOutput{
				Users: []*cognitoidentityprovider.UserType{
					{
						Username: &name,
					},
				},
			}
			return listUsersInGroup, nil
		}

		listOfUsersResponse, errorResponse := api.getUsersInAGroup(nil, getGroupData)

		So(listOfUsersResponse, ShouldResemble, listOfUsers)
		So(errorResponse, ShouldBeNil)

	})

	Convey("When there is a next token cognito is called more than once and the appended list of users in returned", t, func() {
		listOfUsers := []*cognitoidentityprovider.UserType{
			{
				Username: &name,
			},
			{
				Username: &name,
			},
		}

		m.ListUsersInGroupFunc = func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
			nextToken := "nextToken"

			if input.NextToken != nil {
				listUsersInGroup := &cognitoidentityprovider.ListUsersInGroupOutput{
					NextToken: nil,
					Users: []*cognitoidentityprovider.UserType{
						{
							Username: &name,
						},
					},
				}
				return listUsersInGroup, nil
			} else {
				listUsersInGroup := &cognitoidentityprovider.ListUsersInGroupOutput{
					NextToken: &nextToken,
					Users: []*cognitoidentityprovider.UserType{
						{
							Username: &name,
						},
					},
				}
				return listUsersInGroup, nil
			}
		}

		listOfUsersResponse, errorResponse := api.getUsersInAGroup(nil, getGroupData)

		So(listOfUsersResponse, ShouldResemble, listOfUsers)
		So(errorResponse, ShouldBeNil)

	})

	Convey("When GetUsersInAGroup in called with a list of users the appended list of users in returned", t, func() {

		listOfUsers := []*cognitoidentityprovider.UserType{
			{
				Username: &name,
			},
		}

		returnedListOfUsers := []*cognitoidentityprovider.UserType{
			{
				Username: &name,
			},
			{
				Username: &name,
			},
		}

		m.ListUsersInGroupFunc = func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
			listUsersInGroup := &cognitoidentityprovider.ListUsersInGroupOutput{
				Users: []*cognitoidentityprovider.UserType{
					{
						Username: &name,
					},
				},
			}
			return listUsersInGroup, nil
		}

		listOfUsersResponse, errorResponse := api.getUsersInAGroup(listOfUsers, getGroupData)

		So(listOfUsersResponse, ShouldResemble, returnedListOfUsers)
		So(errorResponse, ShouldBeNil)
	})
}

func TestCreateNewGroup(t *testing.T) {
	var (
		internalErrorDescription string = "internal error"
	)

	api, w, m := apiSetup()

	// ListGroupsFunction template - success
	listGroupsFuncSuccess := func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error){
		d := "thisisamocktestname"
		g := "123e4567-e89b-12d3-a456-426614174000"
		p := int64(12)
		groupsList := cognitoidentityprovider.ListGroupsOutput{
			NextToken: nil,
			Groups: []*cognitoidentityprovider.GroupType{
				{
					Description: &d,
					GroupName:   &g,
					Precedence:  &p,
				},
			},
		}

		return &groupsList, nil
	}

	Convey("Create a new group - check responses", t, func() {
		createGroupTests := []struct {
			createNewGroupFunction func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error)
			listGroupsFunction     func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error)
			createGroupInput,
			expectedResponse map[string]interface{}
			assertions func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			// 201 response - group created
			{
				func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error) {
					return &cognitoidentityprovider.CreateGroupOutput{}, nil
				},
				listGroupsFuncSuccess,
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  22,
				},
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  22,
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusCreated)
					So(errorResponse, ShouldBeNil)
				},
			},
			// 400 response - no description field in request body
			{
				nil,
				listGroupsFuncSuccess,
				map[string]interface{}{
					"precedence": 22,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupName)
					So(castErr.Description, ShouldEqual, models.MissingGroupName)
				},
			},
			// 400 response - no precedence field in request body
			{
				nil,
				listGroupsFuncSuccess,
				map[string]interface{}{
					"name": "This is a test name",
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupPrecedence)
					So(castErr.Description, ShouldEqual, models.MissingGroupPrecedence)
				},
			},
			// 400 response - group description begins with reserved string `role-`
			{
				nil,
				listGroupsFuncSuccess,
				map[string]interface{}{
					"name": "role-This is a test name",
					"precedence":  22,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupName)
					So(castErr.Description, ShouldEqual, models.IncorrectPatternInGroupName)
				},
			},
			// 400 response - group precedence setting not minimum of `10`
			{
				nil,
				listGroupsFuncSuccess,
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  1,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupPrecedence)
					So(castErr.Description, ShouldEqual, models.GroupPrecedenceIncorrect)
				},
			},
			// 400 response - group name already exists
			{
				nil,
				func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error){
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				map[string]interface{}{
					"name": "This&^ is- a MOCK. test**() NAMe",
					"precedence":  12,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)
				},
			},
			// 500 response - internal server error from Cognito
			{
				func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				listGroupsFuncSuccess,
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  12,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)
				},
			},
		}

		for _, tt := range createGroupTests {
			m.CreateGroupFunc = tt.createNewGroupFunction
			m.ListGroupsFunc = tt.listGroupsFunction
			body, _ := json.Marshal(tt.createGroupInput)
			r := httptest.NewRequest(http.MethodPost, createGroupEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.CreateGroupHandler(context.Background(), w, r)

			tt.assertions(successResponse, errorResponse)
		}
	})
}

func TestUpdateGroup(t *testing.T) {
	var (
		internalErrorDescription, notFoundErrorDescription string = "internal error", "not found error"
	)

	api, w, m := apiSetup()

	Convey("Update a group - check responses", t, func() {
		createGroupTests := []struct {
			updateGroupFunction func(input *cognitoidentityprovider.UpdateGroupInput) (*cognitoidentityprovider.UpdateGroupOutput, error)
			updateGroupInput,
			expectedResponse map[string]interface{}
			assertions func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			// 200 response - group updated
			{
				func(input *cognitoidentityprovider.UpdateGroupInput) (*cognitoidentityprovider.UpdateGroupOutput, error) {
					return &cognitoidentityprovider.UpdateGroupOutput{}, nil
				},
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  22,
				},
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  22,
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(errorResponse, ShouldBeNil)
				},
			},
			// 400 response - no description field in request body
			{
				nil,
				map[string]interface{}{
					"precedence": 22,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupName)
					So(castErr.Description, ShouldEqual, models.MissingGroupName)
				},
			},
			// 400 response - no precedence field in request body
			{
				nil,
				map[string]interface{}{
					"name": "This is a test name",
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupPrecedence)
					So(castErr.Description, ShouldEqual, models.MissingGroupPrecedence)
				},
			},
			// 400 response - group description begins with reserved string `role-`
			{
				nil,
				map[string]interface{}{
					"name": "role-This is a test name",
					"precedence":  22,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupName)
					So(castErr.Description, ShouldEqual, models.IncorrectPatternInGroupName)
				},
			},
			// 400 response - group precedence setting not minimum of `10`
			{
				nil,
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  1,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupPrecedence)
					So(castErr.Description, ShouldEqual, models.GroupPrecedenceIncorrect)
				},
			},
			// 500 response - internal server error from Cognito
			{
				func(input *cognitoidentityprovider.UpdateGroupInput) (*cognitoidentityprovider.UpdateGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  12,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)
				},
			},
			// 404 response - resource not found from Cognito
			{
				func(input *cognitoidentityprovider.UpdateGroupInput) (*cognitoidentityprovider.UpdateGroupOutput, error) {
					var notFoundError cognitoidentityprovider.ResourceNotFoundException
					notFoundError.Message_ = &notFoundErrorDescription
					return nil, &notFoundError
				},
				map[string]interface{}{
					"name": "This is a test name",
					"precedence":  12,
				},
				nil,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, notFoundErrorDescription)
				},
			},
		}

		for _, tt := range createGroupTests {
			m.UpdateGroupFunc = tt.updateGroupFunction
			body, _ := json.Marshal(tt.updateGroupInput)
			r := httptest.NewRequest(http.MethodPut, updateGroupEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.UpdateGroupHandler(context.Background(), w, r)

			tt.assertions(successResponse, errorResponse)
		}
	})
}

func TestGetListGroups(t *testing.T) {

	api, _, m := apiSetup()

	Convey("When there is no next token cognito is called once and an empty list of groups is returned", t, func() {
		listOfGroups := []*cognitoidentityprovider.GroupType{
			{},
		}
		var count int = 0
		m.ListGroupsFunc = func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
			count++
			listGroups := &cognitoidentityprovider.ListGroupsOutput{
				NextToken: nil,
				Groups: []*cognitoidentityprovider.GroupType{
					{},
				},
			}
			return listGroups, nil
		}

		listOfGroupsResponse, errorResponse := api.GetListGroups()

		So(errorResponse, ShouldBeNil)

		So(listOfGroupsResponse.Groups, ShouldResemble, listOfGroups)
		So(listOfGroupsResponse.Groups, ShouldHaveLength, len(listOfGroups))
		So(listOfGroupsResponse.NextToken, ShouldBeNil)
		So(count, ShouldEqual, 1)

	})

	Convey("When there is no next token cognito is called with 1  entry list of groups in returned", t, func() {
		var (
			description, group_name string = "The publishing admins", "role-admin"
			precedence              int64  = 1
			count                   int    = 0
		)
		listOfGroups := []*cognitoidentityprovider.GroupType{
			{
				Description: &description,
				GroupName:   &group_name,
				Precedence:  &precedence,
			},
		}

		m.ListGroupsFunc = func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
			count++
			listGroups := &cognitoidentityprovider.ListGroupsOutput{

				NextToken: nil,
				Groups: []*cognitoidentityprovider.GroupType{
					{
						Description: &description,
						GroupName:   &group_name,
						Precedence:  &precedence,
					},
				},
			}
			return listGroups, nil
		}

		listOfGroupsResponse, errorResponse := api.GetListGroups()

		So(errorResponse, ShouldBeNil)
		So(listOfGroupsResponse.NextToken, ShouldBeNil)
		So(listOfGroupsResponse.Groups, ShouldResemble, listOfGroups)
		So(listOfGroupsResponse.Groups, ShouldHaveSameTypeAs, listOfGroups)
		So(listOfGroupsResponse.Groups, ShouldHaveLength, len(listOfGroups))
		So(count, ShouldEqual, 1)

	})
}
func TestListGroupsHandler(t *testing.T) {

	var (
		ctx       = context.Background()
		timestamp = time.Now()
		// internalErrorDescription string = "internal error"
		// next_token                      = "next_token"
		groups = []*cognitoidentityprovider.GroupType{
			{
				CreationDate:     &timestamp,
				Description:      aws.String("A test group1"),
				GroupName:        aws.String("test-group1"),
				LastModifiedDate: &timestamp,
				Precedence:       aws.Int64(4),
				RoleArn:          aws.String(""),
				UserPoolId:       aws.String(""),
			},
			{
				CreationDate:     &timestamp,
				Description:      aws.String("A test group1"),
				GroupName:        aws.String("test-group1"),
				LastModifiedDate: &timestamp,
				Precedence:       aws.Int64(4),
				RoleArn:          aws.String(""),
				UserPoolId:       aws.String(""),
			},
		}
	)

	api, w, m := apiSetup()

	Convey("List groups -check expected responses", t, func() {
		internalErrorDescription := ""
		listGroupsTest := []struct {
			description           string
			next_token            string
			getListGroupsFunction func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error)
			assertions            func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				"200 response from Cognito with empty NextToken",
				"",
				func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
					return &cognitoidentityprovider.ListGroupsOutput{
						Groups:    groups,
						NextToken: nil,
					}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {

					So(errorResponse, ShouldBeNil)

					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Body, ShouldNotBeNil)

					var responseBody = models.ListUserGroups{}
					json.Unmarshal(successResponse.Body, &responseBody)

					So(responseBody.NextToken, ShouldBeNil)
					So(responseBody.Count, ShouldEqual, 2)
					So(responseBody.Groups, ShouldNotBeNil)
					So(responseBody.Groups, ShouldHaveLength, responseBody.Count)
					So(*responseBody.Groups[0].Description, ShouldEqual, *groups[0].Description)
				},
			},
			{
				"200 response from Cognito with no groups",
				"",
				func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
					return &cognitoidentityprovider.ListGroupsOutput{
						Groups:    []*cognitoidentityprovider.GroupType{},
						NextToken: nil,
					}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {

					So(errorResponse, ShouldBeNil)

					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Body, ShouldNotBeNil)

					var responseBody = models.ListUserGroups{}
					json.Unmarshal(successResponse.Body, &responseBody)
					So(responseBody.NextToken, ShouldBeNil)
					So(responseBody.Count, ShouldEqual, 0)
					So(responseBody.Groups, ShouldBeNil)
					So(responseBody.Groups, ShouldHaveLength, responseBody.Count)
				},
			},
			{
				"200 response from Cognito with populated NextToken",
				"next_token",
				func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
					return &cognitoidentityprovider.ListGroupsOutput{
						Groups:    groups,
						NextToken: nil,
					}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {

					So(errorResponse, ShouldBeNil)

					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Body, ShouldNotBeNil)

					var responseBody = models.ListUserGroups{}
					json.Unmarshal(successResponse.Body, &responseBody)

					So(responseBody.NextToken, ShouldBeNil)
					So(responseBody.Count, ShouldEqual, 2)
					So(responseBody.Groups, ShouldNotBeNil)
					So(responseBody.Groups, ShouldHaveLength, responseBody.Count)
					So(*responseBody.Groups[0].Description, ShouldEqual, *groups[0].Description)
				},
			},

			{
				"500 response from Cognito",
				"",
				func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := internalErrorDescription
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)

					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)

					So(successResponse, ShouldBeNil)
				},
			},
		}

		for _, tt := range listGroupsTest {
			Convey(tt.description, func() {
				m.ListGroupsFunc = tt.getListGroupsFunction

				postBody := map[string]interface{}{"NextToken": tt.next_token}
				body, err := json.Marshal(postBody)
				So(err, ShouldBeNil)

				r := httptest.NewRequest(http.MethodGet, getListGroupsEndPoint, bytes.NewReader(body))

				urlVars := map[string]string{
					"id": "efgh5678",
				}
				r = mux.SetURLVars(r, urlVars)

				successResponse, errorResponse := api.ListGroupsHandler(ctx, w, r)

				tt.assertions(successResponse, errorResponse)
			})
		}
	})
}

func TestGetGroupHandler(t *testing.T) {

	var (
		ctx       = context.Background()
		timestamp = time.Now()
		getgroup  = cognitoidentityprovider.GroupType{
			CreationDate:     &timestamp,
			Description:      aws.String("A test group1"),
			GroupName:        aws.String("test-group1"),
			LastModifiedDate: &timestamp,
			Precedence:       aws.Int64(4),
			RoleArn:          aws.String(""),
			UserPoolId:       aws.String(""),
		}
	)

	api, w, m := apiSetup()

	Convey("Get group -check expected responses", t, func() {
		GetGroupTest := []struct {
			description      string
			getGroupFunction func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
			assertions       func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				"200 response from Cognito ",
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					return &cognitoidentityprovider.GetGroupOutput{
						Group: &getgroup,
					}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {

					So(errorResponse, ShouldBeNil)

					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Body, ShouldNotBeNil)

					var responseBody = models.ListUserGroups{}
					json.Unmarshal(successResponse.Body, &responseBody)

					So(responseBody, ShouldNotBeNil)

				},
			},
			{
				"404 response from Cognito ",
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {

					So(errorResponse, ShouldNotBeNil)

					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
					So(successResponse, ShouldBeNil)
				},
			},
			{
				"500 response from Cognito ",
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {

					So(errorResponse, ShouldNotBeNil)

					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					So(successResponse, ShouldBeNil)
				},
			}}

		for _, tt := range GetGroupTest {
			Convey(tt.description, func() {
				m.GetGroupFunc = tt.getGroupFunction

				postBody := map[string]interface{}{"GroupName": "group_name_test"}
				body, err := json.Marshal(postBody)
				So(err, ShouldBeNil)

				r := httptest.NewRequest(http.MethodGet, getListGroupsEndPoint, bytes.NewReader(body))

				urlVars := map[string]string{
					"id": "efgh5678",
				}
				r = mux.SetURLVars(r, urlVars)

				successResponse, errorResponse := api.GetGroupHandler(ctx, w, r)

				tt.assertions(successResponse, errorResponse)
			})
		}
	})

}

func TestDeleteGroupHandler(t *testing.T) {

	var (
		ctx = context.Background()
	)

	api, w, m := apiSetup()

	Convey("Delete group -check expected responses", t, func() {
		DeleteGroupTest := []struct {
			description         string
			DeleteGroupFunction func(input *cognitoidentityprovider.DeleteGroupInput) (*cognitoidentityprovider.DeleteGroupOutput, error)
			assertions          func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{{
			"204 response from Cognito ",
			func(input *cognitoidentityprovider.DeleteGroupInput) (*cognitoidentityprovider.DeleteGroupOutput, error) {
				return &cognitoidentityprovider.DeleteGroupOutput{}, nil
			},
			func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
				So(errorResponse, ShouldBeNil)
				So(successResponse.Status, ShouldEqual, http.StatusNoContent)
				So(successResponse, ShouldNotBeNil)
			},
		},
			{
				"404 response from Cognito ",
				func(input *cognitoidentityprovider.DeleteGroupInput) (*cognitoidentityprovider.DeleteGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {

					So(errorResponse, ShouldNotBeNil)

					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
					So(successResponse, ShouldBeNil)
				},
			},
			{
				"500 response from Cognito ",
				func(input *cognitoidentityprovider.DeleteGroupInput) (*cognitoidentityprovider.DeleteGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {

					So(errorResponse, ShouldNotBeNil)

					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					So(successResponse, ShouldBeNil)
				},
			}}

		for _, tt := range DeleteGroupTest {
			Convey(tt.description, func() {
				m.DeleteGroupFunc = tt.DeleteGroupFunction

				postBody := map[string]interface{}{"GroupName": "group_name_test"}
				body, err := json.Marshal(postBody)
				So(err, ShouldBeNil)

				r := httptest.NewRequest(http.MethodGet, getListGroupsEndPoint, bytes.NewReader(body))

				urlVars := map[string]string{
					"id": "efgh5678",
				}
				r = mux.SetURLVars(r, urlVars)

				successResponse, errorResponse := api.DeleteGroupHandler(ctx, w, r)

				tt.assertions(successResponse, errorResponse)
			})
		}
	})

}
