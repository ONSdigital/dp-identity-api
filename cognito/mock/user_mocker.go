package mock

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	email 		string
	username 	string
	givenName 	string
	familyName 	string
	enabled 	bool
	created		time.Time
	status 		string
}

func (m *CognitoIdentityProviderClientStub) AddRandomUser() {
	m.Users = append(m.Users, m.GenerateUser("user@ons.gov.uk", uuid.NewV4().String(), "Jane","Doe"))
}

func (m *CognitoIdentityProviderClientStub) AddUserWithUsername(username string) {
	m.Users = append(m.Users, m.GenerateUser("user@ons.gov.uk", username, "Jane","Doe"))
}

func (m *CognitoIdentityProviderClientStub) GenerateUser(email string, username string, givenName string, familyName string) User {
	layout := "2006-01-02T15:04:05.000Z"
	createdTime, _ := time.Parse(layout, "2021-01-01T12:00:00.000Z")
	return User{
		email: email,
		username: username,
		givenName: givenName,
		familyName: familyName,
		enabled: true,
		created: createdTime,
		status: "CONFIRMED",
	}
}
