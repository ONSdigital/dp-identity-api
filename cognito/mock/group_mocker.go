package mock

import (
	"errors"
	"fmt"
	"time"
)

type Group struct {
	Name        string
	Description string
	Precedence  int32
	Created     time.Time
	Members     []*User
}

func (m *CognitoIdentityProviderClientStub) AddGroupWithName(name string) error {
	newGroup, err := m.GenerateGroup(name, "", 0)
	if err != nil {
		return err
	}
	m.Groups = append(m.Groups, newGroup)
	return nil
}

func (m *CognitoIdentityProviderClientStub) AddGroupWithNameAndPrecedence(name string, precedence int32) error {
	newGroup, err := m.GenerateGroup(name, "", precedence)
	if err != nil {
		return err
	}
	m.Groups = append(m.Groups, newGroup)
	return nil
}

func (m *CognitoIdentityProviderClientStub) GenerateGroup(name, description string, precedence int32) (*Group, error) {
	if name == "" {
		name = "TestGroup"
	}
	if description == "" {
		description = "A test group"
	}
	if precedence == 0 {
		precedence = 100
	}

	createdTime, err := time.Parse("2006-Jan-1", "2010-Jan-1")
	if err != nil {
		return nil, err
	}

	return &Group{
		Name:        name,
		Description: description,
		Precedence:  precedence,
		Created:     createdTime,
		Members:     []*User{},
	}, nil
}

func (m *CognitoIdentityProviderClientStub) ReadGroup(groupName string) *Group {
	for _, group := range m.Groups {
		if group.Name == groupName {
			return group
		}
	}
	return nil
}

func (m *CognitoIdentityProviderClientStub) AddUserToGroup(username, groupName string) error {
	group := m.ReadGroup(groupName)
	if group == nil {
		return errors.New("could not find the group")
	}

	user := m.ReadUser(username)
	if user == nil {
		return errors.New("could not find the group")
	}

	user.Groups = append(user.Groups, group)
	group.Members = append(group.Members, user)
	return nil
}

// BulkGenerateGroups - bulk generate 'n' groups for testing purposes
//                    if groupnames array is nil or length is different, will autofill with

func (m *CognitoIdentityProviderClientStub) BulkGenerateGroups(groupCount int) {
	for i := 0; i < groupCount; i++ {
		D := "group name description " + fmt.Sprint(i)
		G := "group_name_" + fmt.Sprint(i)
		P := int32(i + 13)

		group := Group{
			Name:        G,
			Description: D,
			Precedence:  P,
		}

		m.Groups = append(m.Groups, &group)
	}
}
