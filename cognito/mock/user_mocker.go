package mock

type User struct {
	ID         string
	Email      string
	Password   string
	GivenName  string
	FamilyName string
	Groups     []string
	Status     string
}

func (m *CognitoIdentityProviderClientStub) AddUserWithUsername(email, password string, isConfirmed bool) {
	m.Users = append(m.Users, m.GenerateUser(email, password, isConfirmed))
}

func (m *CognitoIdentityProviderClientStub) GenerateUser(email, password string, isConfirmed bool) User {
	statusString := "FORCE_CHANGE_PASSWORD"
	if isConfirmed {
		statusString = "CONFIRMED"
	}

	return User{
		ID:         "aaaabbbbcccc",
		Email:      email,
		Password:   password,
		GivenName:  "Bob",
		FamilyName: "Smith",
		Groups:     []string{},
		Status:     statusString,
	}
}
