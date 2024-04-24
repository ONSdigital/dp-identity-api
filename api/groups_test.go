package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
	getGroupsReportEndPoint     = "http://localhost:25600/v1/groups/groups-report"
	addUserToGroupEndPoint      = "http://localhost:25600/v1/groups/efgh5678/members"
	removeUserFromGroupEndPoint = "http://localhost:25600/v1/groups/efgh5678/members/abcd1234"
	getUsersInGroupEndPoint     = "http://localhost:25600/v1/groups/efgh5678/members"
	createGroupEndPoint         = "http://localhost:25600/v1/groups"
	getListGroupsEndPoint       = "http://localhost:25600/v1/groups"
	updateGroupEndPoint         = "http://localhost:25600/v1/groups/123e4567-e89b-12d3-a456-426614174000"
)

var (
	groupNotFoundDescription,
	internalErrorDescription,
	userNotFoundDescription = "group not found", "internal error", "user not found"
	ctx = context.Background()
)

func TestAddUserToGroupHandler(t *testing.T) {
	var (
		userId = "abcd1234"
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
		addUserToGroupTests := []struct {
			description               string
			addUserToGroupFunction    func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
			getGroupFunction          func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
			listUsersForGroupFunction func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			assertions                func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{

			{
				"200 response - user added to group",
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
			{
				"Cognito 404 response - getGroup group not found",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldResemble, groupNotFoundDescription)
				},
			},
			{
				"Cognito 500 response - getGroup internal error",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var exception cognitoidentityprovider.InternalErrorException
					exception.Message_ = &internalErrorDescription
					return nil, &exception
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
			{
				"500 response - addUserToGroup",
				func(
					userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (
					*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (
					*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
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
			{
				"400 response - addUserToGroup group not found",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (
					*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (
					*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldResemble, groupNotFoundDescription)
				},
			},
			{
				"404 response - addUserToGroup user not found",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (
					*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var userNotFoundException cognitoidentityprovider.UserNotFoundException
					userNotFoundException.Message_ = &userNotFoundDescription
					return nil, &userNotFoundException
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (
					*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.UserNotFoundError)
					So(castErr.Description, ShouldResemble, userNotFoundDescription)
				},
			},
		}
		for _, tt := range addUserToGroupTests {
			Convey(tt.description, func() {
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
			})
		}
	})
	Convey("Add a user to a group - check expected responses", t, func() {
		addUserToGroupTests := []struct {
			description            string
			addUserToGroupFunction func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (
				*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
			getGroupFunction func(input *cognitoidentityprovider.GetGroupInput) (
				*cognitoidentityprovider.GetGroupOutput, error)
			listUsersForGroupFunction func(input *cognitoidentityprovider.ListUsersInGroupInput) (
				*cognitoidentityprovider.ListUsersInGroupOutput, error)
			userID     string
			groupID    string
			assertions func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				"Cognito 400 response - User validation internal error",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (
					*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (
					*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				"",
				"test_group",
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidUserIdError)
					So(castErr.Description, ShouldEqual, models.MissingUserIdErrorDescription)
				},
			},
			{
				"Cognito 400 response - group validation internal error",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (
					*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (
					*cognitoidentityprovider.GetGroupOutput, error) {
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: getGroupData,
					}
					return group, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				"test_user",
				"",
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupIDError)
					So(castErr.Description, ShouldEqual, models.MissingGroupIDErrorDescription)
				},
			},
		}
		for _, tt := range addUserToGroupTests {
			Convey(tt.description, func() {
				m.AdminAddUserToGroupFunc = tt.addUserToGroupFunction
				m.GetGroupFunc = tt.getGroupFunction
				m.ListUsersInGroupFunc = tt.listUsersForGroupFunction
				postBody := map[string]interface{}{"user_id": tt.userID}
				body, _ := json.Marshal(postBody)
				r := httptest.NewRequest(http.MethodPost, addUserToGroupEndPoint, bytes.NewReader(body))
				urlVars := map[string]string{
					"id": tt.groupID,
				}
				r = mux.SetURLVars(r, urlVars)
				successResponse, errorResponse := api.AddUserToGroupHandler(ctx, w, r)
				tt.assertions(successResponse, errorResponse)

			})
		}
	})
}

func TestRemoveUserFromGroupHandler(t *testing.T) {
	api, w, m := apiSetup()
	timeStamp := time.Now()
	getGroupData := &cognitoidentityprovider.GroupType{
		Description:  aws.String("a test group"),
		GroupName:    aws.String("test-group"),
		Precedence:   aws.Int64(100),
		CreationDate: &timeStamp,
	}
	Convey("Remove a user from a group - check expected responses", t, func() {
		removeUsersFromGroupTests := []struct {
			description                 string
			removeUserFromGroupFunction func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
			getGroupFunction            func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
			listUsersForGroupFunction   func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			assertions                  func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{

			{
				"202 response - user removed to group",
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
			{
				"Cognito 404 response - getGroup group not found",
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldResemble, groupNotFoundDescription)
				},
			},
			{
				"Cognito 500 response - getGroup internal error",
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(inputData *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var exception cognitoidentityprovider.InternalErrorException
					exception.Message_ = &internalErrorDescription
					return nil, &exception
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
			{
				"500 response - removeUserfromGroup",
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
			{
				"400 response - removeUserfromGroup group not found",
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
					So(castErr.Description, ShouldResemble, groupNotFoundDescription)
				},
			},
			{
				"404 response - removeUserfromGroup user not found",
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
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.UserNotFoundError)
					So(castErr.Description, ShouldResemble, userNotFoundDescription)
				},
			},
		}

		for _, tt := range removeUsersFromGroupTests {
			Convey(tt.description, func() {
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
			})
		}
	})
	Convey("Remove a user from a group - check expected responses", t, func() {
		removeUsersFromGroupTests := []struct {
			description                 string
			removeUserFromGroupFunction func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
			getGroupFunction            func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
			listUsersForGroupFunction   func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			userID                      string
			groupID                     string
			assertions                  func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				"Cognito 400 response - User validation internal error",
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
				"",
				"test_group",
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidUserIdError)
					So(castErr.Description, ShouldEqual, models.MissingUserIdErrorDescription)
				},
			},
			{
				"Cognito 400 response - group validation internal error",
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
				"test_user",
				"",
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidGroupIDError)
					So(castErr.Description, ShouldEqual, models.MissingGroupIDErrorDescription)
				},
			},
		}

		for _, tt := range removeUsersFromGroupTests {
			Convey(tt.description, func() {
				m.AdminRemoveUserFromGroupFunc = tt.removeUserFromGroupFunction
				m.GetGroupFunc = tt.getGroupFunction
				m.ListUsersInGroupFunc = tt.listUsersForGroupFunction

				r := httptest.NewRequest(http.MethodDelete, removeUserFromGroupEndPoint, bytes.NewReader(nil))

				urlVars := map[string]string{
					"id":      tt.groupID,
					"user_id": tt.userID,
				}
				r = mux.SetURLVars(r, urlVars)

				successResponse, errorResponse := api.RemoveUserFromGroupHandler(ctx, w, r)

				tt.assertions(successResponse, errorResponse)
			})
		}
	})
}

