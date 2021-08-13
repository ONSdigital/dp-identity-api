package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

//TokensHandler uses submitted email address and password to sign a user in against Cognito and returns a http handler interface
func (api *API) TokensHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	var userSignIn models.UserSignIn
	err = json.Unmarshal(body, &userSignIn)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	validationErrs := userSignIn.ValidateCredentials(ctx)
	if validationErrs != nil {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, *validationErrs...)
	}

	input := userSignIn.BuildCognitoRequest(api.ClientId, api.ClientSecret, api.ClientAuthFlow)
	result, authErr := api.CognitoClient.InitiateAuth(input)

	if authErr != nil {
		responseErr := models.NewCognitoError(ctx, authErr, "Cognito InitiateAuth request from sign in handler")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}

		switch responseErr.Description {
		case models.SignInFailedDescription:
			// Returning `WWW-Authenticate` in header as part of http.StatusUnauthorized response
			// See here: https://datatracker.ietf.org/doc/html/rfc7235#section-4.1
			headers := map[string]string{
				WWWAuthenticateName: "Bearer realm=\"" + ONSRealm + "\", charset=\"" + Charset + "\"",
			}
			return nil, models.NewErrorResponse(http.StatusUnauthorized, headers, responseErr)
		case models.SignInAttemptsExceededDescription:
			// Cognito returns the same Code for invalid credentials and too many attempts errors, changing our Error.Code to enable differentiation in the client
			responseErr.Code = models.TooManyFailedAttemptsError
			return nil, models.NewErrorResponse(http.StatusForbidden, nil, responseErr)
		default:
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		}
	}

	jsonResponse, responseErr := userSignIn.BuildSuccessfulJsonResponse(ctx, result)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	// success headers
	var headers map[string]string
	if result.AuthenticationResult != nil {
		//headers = map[string]string{
		//	AccessTokenHeaderName:  "Bearer " + *result.AuthenticationResult.AccessToken,
		//	IdTokenHeaderName:      *result.AuthenticationResult.IdToken,
		//	RefreshTokenHeaderName: *result.AuthenticationResult.RefreshToken,
		//}
		testRefreshToken := "osidhgowiehgoewirhjgoweirhgeriugherpiuwgherpigheriupghapurghrepiughawipurgberpiguharepiguherpiugheaiouhgrpeirtughaioywgfpiaeughrporgcniohgpxuhrgergnpixuhmgpirehncgipuhzmpiuwghxnicwgnxmzipwhgcoirhgimxzphweprgiuxenipzmhwgpinoernxpmizhwrgxicxrghxmipzhxnpifniozpruhmgpiwueghfoqerybmhipesxefgouyrgsidhmslgkxergxwefhmoirygnzoiwrengmxirhionzgliabfmriacnoxughzmo8ysgebmgrxylntvliaeoygxobmzlikriybnxillaiehfmgro8yiuaehdfmoy3gnoemzflknuhviklorisugheriuvnsdiygawiuefbeyiuagncbifawgnxiuawieufxteviufzabsviuetfvaiuwybvzniauwtefnuwgybnosidhgowiehgoewirhjgoweirhgeriugherpiuwgherpigheriupghapurghrepiughawipurgberpiguharepiguherpiugheaiouhgrpeirtughaioywgfpiaeughrporgcniohgpxuhrgergnpixuhmgpirehncgipuhzmpiuwghxnicwgnxmzipwhgcoirhgimxzphweprgiuxenipzmhwgpinoernxpmizhwrgxicxrghxmipzhxnpifniozpruhmgpiwueghfoqerybmhipesxefgouyrgsidhmslgkxergxwefhmoirygnzoiwrengmxirhionzgliabfmriacnoxughzmo8ysgebmgrxylntvliaeoygxobmzlikriybnxillaiehfmgro8yiuaehdfmoy3gnoemzflknuhviklorisugheriuvnsdiygawiuefbeyiuagncbifawgnxiuawieufxteviufzabsviuetfvaiuwybvzniauwtefnuwgybnosidhgowiehgoewirhjgoweirhgeriugherpiuwgherpigheriupghapurghrepiughawipurgberpiguharepiguherpiugheaiouhgrpeirtughaioywgfpiaeughrporgcniohgpxuhrgergnpixuhmgpirehncgipuhzmpiuwghxnicwgnxmzipwhgcoirhgimxzphweprgiuxenipzmhwgpinoernxpmizhwrgxicxrghxmipzhxnpifniozpruhmgpiwueghfoqerybmhipesxefgouyrgsidhmslgkxergxwefhmoirygnzoiwrengmxirhionzgliabfmriacnoxughzmo8ysgebmgrxylntvliaeoygxobmzlikriybnxillaiehfmgro8yiuaehdfmoy3gnoemzflknuhviklorisugheriuvnsdiygawiuefbeyiuagncbifawgnxiuawieufxteviufzabsviuetfvaiuwybvzniauwtefnuwgybnosidhgowiehgoewirhjgoweirhgeriugherpiuwgherpigheriupghapurghrepiughawipurgberpiguharepiguherpiugheaiouhgrpeirtughaioywgfpiaeughrporgcniohgpxuhrgergnpixuhmgpirehncgipuhzmpiuwghxnicwgnxmzipwhgcoirhgimxzphweprgiuxenipzmhwg"
		headers = map[string]string{
			AccessTokenHeaderName:  "Bearer eyJraWQiOiJBQzBnOXBzZzBwTEJ1Q2Nqa00yVkZEbXlzUlNxNm5KWlNxbkNXd1wvMFk1RT0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIwOTQ2MDIwNy1lMmZiLTQwNTQtOWVkOC03MDhlNjhjNDMxMmIiLCJjb2duaXRvOmdyb3VwcyI6WyJyb2xlLWFkbWluIl0sImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS13ZXN0LTEuYW1hem9uYXdzLmNvbVwvZXUtd2VzdC0xX0JuN0RhSXU3SCIsImNsaWVudF9pZCI6ImdoMGg4aTdja2N1OGJmMjFwOHIwb2pta2QiLCJvcmlnaW5fanRpIjoiOWEzYTA4M2UtMDc5Zi00MWFlLWJhMDMtZWYyZTZjOTI4MjgyIiwiZXZlbnRfaWQiOiJlNTY4NGNhYi01YzBhLTRkN2UtOWVmYy03NzM5MjJlZGI2NGEiLCJ0b2tlbl91c2UiOiJhY2Nlc3MiLCJzY29wZSI6ImF3cy5jb2duaXRvLnNpZ25pbi51c2VyLmFkbWluIiwiYXV0aF90aW1lIjoxNjI4ODY0NjA5LCJleHAiOjE2Mjg4NjgyMDksImlhdCI6MTYyODg2NDYwOSwianRpIjoiYWIzMGZjYjAtNjM1MC00NDVhLTg0YzktNmE4OTU3YWEzNDI4IiwidXNlcm5hbWUiOiJlODc3MDk4Ny1mYTQxLTRmNDktYWFlYi0zZWQ4ZjQwYzY2NmQifQ.e7psdzHXl2YR1zwH0-f1GgQjyywicym0NjWoUTU_WFhJFq-2K4mIDhnmto7-DSEG7Q_s5YUHjPus44zM7XT2UuclhaJd83ywW4FkPAwlvCgaV2Tc-JV7weilr57dEWN5l6Iz9L6kQAAM5uUXeU6vdTQezkbZ0PySX4cQ2k3fH0d1jeHg315NlRVMgpGtj6ApL2LDUbZOzkkcv9BrNCMyIZc9retrlECB5ctlHiiePGbw8yRq1p77suc06AWnhlWrzmAhPREDQb-Sa60my-zfZxxdx2JsFiC46Yausw1-E00aNvRpVIOczKa_Sg5-miCyYAjL4MmJN8HGd1Q84MO5tw",
			IdTokenHeaderName:      "eyJraWQiOiJHUkJldklyb0p6UEJ2YUdhTDl4bTR4XC82clFHa2JLeGkzd0x0Y1RpR3ltRT0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJjNjBjOTU4OS1lZGNiLTQzNjYtYTJlMC02YzVmOTU1NzU3YjMiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUm5tYTlscDJxIiwiY29nbml0bzp1c2VybmFtZSI6IjEyMjY0MGMxLTJmOGUtNDhkZC1hOGU4LWFhNDdhN2Q3NTcxNSIsImdpdmVuX25hbWUiOiJNYXR0Iiwib3JpZ2luX2p0aSI6ImFmOTAyOTUxLTVhMTgtNGYzZS1iZWY0LTNjNjRjMWQ4YTNhMyIsImF1ZCI6ImRmY200b25rNjBybXNzYThybjdzaGpraXYiLCJldmVudF9pZCI6IjAxMTEwNTlhLTg1NzktNDZkMi1iMTdlLThjMDg3MzU3MWFmNyIsInRva2VuX3VzZSI6ImlkIiwiYXV0aF90aW1lIjoxNjI4NjAxNzEzLCJleHAiOjE2Mjg2MDUzMTMsImlhdCI6MTYyODYwMTcxMywiZmFtaWx5X25hbWUiOiJOaWNrcyIsImp0aSI6IjFhNmVmN2EzLWNjMTUtNDU1NC1iNzk5LWIyYTY0OGRkYjc1YiIsImVtYWlsIjoibWF0dC5uaWNrc0Bob3RtYWlsLmNvLnVrIn0.FXnTKcFOhZrUqejjTXtkNjAittWbkLesHROOU2MOmbinRQqWZb__97wyBGfGes_qYIE9u-B5iQ1LgEFy1nmanESNCWwehjm4m4-Ms1F__BFNZI0M5fIAIeazhZycL3Tl9IbmQH-S95UsS3bl7NHAnFUnQzKlNopjCFcro4BeY-yDY2o3SsEXuUcJd6tBT13B471sQohyWRWhMxIm4tDCdQ6ieeR7Y9k0H7phfQHlKxxOqkQLQqfVXlHJsLuXHBRTqJ5eGbumGTuKYiDv9FtRQ7oBkiNZUY_qVqP3Syqp2_WdlotmQD-ndnlccxKag6koRKwrjJKD7rMSp_tIbTLKRQ",
			RefreshTokenHeaderName: testRefreshToken,
		}
	} else {
		headers = nil
	}

	// response - http.StatusCreated by default
	httpStatus := http.StatusCreated
	if result.ChallengeName != nil && *result.ChallengeName == NewPasswordChallenge {
		httpStatus = http.StatusAccepted
	}

	return models.NewSuccessResponse(jsonResponse, httpStatus, headers), nil
}

