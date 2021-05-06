package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateUserHandler(t *testing.T) {

	r := mux.NewRouter()
	ctx := context.Background()

	name := "Foo Bar"
	password := "temp1234"
	email := "foo_bar123@foobar.io.me"

	m := &mock.MockCognitoIdentityProviderClient{}

	m.AdminCreateUserFunc = func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {

		status := "UNCONFIRMED"

		v := &models.CreateUserOutput{
			&cognitoidentityprovider.AdminCreateUserOutput{
				User: &cognitoidentityprovider.UserType{
					Username:   &name,
					UserStatus: &status,
				},
			},
		}

		return v.AdminCreateUserOutput, nil
	}

	api := Setup(ctx, r, m, "us-west-11_bxushuds")

	Convey("Admin create user returns 201: success", t, func() {
		postBody := map[string]interface{}{"username": name, "password": password, "email": email}

		body, _ := json.Marshal(postBody)

		req := httptest.NewRequest("POST", "localhost:25600/users", bytes.NewReader(body))

		createUserHandler := api.CreateUserHandler(ctx)
		createUserHandler.ServeHTTP(httptest.NewRecorder(), req)
		res := req.Response

		// this will break until implemented!
		So(res.StatusCode, ShouldEqual, http.StatusCreated)
	})
}
