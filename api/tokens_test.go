package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	. "github.com/smartystreets/goconvey/convey"
)

const signInEndPoint = "http://localhost:25600/tokens"
const signOutEndPoint = "http://localhost:25600/tokens/self"
const tokenRefreshEndPoint = "http://localhost:25600/tokens/self"

func TestDeprecatedWriteErrorResponse(t *testing.T) {
	Convey("A status code and an error body with two errors is written to a http response", t, func() {

		errorResponseBodyExample := `{"errors":[{"code":"Invalid email","description":"Unable to validate the email in the request"},{"code":"Invalid email","description":"Unable to validate the email in the request"}]}`

		var errorList []apierrorsdeprecated.Error

		errorList = nil

		errInvalidEmail := errors.New("Invalid email")
		invalidErrorDescription := "Unable to validate the email in the request"
		invalidEmailErrorBody := apierrorsdeprecated.IndividualErrorBuilder(errInvalidEmail, invalidErrorDescription)
		errorList = append(errorList, invalidEmailErrorBody)
		errorList = append(errorList, invalidEmailErrorBody)

		ctx := context.Background()
		resp := httptest.NewRecorder()
		statusCode := 400
		errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)

		apierrorsdeprecated.WriteErrorResponse(ctx, resp, statusCode, errorResponseBody)

		So(resp.Code, ShouldEqual, http.StatusBadRequest)
		So(resp.Body.String(), ShouldResemble, errorResponseBodyExample)
	})
}

func TestHandleUnexpectedError(t *testing.T) {
	Convey("An error and an error description is logged and written to a http response", t, func() {

		errorResponseBodyExample := `{"errors":[{"code":"unexpected error","description":"something unexpected has happened"}]}`

		ctx := context.Background()
		unexpectedError := errors.New("unexpected error")
		unexpectedErrorDescription := "something unexpected has happened"

		resp := httptest.NewRecorder()

		apierrorsdeprecated.HandleUnexpectedError(ctx, resp, unexpectedError, unexpectedErrorDescription)

		So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		So(resp.Body.String(), ShouldResemble, errorResponseBodyExample)
	})
}

func TestCognitoResponseHeaderBuild(t *testing.T) {
	Convey("build 201 response using an InitiateAuthOutput from Cognito", t, func() {
		w := httptest.NewRecorder()
		ctx := context.Background()
		accessToken := "accessToken"
		var expiration int64 = 123
		idToken := "idToken"
		Refresh := "refreshToken"

		initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    &expiration,
				IdToken:      &idToken,
				RefreshToken: &Refresh,
			},
		}

		buildSuccessfulResponse(initiateAuthOutput, w, ctx)

		So(w.Result().StatusCode, ShouldEqual, 201)
		So(w.Result().Header["Content-Type"], ShouldResemble, []string{"application/json"})
		So(w.Result().Header["Authorization"], ShouldResemble, []string{"Bearer " + accessToken})
		So(w.Result().Header["Id"], ShouldResemble, []string{idToken})
		So(w.Result().Header["Refresh"], ShouldResemble, []string{Refresh})

		var obj map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.String()), &obj)

		//there should be one entry in body
		So(len(obj), ShouldEqual, 1)

		type kv struct {
			Key   string
			Value interface{}
		}

		var ss []kv
		for k, v := range obj {
			ss = append(ss, kv{k, v})
		}
		str := fmt.Sprintf("%v", ss[0].Value)

		So(ss[0].Key, ShouldResemble, "expirationTime")
		So(str[:19], ShouldResemble, time.Now().UTC().Add(time.Second * 123).String()[:19])

	})

	Convey("build 500 response if the InitiateAuthOutput has an unexpected format", t, func() {
		w := httptest.NewRecorder()
		ctx := context.Background()

		initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{}
		buildSuccessfulResponse(initiateAuthOutput, w, ctx)

		So(w.Result().StatusCode, ShouldEqual, 500)
	})
}

func TestBuildJson(t *testing.T) {
	w := httptest.NewRecorder()

	Convey("build json", t, func() {
		w := httptest.NewRecorder()
		ctx := context.Background()

		testBody := map[string]interface{}{"expirationTime": "123"}
		buildjson(testBody, w, ctx)
		So(w.Body.String(), ShouldResemble, "{\"expirationTime\":\"123\"}")

	})

	Convey("build json err", t, func() {

		ctx := context.Background()

		testBody := map[string]interface{}{
			"foo": make(chan int),
		}
		buildjson(testBody, w, ctx)
		So(w.Body.String(), ShouldResemble, "{\"errors\":[{\"code\":\"json: unsupported type: chan int\",\"description\":\"failed to marshal the error\"}]}")
		So(w.Result().StatusCode, ShouldEqual, 500)
		So(w.Result().Header["Content-Type"], ShouldResemble, []string{"application/json"})
	})
}