func TestGetUsersFromGroupHandler(t *testing.T) {
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
		name = "name"
	)

	getGroupData := models.Group{
		ID: "test-group",
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
		internalErrorDescription = "internal error"
	)

	api, w, m := apiSetup()

	// ListGroupsFunction template - success
	listGroupsFuncSuccess := func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
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
					"name":       "This is a test name",
					"precedence": 22,
				},
				map[string]interface{}{
					"name":       "This is a test name",
					"precedence": 22,
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
					"name":       "role-This is a test name",
					"precedence": 22,
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
					"name":       "This is a test name",
					"precedence": 1,
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
				func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				map[string]interface{}{
					"name":       "This&^ is- a MOCK. test**() NAMe",
					"precedence": 12,
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
					"name":       "This is a test name",
					"precedence": 12,
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
		internalErrorDescription, notFoundErrorDescription = "internal error", "not found error"
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
					"name":       "This is a test name",
					"precedence": 22,
				},
				map[string]interface{}{
					"name":       "This is a test name",
					"precedence": 22,
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(errorResponse, ShouldBeNil)
				},
			},
			// 200 response - group updated, no precedence
			{
				func(input *cognitoidentityprovider.UpdateGroupInput) (*cognitoidentityprovider.UpdateGroupOutput, error) {
					return &cognitoidentityprovider.UpdateGroupOutput{}, nil
				},
				map[string]interface{}{
					"name": "This is a test name",
				},
				map[string]interface{}{
					"name": "This is a test name",
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
			// 400 response - group description begins with reserved string `role-`
			{
				nil,
				map[string]interface{}{
					"name":       "role-This is a test name",
					"precedence": 22,
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
					"name":       "This is a test name",
					"precedence": 1,
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
					"name":       "This is a test name",
					"precedence": 12,
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
					"name":       "This is a test name",
					"precedence": 12,
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
		var count = 0
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

		listOfGroupsResponse, errorResponse := api.GetListGroups("")

		So(errorResponse, ShouldBeNil)

		So(listOfGroupsResponse.Groups, ShouldResemble, listOfGroups)
		So(listOfGroupsResponse.Groups, ShouldHaveLength, len(listOfGroups))
		So(listOfGroupsResponse.NextToken, ShouldBeNil)
		So(count, ShouldEqual, 1)

	})

	Convey("When there is no next token cognito is called with 1  entry list of groups in returned", t, func() {
		var (
			description, group_name       = "The publishing admins", "role-admin"
			precedence              int64 = 1
			count                         = 0
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

		listOfGroupsResponse, errorResponse := api.GetListGroups("")

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
		timestamp = time.Now()
		groups    = []*cognitoidentityprovider.GroupType{
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
					So(*responseBody.Groups[0].Name, ShouldEqual, *groups[0].Description)
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
					So(*responseBody.Groups[0].Name, ShouldEqual, *groups[0].Description)
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

func TestSetGroupUsersHandler(t *testing.T) {

	var (
		name1     = "user-1"
		name2     = "user-2"
		name3     = "user-3"
		timestamp = time.Now()
		getgroup  = cognitoidentityprovider.GroupType{
			CreationDate:     &timestamp,
			Description:      aws.String("A test group1"),
			GroupName:        aws.String("test-group1"),
			LastModifiedDate: &timestamp,
			Precedence:       aws.Int64(4),
			RoleArn:          aws.String(""),
			UserPoolId:       aws.String("")}
	)

	api, w, m := apiSetup()

	Convey("Get group -check expected responses", t, func() {

		GetGroupTest := []struct {
			description                    string
			postbody                       []map[string]string
			mock_getGroupfunc              func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
			mock_listUsersInGroupfunc      func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			mock_setGroupUsersfunc         func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse)
			mockAddUserToGroupFunction     func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse)
			mock_addUserToGroupfunc        func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
			mockRemoveUserToGroupFunction  func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse)
			mock_removeUserToGroupFunction func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
			assertions                     func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				"200 response from Cognito  with input and output",
				[]map[string]string{
					{"user_id": name1},
					{"user_id": name2},
				},
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					return &cognitoidentityprovider.GetGroupOutput{
						Group: &getgroup,
					}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String(name2),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String(name3),
							},
						},
					}, nil
				},
				func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, &models.ErrorResponse{}
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{
							Users: []models.UserParams{{ID: name1}, {ID: name2}},
							Count: 1,
						},
						nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(errorResponse, ShouldBeNil)
					var responseBody = models.UsersList{}
					json.Unmarshal(successResponse.Body, &responseBody)
					So(responseBody, ShouldNotBeNil)
					So(responseBody.Count, ShouldEqual, 2)
				},
			},
			{
				"200 response from Cognito  with input and output zero input",
				nil,
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					return &cognitoidentityprovider.GetGroupOutput{
						Group: &getgroup,
					}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{},
					}, nil
				},
				func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, &models.ErrorResponse{}
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(errorResponse, ShouldBeNil)
					var responseBody = models.UsersList{}
					json.Unmarshal(successResponse.Body, &responseBody)
					So(responseBody, ShouldNotBeNil)
					So(responseBody.Count, ShouldEqual, 0)
				},
			},
			{
				"400 response from Cognito user does not exits",
				[]map[string]string{
					{"user": name1},
					{"user": name2},
				},
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					return &cognitoidentityprovider.GetGroupOutput{
						Group: &getgroup,
					}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{},
					}, nil
				},
				func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
					errorResponse := models.ErrorResponse{
						Errors: []error{models.NewValidationError(ctx, models.InvalidUserIdError, userNotFoundDescription)},
						Status: http.StatusNotFound,
					}
					return nil, &errorResponse
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, 400)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidUserIdError)
					So(castErr.Description, ShouldEqual, "the user id was missing")

				},
			},
			{
				"500 response from Cognito ",
				nil,
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					return &cognitoidentityprovider.GetGroupOutput{
						Group: &getgroup,
					}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := internalErrorDescription
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, &models.ErrorResponse{}
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)
				},
			},
			{
				"400 response from Cognito  getgroup error",
				nil,
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{},
					}, nil
				},
				func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, &models.ErrorResponse{}
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)

				},
			},
			{
				"500 response from Cognito  getgroup error",
				nil,
				func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
					var InternalException cognitoidentityprovider.InternalErrorException
					InternalException.Message_ = &internalErrorDescription
					return nil, &InternalException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{},
					}, nil
				},
				func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, &models.ErrorResponse{}
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)

				},
			},
		}

		for _, tt := range GetGroupTest {
			Convey(tt.description, func() {
				m.GetGroupFunc = tt.mock_getGroupfunc
				m.ListUsersInGroupFunc = tt.mock_listUsersInGroupfunc
				m.SetGroupUsersfunc = tt.mock_setGroupUsersfunc
				m.AddUserToGroupfunc = tt.mockAddUserToGroupFunction
				m.AdminAddUserToGroupFunc = tt.mock_addUserToGroupfunc
				m.RemoveUserFromGroupfunc = tt.mockRemoveUserToGroupFunction
				m.AdminRemoveUserFromGroupFunc = tt.mock_removeUserToGroupFunction
				postBody := tt.postbody
				body, err := json.Marshal(postBody)
				So(err, ShouldBeNil)
				r := httptest.NewRequest(http.MethodPut, addUserToGroupEndPoint, bytes.NewReader(body))
				urlVars := map[string]string{
					"id": "efgh5678",
				}
				r = mux.SetURLVars(r, urlVars)
				successResponse, errorResponse := api.SetGroupUsersHandler(ctx, w, r)
				tt.assertions(successResponse, errorResponse)
			})
		}
	})
}

