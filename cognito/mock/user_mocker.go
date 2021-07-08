package mock

type User struct {
	ID          string
	Email       string
	Password    string
	GivenName   string
	FamilyName  string
	Groups      []string
	Status      string
	Active      bool
	StatusNotes string
}

func (m *CognitoIdentityProviderClientStub) AddUserWithEmail(email, password string, isConfirmed bool) {
	m.Users = append(m.Users, m.GenerateUser("", email, password, "", "", isConfirmed))
}

func (m *CognitoIdentityProviderClientStub) AddUserWithUsername(username, email string, isConfirmed bool) {
	m.Users = append(m.Users, m.GenerateUser(username, email, "", "", "", isConfirmed))
}

func (m *CognitoIdentityProviderClientStub) GenerateUser(id, email, password, givenName, familyName string, isConfirmed bool) *User {
	statusString := "FORCE_CHANGE_PASSWORD"
	if isConfirmed {
		statusString = "CONFIRMED"
	}
	if id == "" {
		id = "aaaabbbbcccc"
	}
	if email == "" {
		email = "email@ons.gov.uk"
	}
	if password == "" {
		password = "Passw0rd!"
	}
	if givenName == "" {
		givenName = "Bob"
	}
	if familyName == "" {
		familyName = "Smith"
	}

	return &User{
		ID:         id,
		Email:      email,
		Password:   password,
		GivenName:  givenName,
		FamilyName: familyName,
		Groups:     []string{},
		Status:     statusString,
		Active:     true,
	}
}

func (m *CognitoIdentityProviderClientStub) SetUserActiveState(username, active string) {
	for _, user := range m.Users {
		if user.ID == username {
			if active == "true" {
				user.Active = true
			} else {
				user.Active = true
			}
			return
		}
	}
}
