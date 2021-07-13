package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const addUserToGroupEndPoint = "http://localhost:25600/v1/groups/efgh5678/memebers"

func TestAddUserToGroupHandler(t *testing.T) {

	var (
		ctx                                                                                = context.Background()
		userId                                                                      string = "abcd1234"
		userNotFoundDescription, groupNotFoundDescription, internalErrorDescription string = "user not found", "group not found", "internal error"
	)

	api, w, m := apiSetup()

	Convey("Add a user to a group - check expected responses", t, func() {
		adminCreateUsersTests := []struct {
			addUserToGroupFunction func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
			assertions             func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			// 200 response - user added to group
			{
				func(userInput *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
					return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
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

		for _, tt := range adminCreateUsersTests {
			m.AdminAddUserToGroupFunc = tt.addUserToGroupFunction

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
			// missing email
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
