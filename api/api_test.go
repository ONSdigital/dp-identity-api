package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/models"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-identity-api/cognito/mock"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		api, err := Setup(ctx, r, &mock.MockCognitoIdentityProviderClient{}, "us-west-2_aaaaaaaaa", "client-aaa-bbb", "secret-ccc-ddd", "authflow")

		Convey("When created the following route(s) should have been added", func() {
			So(hasRoute(api.Router, "/v1/tokens", "POST"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/tokens/self", "DELETE"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/tokens/self", "PUT"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users", "POST"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users/self/password", "PUT"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/password-reset", "POST"), ShouldBeTrue)
		})

		Convey("No error returned when user pool id supplied", func() {
			So(err, ShouldBeNil)
		})

		Convey("Ensure cognito client has been added to api", func() {
			So(api.CognitoClient, ShouldNotBeNil)
		})
	})

	Convey("Given an API instance with an empty required parameter passed", t, func() {
		paramCheckTests := []struct {
			testName       string
			userPoolId     string
			clientId       string
			clientSecret   string
			clientAuthFlow string
		}{
			// missing userPoolId
			{
				"missing userPoolId",
				"",
				"client-aaa-bbb",
				"secret-ccc-ddd",
				"authflow",
			},
			// missing clientId
			{
				"missing clientId",
				"eu-west-22_bdsjhids2",
				"",
				"secret-ccc-ddd",
				"authflow",
			},
			// missing clientSecret
			{
				"missing clientSecret",
				"eu-west-22_bdsjhids2",
				"client-aaa-bbb",
				"",
				"authflow",
			},
			// missing clientAuthFlow
			{
				"missing clientAuthFlow",
				"eu-west-22_bdsjhids2",
				"client-aaa-bbb",
				"secret-ccc-ddd",
				"",
			},
		}

		for _, tt := range paramCheckTests {
			r := mux.NewRouter()
			ctx := context.Background()
			_, err := Setup(ctx, r, &mock.MockCognitoIdentityProviderClient{}, tt.userPoolId, tt.clientId, tt.clientSecret, tt.clientAuthFlow)

			Convey("Error should not be nil if require parameter is empty: "+tt.testName, func() {
				So(err.Error(), ShouldEqual, models.MissingConfigError+": "+models.MissingConfigDescription)
				castErr := err.(*models.Error)
				So(castErr.Code, ShouldEqual, models.MissingConfigError)
				So(castErr.Description, ShouldEqual, models.MissingConfigDescription)
			})
		}
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

func TestWriteErrorResponse(t *testing.T) {
	Convey("the status code and the list of errors from the ErrorResponse object are written to a http response", t, func() {
		ctx := context.Background()

		errorResponseBodyExample := `{"errors":[{"code":"TestError","description":"a error generated for testing purposes"},{"code":"TestError","description":"a error generated for testing purposes"}]}`
		var errorResponse models.ErrorResponse

		errCode := "TestError"
		errDescription := "a error generated for testing purposes"
		statusCode := http.StatusBadRequest

		headerMessage := "Test header message."

		errorResponse.Errors = append(errorResponse.Errors, models.NewValidationError(ctx, errCode, errDescription))
		errorResponse.Errors = append(errorResponse.Errors, models.NewValidationError(ctx, errCode, errDescription))
		errorResponse.Status = statusCode
		errorResponse.Headers = map[string]string{
			WWWAuthenticateName: headerMessage,
		}

		resp := httptest.NewRecorder()

		writeErrorResponse(ctx, resp, &errorResponse)

		So(resp.Code, ShouldEqual, http.StatusBadRequest)
		So(resp.Body.String(), ShouldResemble, errorResponseBodyExample)
		So(resp.Result().Header.Get(WWWAuthenticateName), ShouldEqual, headerMessage)
	})
}

func TestWriteSuccessResponse(t *testing.T) {
	Convey("test that authentication header data is successfully written in success response", t, func() {
		ctx := context.Background()
		body, err := json.Marshal(map[string]interface{}{"expirationTime": "12/12/2021T12:00:00Z"})
		So(err, ShouldBeNil)
		successResponseBodyExample := `{"expirationTime":"12/12/2021T12:00:00Z"}`
		var (
			accessTokenHeaderMessage, idTokenHeaderMessage, refreshTokenHeaderMessage string = "test-access-token-1", "test-id-token-1", "test-refresh-token-1"
		)
		successResponse := models.SuccessResponse{
			Body:   body,
			Status: http.StatusCreated,
			Headers: map[string]string{
				AccessTokenHeaderName:  accessTokenHeaderMessage,
				IdTokenHeaderName:      idTokenHeaderMessage,
				RefreshTokenHeaderName: refreshTokenHeaderMessage,
			},
		}

		resp := httptest.NewRecorder()

		writeSuccessResponse(ctx, resp, &successResponse)

		So(resp.Code, ShouldEqual, http.StatusCreated)
		So(resp.Body.String(), ShouldResemble, successResponseBodyExample)
		So(resp.Result().Header.Get(AccessTokenHeaderName), ShouldEqual, accessTokenHeaderMessage)
		So(resp.Result().Header.Get(IdTokenHeaderName), ShouldEqual, idTokenHeaderMessage)
		So(resp.Result().Header.Get(RefreshTokenHeaderName), ShouldEqual, refreshTokenHeaderMessage)
	})
}

func TestHandleBodyReadError(t *testing.T) {
	Convey("returns an ErrorResponse with a BodyReadError and 500 status", t, func() {
		ctx := context.Background()
		err := errors.New("TestError")
		errResponse := handleBodyReadError(ctx, err)

		So(errResponse.Status, ShouldEqual, http.StatusInternalServerError)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.BodyReadError)
		So(castErr.Description, ShouldEqual, models.BodyReadFailedDescription)
	})
}

func TestHandleBodyUnmarshalError(t *testing.T) {
	Convey("returns an ErrorResponse with a JSONUnmarshalError and 500 status", t, func() {
		ctx := context.Background()
		err := errors.New("TestError")
		errResponse := handleBodyUnmarshalError(ctx, err)

		So(errResponse.Status, ShouldEqual, http.StatusInternalServerError)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.JSONUnmarshalError)
		So(castErr.Description, ShouldEqual, models.ErrorUnmarshalFailedDescription)
	})
}