func TestSetGroupUsers(t *testing.T) {
	api, _, m := apiSetup()

	Convey("Get group -check expected responses", t, func() {

		GetGroupTest := []struct {
			description                    string
			mockAddUserToGroupFunction     func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse)
			mock_addUserToGroupfunc        func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
			mockRemoveUserToGroupFunction  func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse)
			mock_removeUserToGroupFunction func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
			mock_listUsersInGroupfunc      func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			group                          models.Group
			users                          models.UsersList
			assertions                     func(successResponse *models.UsersList, errorResponse *models.ErrorResponse)
		}{
			{
				"200 response from Cognito  with input and output",
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{
							Users: []models.UserParams{{ID: "user_1"}},
							Count: 1,
						},
						nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				models.Group{
					ID: "test-group",
				},
				models.UsersList{
					Users: []models.UserParams{{ID: "user_1"}},
					Count: 1,
				},
				func(successResponse *models.UsersList, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
				},
			},
			{
				"404 response from Cognito ListUsers",
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{
							Users: []models.UserParams{{ID: "user_1"}},
							Count: 1,
						},
						nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				models.Group{
					ID: "test-group",
				},
				models.UsersList{
					Users: []models.UserParams{{ID: "user_1"}},
					Count: 1,
				},
				func(successResponse *models.UsersList, errorResponse *models.ErrorResponse) {
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.NotFoundError)
					So(castErr.Description, ShouldEqual, groupNotFoundDescription)
					So(successResponse, ShouldBeNil)
				},
			},
			{
				"500 response from Cognito listUsers ",
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{
							Users: []models.UserParams{{ID: "user_1"}},
							Count: 1,
						},
						nil
				},
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse) {
					return &models.UsersList{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := internalErrorDescription
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				models.Group{
					ID: "test-group",
				},
				models.UsersList{
					Users: []models.UserParams{{ID: "user_1"}},
					Count: 1,
				},
				func(successResponse *models.UsersList, errorResponse *models.ErrorResponse) {
					So(errorResponse.Status, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
					castErr := errorResponse.Errors[0].(*models.Error)
					So(castErr.Code, ShouldEqual, models.InternalError)
					So(castErr.Description, ShouldEqual, internalErrorDescription)
					So(successResponse, ShouldBeNil)
				},
			},
		}

		for _, tt := range GetGroupTest {
			Convey(tt.description, func() {
				m.AddUserToGroupfunc = tt.mockAddUserToGroupFunction
				m.AdminAddUserToGroupFunc = tt.mock_addUserToGroupfunc
				m.RemoveUserFromGroupfunc = tt.mockRemoveUserToGroupFunction
				m.AdminRemoveUserFromGroupFunc = tt.mock_removeUserToGroupFunction
				m.ListUsersInGroupFunc = tt.mock_listUsersInGroupfunc
				successResponse, errorResponse := api.SetGroupUsers(ctx, tt.group, tt.users)
				tt.assertions(successResponse, errorResponse)
			})
		}
	})
}