//SignOutHandler invalidates a users access token signing them out and returns a http handler interface
func (api *API) SignOutHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	accessToken := models.AccessToken{
		AuthHeader: req.Header.Get(AccessTokenHeaderName),
	}
	validationErr := accessToken.Validate(ctx)
	if validationErr != nil {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErr)
	}

	_, err := api.CognitoClient.GlobalSignOut(accessToken.GenerateSignOutRequest())

	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "Cognito GlobalSignOut request for sign out")
		if responseErr.Code == models.NotFoundError || responseErr.Code == models.NotAuthorisedError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		} else {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
	}

	return models.NewSuccessResponse(nil, http.StatusNoContent, nil), nil
}

//RefreshHandler refreshes a users access token and returns new access and ID tokens, expiration time and the refresh token
func (api *API) RefreshHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var validationErrs []error
	refreshToken := models.RefreshToken{TokenString: req.Header.Get(RefreshTokenHeaderName)}
	validationErr := refreshToken.Validate(ctx)
	if validationErr != nil {
		validationErrs = append(validationErrs, validationErr)
	}

	idToken := models.IdToken{TokenString: req.Header.Get(IdTokenHeaderName)}
	validationErr = idToken.Validate(ctx)
	if validationErr != nil {
		validationErrs = append(validationErrs, validationErr)
	}

	if len(validationErrs) > 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	authInput := refreshToken.GenerateRefreshRequest(api.ClientSecret, idToken.Claims.CognitoUser, api.ClientId)
	result, authErr := api.CognitoClient.InitiateAuth(authInput)

	if authErr != nil {
		responseErr := models.NewCognitoError(ctx, authErr, "Cognito InitiateAuth request for token refresh")
		if responseErr.Code == models.NotAuthorisedError {
			return nil, models.NewErrorResponse(http.StatusForbidden, nil, responseErr)
		} else {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
	}

	jsonResponse, responseErr := refreshToken.BuildSuccessfulJsonResponse(ctx, result)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	headers := map[string]string{
		AccessTokenHeaderName: "Bearer " + *result.AuthenticationResult.AccessToken,
		IdTokenHeaderName:     *result.AuthenticationResult.IdToken,
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusCreated, headers), nil
}

