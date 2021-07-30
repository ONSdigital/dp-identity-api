package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const addUserToGroupEndPoint = "http://localhost:25600/v1/groups/efgh5678/memebers"
const removeUserFromGroupEndPoint = "http://localhost:25600/v1/groups/efgh5678/memebers/abcd1234"
const getUsersInGroupEndPoint = "http://localhost:25600/v1/groups/efgh5678/members"

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
				"id":      "efgh5678",
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

			//postBody := map[string]interface{}{"user_id": userId}
			//body, _ := json.Marshal(postBody)
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
