package mock

type Group struct {
	Name        string
	Description string
	Precedence  int64
	Members     []*User
}

func (m *CognitoIdentityProviderClientStub) AddGroupWithName(name string) {
	m.Groups = append(m.Groups, m.GenerateGroup(name, "", 0))
}

func (m *CognitoIdentityProviderClientStub) AddGroupWithNameAndDescription(name, description string) {
	m.Groups = append(m.Groups, m.GenerateGroup(name, description, 0))
}

func (m *CognitoIdentityProviderClientStub) AddGroupWithNameAndPrecedence(name string, precedence int64) {
	m.Groups = append(m.Groups, m.GenerateGroup(name, "", precedence))
}

func (m *CognitoIdentityProviderClientStub) GenerateGroup(name, description string, precedence int64) *Group {
	if name == "" {
		name = "TestGroup"
	}
	if description == "" {
		description = "A test group"
	}
	if precedence == 0 {
		precedence = 100
	}

	return &Group{
		Name:        name,
		Description: description,
		Precedence:  precedence,
		Members:     []*User{},
	}
}

func (m *CognitoIdentityProviderClientStub) ReadGroup(groupName string) *Group {
	for _, group := range m.Groups {
		if group.Name == groupName {
			return group
		}
	}
	return nil
}