//SignOutAllUsersHandler bulk refresh token invalidation for panic sign out handling
func (api *API) SignOutAllUsersHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var (
		userFilterString string = "status=\"Enabled\""
		backOff                 = []time.Duration{
			1 * time.Second,
			3 * time.Second,
			10 * time.Second,
		}
	)
	usersList, awsErr := api.ListUsersWorker(req.Context(), &userFilterString, backOff)
	if awsErr != nil {
		return nil, awsErr
	}
	globalSignOut := &models.GlobalSignOut{
		ResultsChannel:  make(chan string, len(*usersList)),
		BackoffSchedule: backOff,
		RetryAllowed:    true,
	}
	// run api.SignOutUsersWorker concurrently
	go api.SignOutUsersWorker(req.Context(), globalSignOut, usersList)
	return models.NewSuccessResponse(nil, http.StatusAccepted, nil), nil
}

// ListUsersWorker - generates a list of users based on `userFilterString` filter string
func (api *API) ListUsersWorker(ctx context.Context, userFilterString *string, backoffSchedule []time.Duration) (*[]models.UserParams, *models.ErrorResponse) {
	var (
		awsErr                error
		usersList             models.UsersList
		listUsersResp, result *cognitoidentityprovider.ListUsersOutput
		listUserInput         = usersList.BuildListUserRequest(
			*userFilterString,
			"",
			int64(0),
			nil,
			&api.UserPoolId,
		)
		usersListError *models.ErrorResponse
	)
	listUsersResp, awsErr = api.generateListUsersRequest(listUserInput)
	if awsErr != nil {
		err := models.NewCognitoError(ctx, awsErr, "Cognito ListUsers request from signout all users from group endpoint")
		usersListError = models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	} else {
		if listUsersResp.PaginationToken != nil {
			listUserInput.PaginationToken = listUsersResp.PaginationToken
			// set `loadingInProgress` to control requesting new list data
			loadingInProgress := true
			for loadingInProgress {
				for _, backoff := range backoffSchedule {
					result, awsErr = api.generateListUsersRequest(listUserInput)
					if awsErr == nil {
						listUsersResp.Users = append(listUsersResp.Users, result.Users...)
						if result.PaginationToken != nil {
							listUserInput.PaginationToken = result.PaginationToken
							break
						} else {
							loadingInProgress = false
							break
						}
					} else {
						err := models.NewCognitoError(ctx, awsErr, "Cognito ListUsers request from signout all users from group endpoint")
						if err.Code != models.TooManyRequestsError {
							usersListError = models.NewErrorResponse(http.StatusInternalServerError, nil, err)
							loadingInProgress = false
							break
						}
					}
					time.Sleep(backoff)
				}
			}
		}
	}
	if usersListError != nil {
		return nil, usersListError
	} else {
		usersList.MapCognitoUsers(&listUsersResp.Users)
		return &usersList.Users, nil
	}
}

