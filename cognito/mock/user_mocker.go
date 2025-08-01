package mock

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/google/uuid"
)

type User struct {
	ID          string
	Email       string
	Password    string
	GivenName   string
	FamilyName  string
	Groups      []*Group
	Status      types.UserStatusType
	Active      bool
	StatusNotes string
}

func (m *CognitoIdentityProviderClientStub) AddUserWithEmail(email, password string, isConfirmed bool) {
	m.Users = append(m.Users, m.GenerateUser("", email, password, "", "", isConfirmed))
}

func (m *CognitoIdentityProviderClientStub) AddUserWithUsername(username, email string, isConfirmed bool) {
	m.Users = append(m.Users, m.GenerateUser(username, email, "", "", "", isConfirmed))
}
func (m *CognitoIdentityProviderClientStub) AddUserWithAttributes(id, forename, lastname, email, password string, isConfirmed bool) {
	m.Users = append(m.Users, m.GenerateUser(id, email, password, forename, lastname, isConfirmed))
}

func (m *CognitoIdentityProviderClientStub) AddUserWithForename(username, email, forename string, isConfirmed bool) {
	m.Users = append(m.Users, m.GenerateUser(username, email, "", forename, "", isConfirmed))
}

// AddMultipleUsers generates the required number of users in the system
func (m *CognitoIdentityProviderClientStub) AddMultipleUsers(usersCount int) {
	for len(m.Users) < usersCount {
		m.Users = append(m.Users, m.GenerateUser("", "", "", "", "", true))
	}
}

// AddMultipleActiveUsers generates the required number of users in the system
func (m *CognitoIdentityProviderClientStub) AddMultipleActiveUsers(activeusersCount, inactiveusersCount int) {
	for len(m.Users) < activeusersCount {
		m.Users = append(m.Users, m.GenerateUser("", "", "", "", "", true))
	}
	for len(m.Users) < inactiveusersCount {
		m.Users = append(m.Users, m.GenerateUser("", "", "", "", "", false))
	}
}

func (m *CognitoIdentityProviderClientStub) GenerateUser(id, email, password, givenName, familyName string, isConfirmed bool) *User {
	statusString := types.UserStatusTypeForceChangePassword
	if isConfirmed {
		statusString = types.UserStatusTypeConfirmed
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
			user.Active = active == "true"
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
	ctx := context.Background()
	user := m.ReadUser(userName)
	if user != nil {
		for _, group := range m.Groups {
			if !strings.HasPrefix(group.Name, "role-") {
				if err := m.AddUserToGroup(user.ID, group.Name); err != nil {
					log.Warn(ctx, "error adding user to group", log.Data{
						"userName":  userName,
						"groupName": group.Name,
					})
				}
			}
		}
	}
}

// BulkGenerateUsers - bulk generate 'n' users for testing purposes
//
//	if usernames array is nil or length is different, will auto-assign UUIDs
func BulkGenerateUsers(userCount int, userNames []string) *cognitoidentityprovider.ListUsersOutput {
	paginationToken := "abc-123-xyz-345-xxx"
	usersList := &cognitoidentityprovider.ListUsersOutput{}
	for i := 0; i < userCount; i++ {
		var (
			userID = ""
			status = types.UserStatusTypeConfirmed
		)
		if userNames == nil || i > len(userNames)-1 {
			userID = uuid.NewString()
		} else {
			userID = userNames[i]
		}
		user := types.UserType{}
		user.Username = &userID
		user.Enabled = true
		user.UserStatus = status
		usersList.Users = append(usersList.Users, user)
		usersList.PaginationToken = &paginationToken
	}
	return usersList
}
