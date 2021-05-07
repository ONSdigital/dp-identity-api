package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const usersEndPoint = "http://localhost:25600/users"

func TestCreateUserHandler(t *testing.T) {

	r := mux.NewRouter()
	ctx := context.Background()

	name := "Foo Bar"
	password := "temp1234"
	email := "foo_bar123@foobar.io.me"
	status := "UNCONFIRMED"

	m := &mock.MockCognitoIdentityProviderClient{}

	// mock call to: AdminCreateUser(input *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	m.AdminCreateUserFunc = func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
		user := &models.CreateUserOutput{
			UserOutput: &cognitoidentityprovider.AdminCreateUserOutput{
				User: &cognitoidentityprovider.UserType{
					Username:   &name,
					UserStatus: &status,
				},
			},
		}
		return user.UserOutput, nil
	}

	api := Setup(ctx, r, m, "us-west-11_bxushuds")

	Convey("Admin create user returns 201: successfully created user", t, func() {
		postBody := map[string]interface{}{"username": name, "password": password, "email": email}

		body, _ := json.Marshal(postBody)

		r := httptest.NewRequest("POST", usersEndPoint, bytes.NewReader(body))

		w := httptest.NewRecorder()

		api.Router.ServeHTTP(w, r)

		So(w.Code, ShouldEqual, http.StatusCreated)
	})

	Convey("Admin create user returns 500: error unmarshalling request body", t, func() {

		r := httptest.NewRequest("POST", usersEndPoint, bytes.NewReader(nil))

		w := httptest.NewRecorder()

		api.Router.ServeHTTP(w, r)

		body, _ := ioutil.ReadAll(w.Body)
		var e models.ErrorStructure
		json.Unmarshal(body, &e)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(e.Errors[0].Message, ShouldEqual, "api endpoint POST user returned an error unmarshalling request body")
	})

	Convey("Admin create user returns 400: bad request error", t, func() {

		r := httptest.NewRequest("POST", usersEndPoint, bytes.NewReader(nil))

		w := httptest.NewRecorder()

		api.Router.ServeHTTP(w, r)

		body, _ := ioutil.ReadAll(w.Body)
		var e models.ErrorStructure
		json.Unmarshal(body, &e)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(e.Errors[0].Message, ShouldEqual, "api endpoint POST user returned an error unmarshalling request body")
	})
}
