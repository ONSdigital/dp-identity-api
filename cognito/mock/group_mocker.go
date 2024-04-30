package mock

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

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

func (m *CognitoIdentityProviderClientStub) AddGroupWithNameSort(name string) error {
	newGroup, err := m.GenerateGroupSort(name)
	if err != nil {
		return err
	}
	m.GroupsList = append(m.GroupsList, newGroup)
	return nil
}

func (m *CognitoIdentityProviderClientStub) AddGroupWithNameAndDescription(name, description string) error {
	newGroup, err := m.GenerateGroup(name, description, 0)
	if err != nil {
		return err
	}
	m.Groups = append(m.Groups, newGroup)

	pres := int64(0)
	newGroupType := cognitoidentityprovider.GroupType{
		CreationDate:     nil,
		Description:      &description,
		GroupName:        &name,
		LastModifiedDate: nil,
		Precedence:       &pres,
		RoleArn:          nil,
		UserPoolId:       nil,
	}
	newGroupTypeList := []*cognitoidentityprovider.GroupType{}
	newGroupTypeList = append(newGroupTypeList, &newGroupType)
	groupsListOutput := cognitoidentityprovider.ListGroupsOutput{
		Groups:    newGroupTypeList,
		NextToken: nil,
	}
	m.GroupsList = append(m.GroupsList, groupsListOutput)

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

func (m *CognitoIdentityProviderClientStub) GenerateGroupSort(description string) (cognitoidentityprovider.ListGroupsOutput, error) {
	time, _ := time.Parse("2006-Jan-1", "2010-Jan-1")
	emptyString := ""
	var num int64
	num = 1
	return cognitoidentityprovider.ListGroupsOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{
				Description:      &description,
				CreationDate:     &time,
				GroupName:        &emptyString,
				LastModifiedDate: &time,
				Precedence:       &num,
				RoleArn:          &emptyString,
				UserPoolId:       &emptyString,
			},
		},
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

func (m *CognitoIdentityProviderClientStub) AddUserToGroup(username string, groupName string) error {
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

//BulkGenerateGroups - bulk generate 'n' groups for testing purposes
//                    if groupnames array is nil or length is different, will autofill with

func (m *CognitoIdentityProviderClientStub) BulkGenerateGroups(groupCount int) {

	for i := 0; i < groupCount; i++ {
		D := "group name description " + fmt.Sprint(i)
		G := "group_name_" + fmt.Sprint(i)
		P := int64(i + 13)

		group := Group{
			Name:        G,
			Description: D,
			Precedence:  P,
		}

		m.Groups = append(m.Groups, &group)

	}
}

func (m *CognitoIdentityProviderClientStub) BulkGenerateGroupsList(groupCount int) {
	//Type to map for the Cognito GroupType object

	var (
		groupsList  cognitoidentityprovider.ListGroupsOutput
		next_token  string = "next_token"
		totalgroups        = 0
	)

	for i := 1; i <= groupCount; i++ {
		D := "group name description " + fmt.Sprint(i)
		G := "group_name_" + fmt.Sprint(i)
		P := rand.Int63n(100-13) + 13

		group := cognitoidentityprovider.GroupType{
			Description: &D,
			GroupName:   &G,
			Precedence:  &P,
		}

		groupsList.Groups = append(groupsList.Groups, &group)

		if i%60 == 0 && i > 0 {
			groupsList.NextToken = &next_token
			totalgroups = totalgroups + len(groupsList.Groups)
			m.GroupsList = append(m.GroupsList, groupsList)
			groupsList = *new(cognitoidentityprovider.ListGroupsOutput)
		}

	}
	groupsList.NextToken = nil
	m.GroupsList = append(m.GroupsList, groupsList)

}
