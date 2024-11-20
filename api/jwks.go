package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/v2/models"
)

// CognitoPoolJWKSHandler handles the retrieval of pool specific web key set
func (api *API) CognitoPoolJWKSHandler(ctx context.Context, _ http.ResponseWriter, _ *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	keyData, err := api.JWKSManager.JWKSGetKeyset(api.AWSRegion, api.UserPoolID)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusNotFound, nil, err)
	}

	jsonResponse, err := api.JWKSManager.JWKSToRSAJSONResponse(keyData)
	if err != nil {
		return nil, handleJWKSParsingErrors(ctx, err)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

func handleJWKSParsingErrors(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.JWKSParseError, models.JWKSParseErrorDescription),
	)
}