func TestAPI_TokensHandler(t *testing.T) {
	var (
		r                                               = mux.NewRouter()
		ctx                                             = context.Background()
		poolId, clientId, clientSecret, authFlow string = "us-west-11_bxushuds", "client-aaa-bbb", "secret-ccc-ddd", "USER_PASSWORD_AUTH"
		accessToken, idToken, refreshToken       string = "aaaa.bbbb.cccc", "llll.mmmm.nnnn", "zzzz.yyyy.xxxx.wwww.vvvv"
		expireLength                             int64  = 500
	)

	m := &mock.MockCognitoIdentityProviderClient{}

	// mock call to: AdminUserGlobalSignOut(input *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
	m.AdminUserGlobalSignOutFunc = func(signOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
		return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
	}
	// mock call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	m.InitiateAuthFunc = func(signInInput *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
		return &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    &expireLength,
				IdToken:      &idToken,
				RefreshToken: &refreshToken,
			},
		}, nil
	}

	api, _ := Setup(ctx, r, m, poolId, clientId, clientSecret, authFlow)
	w := httptest.NewRecorder()

	Convey("Sign in success: no ErrorResponse, SuccessResponse Status 201", t, func() {
		body := map[string]interface{}{
			"email":    "email@ons.gov.uk",
			"password": "password",
		}
		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)
		request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

		successResponse, errorResponse := api.TokensHandler(w, request, ctx)

		So(errorResponse, ShouldBeNil)
		So(successResponse.Status, ShouldEqual, http.StatusCreated)
		var responseBody map[string]interface{}
		err = json.Unmarshal(successResponse.Body, &responseBody)
		So(err, ShouldBeNil)
		So(responseBody["expirationTime"], ShouldNotBeNil)
	})

	Convey("Sign In validation error: adds an error to the ErrorResponse and sets its Status to 400", t, func() {
		body := map[string]interface{}{
			"email":    "email@ons.gov.uk",
			"password": "",
		}
		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)
		request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

		successResponse, errorResponse := api.TokensHandler(w, request, ctx)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, models.InvalidPasswordError)
	})

	Convey("Sign In Cognito internal error: adds an error to the ErrorResponse and sets its Status to 500", t, func() {
		awsErrCode := "InternalErrorException"
		awsErrMessage := "Something strange happened"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
		m.InitiateAuthFunc = func(signInInput *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
			return nil, awsErr
		}

		body := map[string]interface{}{
			"email":    "email@ons.gov.uk",
			"password": "password",
		}
		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)
		request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

		successResponse, errorResponse := api.TokensHandler(w, request, ctx)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})

	Convey("Sign In Cognito request error: adds an error to the ErrorResponse and sets its Status to 400", t, func() {
		awsErrCode := "NotAuthorizedException"
		awsErrMessage := "User is not authorized"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
		m.InitiateAuthFunc = func(signInInput *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
			return nil, awsErr
		}

		body := map[string]interface{}{
			"email":    "email@ons.gov.uk",
			"password": "password",
		}
		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)
		request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

		successResponse, errorResponse := api.TokensHandler(w, request, ctx)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})
}

func TestAPI_SignOutHandler(t *testing.T) {
	var (
		r                                               = mux.NewRouter()
		ctx                                             = context.Background()
		poolId, clientId, clientSecret, authFlow string = "us-west-11_bxushuds", "client-aaa-bbb", "secret-ccc-ddd", "USER_PASSWORD_AUTH"
	)

	m := &mock.MockCognitoIdentityProviderClient{}

	// mock call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
	m.GlobalSignOutFunc = func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
		return &cognitoidentityprovider.GlobalSignOutOutput{}, nil
	}

	api, _ := Setup(ctx, r, m, poolId, clientId, clientSecret, authFlow)
	w := httptest.NewRecorder()

	Convey("Global Sign Out success: no errors added to ErrorResponse Errors list", t, func() {
		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, nil)
		request.Header.Set("Authorization", "Bearer zzzz-yyyy-xxxx")

		successResponse, errorResponse := api.SignOutHandler(w, request, ctx)

		So(errorResponse, ShouldBeNil)
		So(successResponse.Status, ShouldEqual, http.StatusNoContent)
		So(successResponse.Body, ShouldBeNil)
	})

	Convey("Global Sign Out validation error: adds an error to the ErrorResponse and sets its Status to 400", t, func() {
		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, nil)
		request.Header.Set("Authorization", "")

		successResponse, errorResponse := api.SignOutHandler(w, request, ctx)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, models.InvalidTokenError)
	})

	Convey("Global Sign Out Cognito internal error: adds an error to the ErrorResponse and sets its Status to 500", t, func() {
		awsErrCode := "InternalErrorException"
		awsErrMessage := "Something strange happened"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
		m.GlobalSignOutFunc = func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
			return nil, awsErr
		}

		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, nil)
		request.Header.Set("Authorization", "Bearer zzzz-yyyy-xxxx")

		successResponse, errorResponse := api.SignOutHandler(w, request, ctx)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})

	Convey("Global Sign Out Cognito request error: adds an error to the ErrorResponse and sets its Status to 400", t, func() {
		awsErrCode := "NotAuthorizedException"
		awsErrMessage := "User is not authorized"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
		m.GlobalSignOutFunc = func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
			return nil, awsErr
		}

		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, nil)
		request.Header.Set("Authorization", "Bearer zzzz-yyyy-xxxx")

		successResponse, errorResponse := api.SignOutHandler(w, request, ctx)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})
}
