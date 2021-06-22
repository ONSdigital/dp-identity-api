package mock

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

type User struct {
	email    string
	password string
	Attributes []*cognitoidentityprovider.AttributeType
}

func (m *CognitoIdentityProviderClientStub) AddUserWithUsername(email, password string, addAttributes bool) {
	m.Users = append(m.Users, m.GenerateUser(email, password, addAttributes))
}

func (m *CognitoIdentityProviderClientStub) GenerateUser(email, password string, addAttributes bool) User {
	var user User
	user.email = email
	user.password = password
	// add email verified only if required
	if addAttributes {
		var (
			name, value string = "email_verified", "true"
		)
		user.Attributes = []*cognitoidentityprovider.AttributeType{
			{
				Name: &name,
				Value: &value,
			},
		}
	} else {
		user.Attributes = nil
	}
	return user
}
