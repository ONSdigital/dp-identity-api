package models_test

import (
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-identity-api/models"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestNewErrorResponse(t *testing.T) {
	Convey("successfully constructs an ErrorResponse object", t, func() {
		err := models.Error{
			Cause:       errors.New("TestError"),
			Code:        "TestErrorCode",
			Description: "description of the error",
		}

		errResponse := models.NewErrorResponse([]error{&err}, http.StatusInternalServerError)

		So(errResponse, ShouldNotBeNil)
		So(len(errResponse.Errors), ShouldEqual, 1)
		So(errResponse.Status, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestNewSuccessResponse(t *testing.T) {
	Convey("successfully constructs a SuccessResponse object", t, func() {
		postBody := map[string]interface{}{"expirationTime": "tomorrow"}
		jsonResponse, err := json.Marshal(postBody)
		So(err, ShouldBeNil)

		successResponse := models.NewSuccessResponse(jsonResponse, http.StatusCreated)

		So(successResponse, ShouldNotBeNil)
		So(successResponse.Body, ShouldNotBeNil)
		So(successResponse.Status, ShouldEqual, http.StatusCreated)
	})
}
