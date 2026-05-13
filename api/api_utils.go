package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/ONSdigital/log.go/v2/log"
)

func handleAuthEntityDataError(ctx context.Context, err error, logData log.Data) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.GetAuthEntityDataError, models.GetAuthEntityDataErrorDescription, logData),
	)
}
