package mock

type User struct {
	email    string
	password string
}

func (m *CognitoIdentityProviderClientStub) AddUserWithUsername(email, password string) {
	m.Users = append(m.Users, m.GenerateUser(email, password))
}

func (m *CognitoIdentityProviderClientStub) GenerateUser(email, password string) User {
	return User{
		email:    email,
		password: password,
	}
}