func TestRemoveUserFromGroup(t *testing.T) {
	var (
		userId = "abcd1234"
	)
	api, _, m := apiSetup()
	getGroupData := models.Group{
		ID: "123456789",
	}

	Convey("Remove a user from a group - check expected responses", t, func() {
		RemoveUsersTests := []struct {
			description                   string
			mockRemoveUserToGroupFunction func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
			mock_listUsersInGroupfunc     func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			assertions                    func(successResponse *models.UsersList, errorResponse error)
		}{
			{
				"200 response - user removed from group",
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
				},
			},
			{
				"Cognito 400 response - user not found",
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					var userNotFoundException cognitoidentityprovider.UserNotFoundException
					userNotFoundException.Message_ = &userNotFoundDescription
					return nil, &userNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.UserNotFoundException)
					So(*castErr.Message_, ShouldResemble, "user not found")
				},
			},
			{
				"Cognito 404 response - group not found",
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.ResourceNotFoundException)
					So(*castErr.Message_, ShouldResemble, "group not found")
				},
			},
			{
				"Cognito 500 response - internal error",
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.InternalErrorException)
					So(*castErr.Message_, ShouldResemble, "internal error")
				},
			},
			{
				"Cognito 404 response - listUsers group not found",
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.ResourceNotFoundException)
					So(*castErr.Message_, ShouldResemble, "group not found")
				},
			},
			{
				"Cognito 500 response - listUsers internal error",
				func(userInput *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
					return &cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var internalErrorException cognitoidentityprovider.InternalErrorException
					internalErrorException.Message_ = &internalErrorDescription
					return nil, &internalErrorException
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.InternalErrorException)
					So(*castErr.Message_, ShouldResemble, "internal error")
				},
			},
		}
		for _, tt := range RemoveUsersTests {
			Convey(tt.description, func() {
				m.AdminRemoveUserFromGroupFunc = tt.mockRemoveUserToGroupFunction
				m.ListUsersInGroupFunc = tt.mock_listUsersInGroupfunc
				successResponse, errorResponse := api.RemoveUserFromGroup(ctx, getGroupData, userId)
				tt.assertions(successResponse, errorResponse)
			})
		}
	})
}

