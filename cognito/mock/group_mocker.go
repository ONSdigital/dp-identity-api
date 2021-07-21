package mock

import "time"

type Group struct {
	Name        string
	Description string
	Precedence  int64
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

func (m *CognitoIdentityProviderClientStub) AddGroupWithNameAndDescription(name, description string) error {
	newGroup, err := m.GenerateGroup(name, description, 0)
	if err != nil {
		return err
	}
	m.Groups = append(m.Groups, newGroup)
	return nil
}

func (m *CognitoIdentityProviderClientStub) AddGroupWithNameAndPrecedence(name string, precedence int64) error {
	newGroup, err := m.GenerateGroup(name, "", precedence)
	if err != nil {
		return err
	}
	m.Groups = append(m.Groups, newGroup)
	return nil
}

func (m *CognitoIdentityProviderClientStub) GenerateGroup(name, description string, precedence int64) (*Group, error) {
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
