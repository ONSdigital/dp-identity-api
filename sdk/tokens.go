package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/models"
	apiError "github.com/ONSdigital/dp-identity-api/sdk/errors"
)

type TokenResponse struct {
	Token                      string `json:"-"`
	RefreshToken               string `json:"-"`
	ExpirationTime             string `json:"expirationTime"`
	RefreshTokenExpirationTime string `json:"refreshTokenExpirationTime"`
}

// GetToken attempts to sign in and obtain a JWT token from the API
func (cli *Client) GetToken(ctx context.Context, credentials models.UserSignIn) (*TokenResponse, apiError.Error) {
	path := fmt.Sprintf("%s/tokens", cli.hcCli.URL)

	b, _ := json.Marshal(credentials)

	respInfo, apiErr := cli.callIdentityAPI(ctx, path, http.MethodPost, b)
	if apiErr != nil {
		return nil, apiErr
	}

	var tokenResponse TokenResponse

	if err := json.Unmarshal(respInfo.Body, &tokenResponse); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal tokenResponse - error is: %v", err),
		}
	}

	var headers = respInfo.Headers

	tokenResponse.Token = headers.Get("Authorization")
	tokenResponse.RefreshToken = headers.Get("Refresh")

	return &tokenResponse, nil
}