// SignOutUsersWorker - signs out users globally by invalidating user's refresh token
func (api *API) SignOutUsersWorker(ctx context.Context, g *models.GlobalSignOut, usersList *[]models.UserParams) {
	userSignOutRequestData := g.BuildSignOutUserRequest(usersList, &api.UserPoolId)

	for _, userSignoutRequest := range userSignOutRequestData {
		for _, backoff := range g.BackoffSchedule {
			_, err := api.generateGlobalSignOutRequest(userSignoutRequest)

			// no errors returned - add username to results channel and break to next user in list
			if err == nil {
				// the results channel is not being processed by caller currently - here for possible future extensibility
				g.ResultsChannel <- *userSignoutRequest.Username
				g.RetryAllowed = true

				break
			} else {
				// error returned - process it
				responseErr := models.NewCognitoError(ctx, err, "Cognito AdminUserGlobalSignOut request for sign out")

				// is error code != models.TooManyRequestsError? If so, proceed...
				if responseErr.Code != models.TooManyRequestsError {
					// if g.RetryAllowed is true, allowed to request api again
					if g.RetryAllowed {
						// attempt one more request to api
						g.RetryAllowed = false // 3. Set GlobalSignOut.RetryAllowed to false and request AdminUserGlobalSignOut
						_, retryErr := api.generateGlobalSignOutRequest(userSignoutRequest)

						if retryErr != nil {
							// if error response from request received again, process it
							retryResponseErr := models.NewCognitoError(ctx, err, "Cognito AdminUserGlobalSignOut request for sign out")

							// if error code != models.TooManyRequestsError break to next user
							if retryResponseErr.Code != models.TooManyRequestsError {
								g.RetryAllowed = true

								break
							}
						} else {
							// no error on retry, add user to results channel and break to next user
							// the results channel is not being processed by caller currently - here for possible future extensibility
							g.ResultsChannel <- *userSignoutRequest.Username
							g.RetryAllowed = true

							break
						}
					} else {
						// if GlobalSignOut.RetryAllowed already false break to next user
						g.RetryAllowed = true

						break
					}
				}
			}

			// backoff for predetermined length of time before requesting again
			time.Sleep(backoff)
		}

	}
	close(g.ResultsChannel)
}

// generateGlobalSignOutRequest - local routine to generete the global signout request per user
func (api *API) generateGlobalSignOutRequest(user *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
	return api.CognitoClient.AdminUserGlobalSignOut(user)
}

// generateListUsersRequest - local routine to generate a list users request
func (api *API) generateListUsersRequest(input *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
	return api.CognitoClient.ListUsers(input)
}