func TestAddUserToGroup(t *testing.T) {

	var (
		userId = "abcd1234"
	)
	api, _, m := apiSetup()
	getGroupData := models.Group{
		ID: "123456789",
	}

	Convey("Remove a user from a group - check expected responses", t, func() {
		RemoveUsersTests := []struct {
			description               string
			mock_addUserToGroupfunc   func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
			mock_listUsersInGroupfunc func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			assertions                func(successResponse *models.UsersList, errorResponse error)
		}{
			{
				"200 response - user removed from group",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
				},
			},
			{
				"Cognito 400 response - user not found",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var userNotFoundException cognitoidentityprovider.UserNotFoundException
					userNotFoundException.Message_ = &userNotFoundDescription
					return nil, &userNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)

				},
			},
			{
				"Cognito 404 response - group not found",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.ResourceNotFoundException)
					So(*castErr.Message_, ShouldResemble, "group not found")
				},
			},
			{
				"Cognito 500 response - internal error",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					var internalError cognitoidentityprovider.InternalErrorException
					internalError.Message_ = &internalErrorDescription
					return nil, &internalError
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return &cognitoidentityprovider.ListUsersInGroupOutput{
						Users: []*cognitoidentityprovider.UserType{
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_3"),
							},
							{
								Enabled:    aws.Bool(true),
								UserStatus: aws.String("CONFIRMED"),
								Username:   aws.String("user_4"),
							},
						},
					}, nil
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.InternalErrorException)
					So(*castErr.Message_, ShouldResemble, "internal error")
				},
			},
			{
				"Cognito 404 response - listUsers group not found",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var groupNotFoundException cognitoidentityprovider.ResourceNotFoundException
					groupNotFoundException.Message_ = &groupNotFoundDescription
					return nil, &groupNotFoundException
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.ResourceNotFoundException)
					So(*castErr.Message_, ShouldResemble, "group not found")
				},
			},
			{
				"Cognito 500 response - listUsers internal error",
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
				},
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var internalErrorException cognitoidentityprovider.InternalErrorException
					internalErrorException.Message_ = &internalErrorDescription
					return nil, &internalErrorException
				},
				func(successResponse *models.UsersList, errorResponse error) {
					So(successResponse, ShouldBeNil)
					castErr := errorResponse.(*cognitoidentityprovider.InternalErrorException)
					So(*castErr.Message_, ShouldResemble, "internal error")
				},
			},
		}

		for _, tt := range RemoveUsersTests {
			Convey(tt.description, func() {
				m.AdminAddUserToGroupFunc = tt.mock_addUserToGroupfunc
				m.ListUsersInGroupFunc = tt.mock_listUsersInGroupfunc
				successResponse, errorResponse := api.AddUserToGroup(ctx, getGroupData, userId)
				tt.assertions(successResponse, errorResponse)
			})
		}
	})
}