func TestInitialiseRoleGroups(t *testing.T) {
	Convey("Initialise role groups - check expected responses", t, func() {
		m := &mock.MockCognitoIdentityProviderClient{}

		ctx := context.Background()

		userPoolId := "us-west-11_bxushuds"

		adminCreateUsersTests := []struct {
			createGroupFunction func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error)
			err                 error
		}{
			{
				// neither group exists
				func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error) {
					group := &cognitoidentityprovider.CreateGroupOutput{
						Group: &cognitoidentityprovider.GroupType{},
					}
					return group, nil
				},
				nil,
			},
			{
				// admin group exists
				func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error) {
					if *input.GroupName == models.AdminRoleGroup {
						awsErrCode := "GroupExistsException"
						awsErrMessage := "This group exists"
						awsOrigErr := errors.New(awsErrCode)
						awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
						return nil, awsErr
					} else {
						group := &cognitoidentityprovider.CreateGroupOutput{
							Group: &cognitoidentityprovider.GroupType{},
						}
						return group, nil
					}
				},
				nil,
			},
			{
				// publisher group exists
				func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error) {
					if *input.GroupName == models.PublisherRoleGroup {
						awsErrCode := "GroupExistsException"
						awsErrMessage := "This group exists"
						awsOrigErr := errors.New(awsErrCode)
						awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
						return nil, awsErr
					} else {
						group := &cognitoidentityprovider.CreateGroupOutput{
							Group: &cognitoidentityprovider.GroupType{},
						}
						return group, nil
					}
				},
				nil,
			},
			{
				// create group internal error
				func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something weird happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				models.NewError(ctx, nil, models.InternalError, "Something weird happened"),
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.CreateGroupFunc = tt.createGroupFunction

			err := initialiseRoleGroups(ctx, m, userPoolId)

			if tt.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
			}
		}
	})
}
