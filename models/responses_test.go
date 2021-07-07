package models_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-identity-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewErrorResponse(t *testing.T) {
	Convey("successfully constructs an ErrorResponse object", t, func() {

		Convey("with one error", func() {
			err := models.Error{
				Cause:       errors.New("TestError"),
				Code:        "TestErrorCode",
				Description: "description of the error",
			}

			errResponse := models.NewErrorResponse(http.StatusInternalServerError, nil, &err)

			So(errResponse, ShouldNotBeNil)
			So(len(errResponse.Errors), ShouldEqual, 1)
			So(errResponse.Status, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("with multiple errors", func() {
			err := models.Error{
				Cause:       errors.New("TestError"),
				Code:        "TestErrorCode",
				Description: "description of the error",
			}
			err2 := models.Error{
				Cause:       errors.New("TestError2"),
				Code:        "TestErrorCode",
				Description: "description of the error",
			}
			errList := []error{
				&err,
				&err2,
			}

			errResponse := models.NewErrorResponse(http.StatusInternalServerError, nil, errList...)

			So(errResponse, ShouldNotBeNil)
			So(len(errResponse.Errors), ShouldEqual, 2)
			So(errResponse.Status, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestNewSuccessResponse(t *testing.T) {
	Convey("successfully constructs a SuccessResponse object", t, func() {
		postBody := map[string]interface{}{"expirationTime": "tomorrow"}
		jsonResponse, err := json.Marshal(postBody)
		So(err, ShouldBeNil)

		successResponse := models.NewSuccessResponse(jsonResponse, http.StatusCreated, nil)

		So(successResponse, ShouldNotBeNil)
		So(successResponse.Body, ShouldNotBeNil)
		So(successResponse.Status, ShouldEqual, http.StatusCreated)
	})
}
