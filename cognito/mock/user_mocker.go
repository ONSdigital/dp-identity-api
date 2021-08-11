package mock

import (
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> d47026707c3c9cc5b8217883e34fb1da111b238b
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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

//BulkGenerateGroups - bulk generate 'n' groups for testing purposes
//                    if groupnames array is nil or length is different, will autofill with
func BulkGenerateGroups(groupCount int, groupNames []string) *cognitoidentityprovider.AdminListGroupsForUserOutput {
	nextToken := "abc-123-xyz-345-xxx"
	groupList := &cognitoidentityprovider.AdminListGroupsForUserOutput{}

	for i := 0; i < groupCount; i++ {
		var (
			timestamp         = time.Now()
			randomNum   int64 = int64(rand.Intn((100 - 3) + 3))
			userPoolId        = "aaaa-bbbb-ccc-dddd"
			group_name        = "group_name_" + fmt.Sprint(i)
			description       = "group name description " + fmt.Sprint(i)
			groupName         = ""
		)
		if groupNames == nil || i > len(groupNames)-1 {
			groupName = group_name
		} else {
			groupName = groupNames[i]
		}
		group := &cognitoidentityprovider.GroupType{}
		group.CreationDate = &timestamp
		group.Description = &description
		group.GroupName = &groupName
		group.LastModifiedDate = &timestamp
		group.Precedence = &randomNum
		group.UserPoolId = &userPoolId

		groupList.Groups = append(groupList.Groups, group)
		groupList.NextToken = &nextToken
	}
	return groupList
}