func TestListGroupsUsersHandler(t *testing.T) {
	api, w, m := apiSetup()
	Convey("Check results for json", t, func() {
		listGroupsUsers := []struct {
			description          string
			listUsersInGroupFunc func(input *cognitoidentityprovider.ListUsersInGroupInput) (
				*cognitoidentityprovider.ListUsersInGroupOutput, error)
			listGroupsFunc func(input *cognitoidentityprovider.ListGroupsInput) (
				*cognitoidentityprovider.ListGroupsOutput, error)
			assertions func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				description: "empty group",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(0)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse.Headers, ShouldBeNil)
					So(isJson(successResponse, 0), ShouldBeTrue)
					So(isCSV(successResponse, 1), ShouldBeFalse)
				},
			},
			{
				description: "empty group no users",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(1)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse.Headers, ShouldBeNil)
					So(isJson(successResponse, 0), ShouldBeTrue)
					So(isCSV(successResponse, 1), ShouldBeFalse)
				},
			},
			{
				description: "json 1 group",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return listGroupsUsers(1), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(1)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse.Headers, ShouldBeNil)
					So(isJson(successResponse, 1), ShouldBeTrue)
					So(isCSV(successResponse, 1), ShouldBeFalse)
				},
			},
			{
				description: "json 3 group",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return listGroupsUsers(3), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(3)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse.Headers, ShouldBeNil)
					So(isJson(successResponse, 9), ShouldBeTrue)
					So(isCSV(successResponse, 1), ShouldBeFalse)
				},
			},
			{
				description: "json error getting groups",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					var exception cognitoidentityprovider.InternalErrorException
					exception.Message_ = &internalErrorDescription
					return nil, &exception
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(3)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
			{
				description: "json error getting group membership",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return listGroupsUsers(3), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					var exception cognitoidentityprovider.InternalErrorException
					exception.Message_ = &internalErrorDescription
					return nil, &exception
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse, ShouldNotBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
		}
		for _, tt := range listGroupsUsers {
			Convey(tt.description, func() {
				m.ListUsersInGroupFunc = tt.listUsersInGroupFunc
				m.ListGroupsFunc = tt.listGroupsFunc
				r := httptest.NewRequest(http.MethodGet, getGroupsReportEndPoint, nil)
				urlVars := map[string]string{
					"id": "efgh5678",
				}
				r = mux.SetURLVars(r, urlVars)
				successResponse, errorResponse := api.ListGroupsUsersHandler(ctx, w, r)
				tt.assertions(successResponse, errorResponse)
			})
		}
	})
	Convey("Check results for csv", t, func() {
		listGroupsUsers := []struct {
			description          string
			listUsersInGroupFunc func(input *cognitoidentityprovider.ListUsersInGroupInput) (
				*cognitoidentityprovider.ListUsersInGroupOutput, error)
			listGroupsFunc func(input *cognitoidentityprovider.ListGroupsInput) (
				*cognitoidentityprovider.ListGroupsOutput, error)
			assertions func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				description: "empty group",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(0)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse.Headers, ShouldNotBeNil)
					So(successResponse.Headers["Content-type"], ShouldEqual, "text/csv")
					So(isJson(successResponse, 0), ShouldBeFalse)
					So(isCSV(successResponse, 2), ShouldBeTrue)
				},
			},
			{
				description: "1 group no users",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(1)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse.Headers, ShouldNotBeNil)
					So(successResponse.Headers["Content-type"], ShouldEqual, "text/csv")
					So(isJson(successResponse, 0), ShouldBeFalse)
					So(isCSV(successResponse, 2), ShouldBeTrue)
				},
			},
			{
				description: "csv 1 group 1 user",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return listGroupsUsers(1), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(1)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse.Headers, ShouldNotBeNil)
					So(successResponse.Headers["Content-type"], ShouldEqual, "text/csv")
					So(isJson(successResponse, 0), ShouldBeFalse)
					So(isCSV(successResponse, 3), ShouldBeTrue)
				},
			},
			{
				description: "json 3 group",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return listGroupsUsers(3), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(3)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(successResponse.Headers, ShouldNotBeNil)
					So(successResponse.Headers["Content-type"], ShouldEqual, "text/csv")
					So(isJson(successResponse, 0), ShouldBeFalse)
					So(isCSV(successResponse, 11), ShouldBeTrue)
				},
			},
		}
		for _, tt := range listGroupsUsers {
			Convey(tt.description, func() {
				m.ListUsersInGroupFunc = tt.listUsersInGroupFunc
				m.ListGroupsFunc = tt.listGroupsFunc
				r := httptest.NewRequest(http.MethodGet, getGroupsReportEndPoint, nil)
				r.Header.Set("Accept", "text/csv")
				urlVars := map[string]string{
					"id": "efgh5678",
				}
				r = mux.SetURLVars(r, urlVars)
				successResponse, errorResponse := api.ListGroupsUsersHandler(ctx, w, r)
				tt.assertions(successResponse, errorResponse)
			})
		}
	})
	Convey("Check results for csv", t, func() {
		listGroupsUsers := []struct {
			description          string
			listUsersInGroupFunc func(input *cognitoidentityprovider.ListUsersInGroupInput) (
				*cognitoidentityprovider.ListUsersInGroupOutput, error)
			listGroupsFunc func(input *cognitoidentityprovider.ListGroupsInput) (
				*cognitoidentityprovider.ListGroupsOutput, error)
			assertions func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				description: "empty group",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(0)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(isJson(successResponse, 0), ShouldBeFalse)
					So(isCSV(successResponse, 2), ShouldBeTrue)
				},
			},
			{
				description: "1 group no users",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(1)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(isJson(successResponse, 0), ShouldBeFalse)
					So(isCSV(successResponse, 2), ShouldBeTrue)
				},
			},
			{
				description: "json 1 group 1 user",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return listGroupsUsers(1), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(1)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(isJson(successResponse, 0), ShouldBeFalse)
					So(isCSV(successResponse, 3), ShouldBeTrue)
				},
			},
			{
				description: "json 3 group",
				listUsersInGroupFunc: func(input *cognitoidentityprovider.ListUsersInGroupInput) (
					*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					return listGroupsUsers(3), nil
				},
				listGroupsFunc: func(input *cognitoidentityprovider.ListGroupsInput) (
					*cognitoidentityprovider.ListGroupsOutput, error) {
					output := listGroups(3)
					return &output, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(isJson(successResponse, 0), ShouldBeFalse)
					So(isCSV(successResponse, 11), ShouldBeTrue)
				},
			},
		}
		for _, tt := range listGroupsUsers {
			Convey(tt.description, func() {
				m.ListUsersInGroupFunc = tt.listUsersInGroupFunc
				m.ListGroupsFunc = tt.listGroupsFunc
				r := httptest.NewRequest(http.MethodGet, getGroupsReportEndPoint, nil)
				r.Header.Set("Accept", "text/csv")
				urlVars := map[string]string{
					"id": "",
				}
				r = mux.SetURLVars(r, urlVars)
				successResponse, errorResponse := api.ListGroupsUsersHandler(ctx, w, r)
				tt.assertions(successResponse, errorResponse)
			})
		}
	})
}

