package mock

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
)

type User struct {
	ID          string
	Email       string
	Password    string
	GivenName   string
	FamilyName  string
	Groups      []*Group
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

//Generates the required number of users in the system
func (m *CognitoIdentityProviderClientStub) AddMultipleUsers(usersCount int) {
	for len(m.Users) < usersCount {
		m.Users = append(m.Users, m.GenerateUser("", "", "", "", "", true))
	}
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
		Groups:     []*Group{},
		Status:     statusString,
		Active:     true,
	}
}

func (m *CognitoIdentityProviderClientStub) SetUserActiveState(username, active string) {
	for _, user := range m.Users {
		if user.ID == username {
			user.Active = "true" == active
			return
		}
	}
}

func (m *CognitoIdentityProviderClientStub) ReadUser(username string) *User {
	for _, user := range m.Users {
		if user.ID == username {
			return user
		}
	}
	return nil
}

func (m *CognitoIdentityProviderClientStub) MakeUserMember(userName string) {

	user := m.ReadUser(userName)
	if user != nil {
		for _, group := range m.Groups {
			if !strings.HasPrefix(group.Name, "role-") {
				m.AddUserToGroup(user.ID, group.Name)
			}
		}
	}
}

//BulkGenerateUsers - bulk generate 'n' users for testing purposes
//                    if usernames array is nil or length is different, will auto-assign UUIDs
func BulkGenerateUsers(userCount int, userNames []string) *cognitoidentityprovider.ListUsersOutput {
	paginationToken := "abc-123-xyz-345-xxx"
	usersList := &cognitoidentityprovider.ListUsersOutput{}
	for i := 0; i < userCount; i++ {
		var (
			userId, status string = "", "CONFIRMED"
			enabled        bool   = true
		)
		if userNames == nil || i > len(userNames)-1 {
			userId = uuid.NewString()
		} else {
			userId = userNames[i]
		}
		user := &cognitoidentityprovider.UserType{}
		user.Username = &userId
		user.Enabled = &enabled
		user.UserStatus = &status
		usersList.Users = append(usersList.Users, user)
		usersList.PaginationToken = &paginationToken
	}
	return usersList
}
