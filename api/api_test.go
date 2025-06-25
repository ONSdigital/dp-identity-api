package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authorisation "github.com/ONSdigital/dp-authorisation/v2/authorisation/mock"
	"github.com/ONSdigital/dp-identity-api/v2/cognito/mock"
	jwksmock "github.com/ONSdigital/dp-identity-api/v2/jwks/mock"
	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/smithy-go"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	awsErrCode       = "InternalErrorException"
	awsErrMessage    = "Something strange happened"
	awsUNFErrCode    = "UserNotFoundException"
	awsUNFErrMessage = "user could not be found"
	unknownError     = smithy.ErrorFault(0)
	serverError      = smithy.ErrorFault(1)
	clientError      = smithy.ErrorFault(2)
)

var jwksHandler = jwksmock.JWKSStubbed

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()

		m := &mock.MockCognitoIdentityProviderClient{}
		m.CreateGroupFunc = func(_ context.Context, _ *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
			group := &cognitoidentityprovider.CreateGroupOutput{
				Group: &types.GroupType{},
			}
			return group, nil
		}

		api, err := Setup(ctx, r, m,
			"us-west-2_aaaaaaaaa", "client-aaa-bbb", "secret-ccc-ddd", "authflow", "eu-west-1234", true,
			[]string{"@ons.gov.uk", "@ext.ons.gov.uk"}, newAuthorisationMiddlwareMock(), jwksHandler)

		Convey("When created the following route(s) should have been added", func() {
			So(hasRoute(api.Router, "/v1/tokens", http.MethodPost), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/tokens/self", http.MethodDelete), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/tokens/self", http.MethodPut), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users", http.MethodPost), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users/{id}", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users/{id}/groups", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users/{id}", http.MethodPut), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/users/self/password", http.MethodPut), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/password-reset", http.MethodPost), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups-report", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups/{id}", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups/{id}", http.MethodDelete), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups/{id}", http.MethodPut), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups/{id}/members", http.MethodPost), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups/{id}/members", http.MethodGet), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups/{id}/members", http.MethodPut), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/groups/{id}/members/{user_id}", http.MethodDelete), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/jwt-keys", http.MethodGet), ShouldBeTrue)
		})

		Convey("No error returned when user pool id supplied", func() {
			So(err, ShouldBeNil)
		})

		Convey("Ensure cognito client has been added to api", func() {
			So(api.CognitoClient, ShouldNotBeNil)
		})
	})

	Convey("Given an API instance with an empty required parameter passed", t, func() {
		authorisationMiddleware := newAuthorisationMiddlwareMock()
		paramCheckTests := []struct {
			testName            string
			userPoolID          string
			clientID            string
			clientSecret        string
			clientAuthFlow      types.AuthFlowType
			awsRegion           string
			blockPlusAddressing bool
			allowedDomains      []string
		}{
			// missing userPoolID
			{
				"missing userPoolID",
				"",
				"client-aaa-bbb",
				"secret-ccc-ddd",
				"authflow",
				"eu-west-1234",
				true,
				[]string{"@ons.gov.uk", "@ext.ons.gov.uk"},
			},
			// missing clientID
			{
				"missing clientID",
				"eu-west-22_bdsjhids2",
				"",
				"secret-ccc-ddd",
				"authflow",
				"eu-west-1234",
				true,
				[]string{"@ons.gov.uk", "@ext.ons.gov.uk"},
			},
			// missing clientSecret
			{
				"missing clientSecret",
				"eu-west-22_bdsjhids2",
				"client-aaa-bbb",
				"",
				"authflow",
				"eu-west-1234",
				true,
				[]string{"@ons.gov.uk", "@ext.ons.gov.uk"},
			},
			// missing clientAuthFlow
			{
				"missing clientAuthFlow",
				"eu-west-22_bdsjhids2",
				"client-aaa-bbb",
				"secret-ccc-ddd",
				"",
				"eu-west-1234",
				true,
				[]string{"@ons.gov.uk", "@ext.ons.gov.uk"},
			},
			// missing allowedDomains
			{
				"missing allowedDomains",
				"eu-west-22_bdsjhids2",
				"client-aaa-bbb",
				"secret-ccc-ddd",
				"authflow",
				"eu-west-1234",
				true,
				nil,
			},
		}

		for _, tt := range paramCheckTests {
			r := mux.NewRouter()
			ctx := context.Background()
			_, err := Setup(ctx, r, &mock.MockCognitoIdentityProviderClient{}, tt.userPoolID, tt.clientID, tt.clientSecret, tt.awsRegion, tt.clientAuthFlow, tt.blockPlusAddressing, tt.allowedDomains, authorisationMiddleware, jwksHandler)

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
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

func apiMockSetup() (*API, *httptest.ResponseRecorder, *mock.MockCognitoIdentityProviderClient) {
	var (
		ctx                                       = context.Background()
		r                                         = mux.NewRouter()
		poolID, clientID, clientSecret, awsRegion = "us-west-11_bxushuds", "client-aaa-bbb", "secret-ccc-ddd", "eu-west-1234"
		authFlow                                  = types.AuthFlowTypeUserPasswordAuth
		blockPlusAddressing                       = true
		allowedDomains                            = []string{"@ons.gov.uk", "@ext.ons.gov.uk"}
	)

	m := &mock.MockCognitoIdentityProviderClient{}
	m.CreateGroupFunc = func(_ context.Context, _ *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
		group := &cognitoidentityprovider.CreateGroupOutput{
			Group: &types.GroupType{},
		}
		return group, nil
	}

	api, _ := Setup(ctx, r, m, poolID, clientID, clientSecret, awsRegion, authFlow, blockPlusAddressing, allowedDomains, newAuthorisationMiddlwareMock(), jwksHandler)

	w := httptest.NewRecorder()

	return api, w, m
}

func apiMockSetupWithDynamicBlockPlusAddressing(blockPlusAddressing bool) (*API, *httptest.ResponseRecorder, *mock.MockCognitoIdentityProviderClient) {
	var (
		ctx                                       = context.Background()
		r                                         = mux.NewRouter()
		poolID, clientID, clientSecret, awsRegion = "us-west-11_bxushuds", "client-aaa-bbb", "secret-ccc-ddd", "eu-west-1234"
		authFlow                                  = types.AuthFlowTypeUserPasswordAuth
		allowedDomains                            = []string{"@ons.gov.uk", "@ext.ons.gov.uk"}
	)

	m := &mock.MockCognitoIdentityProviderClient{}
	m.CreateGroupFunc = func(_ context.Context, _ *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
		group := &cognitoidentityprovider.CreateGroupOutput{
			Group: &types.GroupType{},
		}
		return group, nil
	}
	m.ListUsersFunc = func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
		// Return a mock response
		user := &cognitoidentityprovider.ListUsersOutput{
			Users: []types.UserType{},
		}
		return user, nil
	}

	api, _ := Setup(ctx, r, m, poolID, clientID, clientSecret, awsRegion, authFlow, blockPlusAddressing, allowedDomains, newAuthorisationMiddlwareMock(), jwksHandler)

	w := httptest.NewRecorder()

	return api, w, m
}

func TestWriteErrorResponse(t *testing.T) {
	Convey("the status code and the list of errors from the ErrorResponse object are written to a http response", t, func() {
		ctx := context.Background()

		errorResponseBodyExample := `{"errors":[{"code":"TestError","description":"a error generated for testing purposes"},{"code":"TestError","description":"another error generated for testing purposes"}]}`
		var errorResponse models.ErrorResponse

		errCode := "TestError"
		errDescription := "a error generated for testing purposes"
		anotherErrDescription := "another error generated for testing purposes"
		statusCode := http.StatusBadRequest

		headerMessage := "Test header message."

		errorResponse.Errors = append(errorResponse.Errors, models.NewValidationError(ctx, errCode, errDescription), models.NewValidationError(ctx, errCode, anotherErrDescription))
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

	Convey("the status code and InternalServerError as desc to http response for internal server errors", t, func() {
		ctx := context.Background()

		headerMsg := "Test header message."
		errorResponse := models.ErrorResponse{
			Errors: []error{models.NewValidationError(ctx, "TestError", "a error generated for testing purposes")},
			Status: http.StatusInternalServerError,
			Headers: map[string]string{
				WWWAuthenticateName: headerMsg,
			},
		}

		resp := httptest.NewRecorder()

		writeErrorResponse(ctx, resp, &errorResponse)

		So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		So(resp.Result().Header.Get(WWWAuthenticateName), ShouldEqual, headerMsg)
		So(resp.Body.String(), ShouldEqual, `{"code":"`+models.InternalError+`","description":"`+models.InternalErrorDescription+`"}`)
	})
}

func TestWriteSuccessResponse(t *testing.T) {
	Convey("test that authentication header data is successfully written in success response", t, func() {
		ctx := context.Background()
		body, err := json.Marshal(map[string]interface{}{"expirationTime": "12/12/2021T12:00:00Z", "refreshTokenExpirationTime": "13/12/2021T11:00:00Z"})
		So(err, ShouldBeNil)
		successResponseBodyExample := `{"expirationTime":"12/12/2021T12:00:00Z","refreshTokenExpirationTime":"13/12/2021T11:00:00Z"}`
		var (
			accessTokenHeaderMessage, idTokenHeaderMessage, refreshTokenHeaderMessage = "test-access-token-1", "test-id-token-1", "test-refresh-token-1"
		)
		successResponse := models.SuccessResponse{
			Body:   body,
			Status: http.StatusCreated,
			Headers: map[string]string{
				AccessTokenHeaderName:  accessTokenHeaderMessage,
				IDTokenHeaderName:      idTokenHeaderMessage,
				RefreshTokenHeaderName: refreshTokenHeaderMessage,
			},
		}

		resp := httptest.NewRecorder()

		writeSuccessResponse(ctx, resp, &successResponse)

		So(resp.Code, ShouldEqual, http.StatusCreated)
		So(resp.Body.String(), ShouldResemble, successResponseBodyExample)
		So(resp.Result().Header.Get(AccessTokenHeaderName), ShouldEqual, accessTokenHeaderMessage)
		So(resp.Result().Header.Get(IDTokenHeaderName), ShouldEqual, idTokenHeaderMessage)
		So(resp.Result().Header.Get(RefreshTokenHeaderName), ShouldEqual, refreshTokenHeaderMessage)
	})
}

func TestHandleBodyReadError(t *testing.T) {
	Convey("returns an ErrorResponse with a BodyReadError and 400 status", t, func() {
		ctx := context.Background()
		err := errors.New("TestError")
		errResponse := handleBodyReadError(ctx, err)

		So(errResponse.Status, ShouldEqual, http.StatusBadRequest)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.BodyReadError)
		So(castErr.Description, ShouldEqual, models.BodyReadFailedDescription)
	})
}

func TestHandleBodyUnmarshalError(t *testing.T) {
	Convey("returns an ErrorResponse with a JSONUnmarshalError and 400 status", t, func() {
		ctx := context.Background()
		err := errors.New("TestError")
		errResponse := handleBodyUnmarshalError(ctx, err)

		So(errResponse.Status, ShouldEqual, http.StatusBadRequest)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.JSONUnmarshalError)
		So(castErr.Description, ShouldEqual, models.ErrorUnmarshalFailedDescription)
	})
}

func TestInitialiseRoleGroups(t *testing.T) {
	Convey("Initialise role groups - check expected responses", t, func() {
		m := &mock.MockCognitoIdentityProviderClient{}

		ctx := context.Background()

		userPoolID := "us-west-11_bxushuds"

		adminCreateUsersTests := []struct {
			createGroupFunction func(ctx context.Context, input *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error)
			err                 error
		}{
			{
				// neither group exists
				func(_ context.Context, _ *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
					group := &cognitoidentityprovider.CreateGroupOutput{
						Group: &types.GroupType{},
					}
					return group, nil
				},
				nil,
			},
			{
				// admin group exists
				func(_ context.Context, input *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
					if *input.GroupName == models.AdminRoleGroup {
						awsErrCode := "GroupExistsException"
						awsErrMessage := "This group exists"
						awsErr := &smithy.GenericAPIError{
							Code:    awsErrCode,
							Message: awsErrMessage,
							Fault:   clientError,
						}
						return nil, awsErr
					}
					group := &cognitoidentityprovider.CreateGroupOutput{
						Group: &types.GroupType{},
					}
					return group, nil
				},
				nil,
			},
			{
				// publisher group exists
				func(_ context.Context, input *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
					if *input.GroupName == models.PublisherRoleGroup {
						awsErrCode := "GroupExistsException"
						awsErrMessage := "This group exists"
						awsErr := &smithy.GenericAPIError{
							Code:    awsErrCode,
							Message: awsErrMessage,
							Fault:   clientError,
						}
						return nil, awsErr
					}
					group := &cognitoidentityprovider.CreateGroupOutput{
						Group: &types.GroupType{},
					}
					return group, nil
				},
				nil,
			},
			{
				// create group internal error
				func(_ context.Context, _ *cognitoidentityprovider.CreateGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   serverError,
					}
					return nil, awsErr
				},
				models.NewError(ctx, nil, models.InternalError, "Something weird happened"),
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.CreateGroupFunc = tt.createGroupFunction

			err := initialiseRoleGroups(ctx, m, userPoolID)

			if tt.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
			}
		}
	})
}

func newAuthorisationMiddlwareMock() *authorisation.MiddlewareMock {
	return &authorisation.MiddlewareMock{
		RequireFunc: func(_ string, handlerFunc http.HandlerFunc) http.HandlerFunc {
			return handlerFunc
		},
	}
}
