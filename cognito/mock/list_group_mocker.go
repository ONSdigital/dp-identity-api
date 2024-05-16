package mock

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

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

func (m *CognitoIdentityProviderClientStub) AddListGroupWithName(name string) error {
	newGroup, err := m.GenerateListGroup(name)
	if err != nil {
		return err
	}
	m.GroupsList = append(m.GroupsList, newGroup)
	return nil
}

func (m *CognitoIdentityProviderClientStub) GenerateListGroup(description string) (cognitoidentityprovider.ListGroupsOutput, error) {
	time, _ := time.Parse("2006-Jan-1", "2010-Jan-1")
	emptyString := ""
	num := int64(1)

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