// isCSV will test that there is more than one slice and a header row
func isCSV(successResponse *models.SuccessResponse, expectedLength int) bool {
	testOutCSV := string(successResponse.Body[:])
	stringSlice := strings.Split(testOutCSV, "\n")
	if len(stringSlice) > 1 && stringSlice[0] == "Group,User" && len(stringSlice) == expectedLength {
		return true
	}
	return false
}

// isJson test that can be unmarshal into the given structure
func isJson(successResponse *models.SuccessResponse, expectedLength int) bool {
	var testOutJSON []models.ListGroupUsersType
	jsonErr := json.Unmarshal(successResponse.Body, &testOutJSON)
	if expectedLength > 0 {
		if jsonErr == nil && len(testOutJSON) == expectedLength {
			return true
		}
	} else {
		if jsonErr == nil && testOutJSON != nil {
			return true
		}
	}
	return false
}

func TestGetTeamsReportLines(t *testing.T) {
	api, _, m := apiSetup()
	Convey("init", t, func() {
		listGroupsUsers := []struct {
			description           string
			groupsList            cognitoidentityprovider.ListGroupsOutput
			listUsersForGroupFunc func(usersInput *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
			assertions            func(Response []models.ListGroupUsersType, errorResponse error)
		}{
			{
				"200 response - no groups",
				listGroups(0),
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l), nil
				},
				func(groupsUsersList []models.ListGroupUsersType, errorResponse error) {
					So(groupsUsersList, ShouldNotBeNil)
					So(errorResponse, ShouldBeNil)
				},
			},
			{
				"200 response - 1 groups",
				listGroups(1),
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l + 1), nil
				},
				func(groupsUsersList []models.ListGroupUsersType, errorResponse error) {
					So(errorResponse, ShouldBeNil)
					So(groupsUsersList, ShouldNotBeNil)
					So(groupsUsersList, ShouldHaveLength, 1)
					So(groupsUsersList[0].GroupName, ShouldResemble, "group 0 description")
					So(groupsUsersList[0].UserEmail, ShouldResemble, "user_0.email@domain.test")
				},
			},
			{
				"200 response - 3 groups",
				listGroups(3),
				func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
					l, _ := strconv.Atoi(string(*input.GroupName)[len(*input.GroupName)-1:])
					return listGroupsUsers(l + 1), nil
				},
				func(groupsUsersList []models.ListGroupUsersType, errorResponse error) {
					So(errorResponse, ShouldBeNil)
					So(groupsUsersList, ShouldNotBeNil)
					So(groupsUsersList, ShouldHaveLength, 6)
					So(groupsUsersList[0].GroupName, ShouldResemble, "group 0 description")
					So(groupsUsersList[0].UserEmail, ShouldResemble, "user_0.email@domain.test")
					So(groupsUsersList[2].GroupName, ShouldResemble, "group 1 description")
					So(groupsUsersList[2].UserEmail, ShouldResemble, "user_1.email@domain.test")
					So(groupsUsersList[4].GroupName, ShouldResemble, "group 2 description")
					So(groupsUsersList[4].UserEmail, ShouldResemble, "user_1.email@domain.test")
					So(groupsUsersList[5].GroupName, ShouldResemble, "group 2 description")
					So(groupsUsersList[5].UserEmail, ShouldResemble, "user_2.email@domain.test")
				},
			},
		}
		for _, tt := range listGroupsUsers {
			Convey(tt.description, func() {
				m.ListUsersInGroupFunc = tt.listUsersForGroupFunc
				groupMembershipList, errorResponse := api.GetTeamsReportLines(&tt.groupsList)
				tt.assertions(*groupMembershipList, errorResponse)
			})
		}
	})
}

// listGroups func to mock cognitoidentityprovider.ListGroupsOutput for use in TestGetTeamsReportLines
func listGroups(noOfGroups int) cognitoidentityprovider.ListGroupsOutput {
	groupList := []*cognitoidentityprovider.GroupType{}
	for i := 0; i < noOfGroups; i++ {
		groupName := fmt.Sprintf("group_%d", i)
		groupDescription := fmt.Sprintf("group %d description", i)
		groups := cognitoidentityprovider.GroupType{
			Description: &groupDescription,
			GroupName:   &groupName,
		}
		groupList = append(groupList, &groups)
	}
	output := cognitoidentityprovider.ListGroupsOutput{
		Groups:    groupList,
		NextToken: nil,
	}
	return output
}

// listGroupsUsers func to mock cognitoidentityprovider.ListUsersInGroupOutput for use in TestGetTeamsReportLines
func listGroupsUsers(noOfUsers int) *cognitoidentityprovider.ListUsersInGroupOutput {
	userList := []*cognitoidentityprovider.UserType{}
	var (
		attributeEmail = "email"
	)

	for i := 0; i < noOfUsers; i++ {
		userAttributes := []*cognitoidentityprovider.AttributeType{}
		userName := fmt.Sprintf("user_%d", i)
		userEmail := userName + ".email@domain.test"
		userAttribute := cognitoidentityprovider.AttributeType{Name: &attributeEmail, Value: &userEmail}
		userAttributes = append(userAttributes, &userAttribute)
		userType := cognitoidentityprovider.UserType{
			Enabled:    aws.Bool(true),
			UserStatus: aws.String("CONFIRMED"),
			Username:   aws.String(userName),
			Attributes: userAttributes,
		}
		userList = append(userList, &userType)
	}

	return &cognitoidentityprovider.ListUsersInGroupOutput{
		NextToken: nil,
		Users:     userList,
	}
}

// listGroupsUsersMock func to mock []models.ListGroupUsersType for use in TestListGroupsUsersHandler
func listGroupsUsersMock(noOfGroups, noOfUsers int) ([]models.ListGroupUsersType, error) {
	var output []models.ListGroupUsersType

	groupsList := listGroups(noOfGroups)
	for _, g := range groupsList.Groups {
		userList := listGroupsUsers(noOfUsers)
		for _, user := range userList.Users {
			for _, attribute := range user.Attributes {
				if strings.ToLower(*attribute.Name) == "email" {
					output = append(output, models.ListGroupUsersType{
						GroupName: *g.Description,
						UserEmail: *attribute.Value,
					})
				}
			}
		}

	}
	return output, nil
}

func TestSortGroups(t *testing.T) {
	Convey("Given a list of groups and a sort order", t, func() {

		groupA := "A Group"
		groupB := "B Group"
		groupC := "C Group"
		groups := []*cognitoidentityprovider.GroupType{
			{GroupName: &groupB},
			{GroupName: &groupC},
			{GroupName: &groupA},
		}
		Convey("When sorting by name in ascending order", func() {
			sortGroups(groups, "name:asc")
			Convey("The groups should be sorted in ascending order", func() {
				So(*groups[0].GroupName, ShouldEqual, "A Group")
				So(*groups[1].GroupName, ShouldEqual, "B Group")
				So(*groups[2].GroupName, ShouldEqual, "C Group")
			})
		})
		Convey("When sorting by name in descending order", func() {
			sortGroups(groups, "name:desc")
			Convey("The groups should be sorted in descending order", func() {
				So(*groups[0].GroupName, ShouldEqual, "C Group")
				So(*groups[1].GroupName, ShouldEqual, "B Group")
				So(*groups[2].GroupName, ShouldEqual, "A Group")
			})
		})
		Convey("When sorting with an invalid sortBy parameter", func() {
			sortGroups(groups, "non:sense")
			Convey("The groups should remain unsorted", func() {
				So(*groups[0].GroupName, ShouldEqual, "B Group")
				So(*groups[1].GroupName, ShouldEqual, "C Group")
				So(*groups[2].GroupName, ShouldEqual, "A Group")
			})
		})
		Convey("When sorting with an invalid asc or desc", func() {
			sortGroups(groups, "name:xyz")
			Convey("The groups should remain unsorted", func() {
				So(*groups[0].GroupName, ShouldEqual, "B Group")
				So(*groups[1].GroupName, ShouldEqual, "C Group")
				So(*groups[2].GroupName, ShouldEqual, "A Group")
			})
		})
		Convey("When sorting with an empty sort order", func() {
			sortGroups(groups, "")
			Convey("The groups should remain unsorted", func() {
				So(*groups[0].GroupName, ShouldEqual, "B Group")
				So(*groups[1].GroupName, ShouldEqual, "C Group")
				So(*groups[2].GroupName, ShouldEqual, "A Group")
			})
		})
	})
}
