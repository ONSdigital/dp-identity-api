package mock

import (
	"crypto/rand"
	"fmt"
	"math/big"
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
	// Type to map for the Cognito GroupType object
	var (
		groupsList  cognitoidentityprovider.ListGroupsOutput
		nextToken   = "next_token"
		totalGroups = 0
	)

	// Preallocate the slice for groups if you have an estimated maximum size
	m.GroupsList = make([]cognitoidentityprovider.ListGroupsOutput, 0, (groupCount/60)+1)

	for i := 1; i <= groupCount; i++ {
		D := "group name description " + fmt.Sprint(i)
		G := "group_name_" + fmt.Sprint(i)

		// Generate a secure random precedence value between 13 and 100
		P, err := secureRandomInt(13, 100)
		if err != nil {
			fmt.Println("Error generating random number:", err)
			return
		}

		group := cognitoidentityprovider.GroupType{
			Description: &D,
			GroupName:   &G,
			Precedence:  &P,
		}

		groupsList.Groups = append(groupsList.Groups, &group)

		if i%60 == 0 {
			groupsList.NextToken = &nextToken
			totalGroups += len(groupsList.Groups)
			m.GroupsList = append(m.GroupsList, groupsList)
			groupsList = cognitoidentityprovider.ListGroupsOutput{} // Reset for the next batch
		}
	}
	if len(groupsList.Groups) > 0 {
		// If there are remaining groups in the last batch
		groupsList.NextToken = nil
		m.GroupsList = append(m.GroupsList, groupsList)
	}
}

// secureRandomInt generates a secure random integer in the range [minValue, maxValue].
func secureRandomInt(minValue, maxValue int64) (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(maxValue-minValue+1))
	if err != nil {
		return 0, err
	}
	return n.Int64() + minValue, nil
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
	parsedTime, _ := time.Parse("2006-Jan-1", "2010-Jan-1")
	emptyString := ""
	num := int64(1)

	return cognitoidentityprovider.ListGroupsOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{
				Description:      &description,
				CreationDate:     &parsedTime,
				GroupName:        &emptyString,
				LastModifiedDate: &parsedTime,
				Precedence:       &num,
				RoleArn:          &emptyString,
				UserPoolId:       &emptyString,
			},
		},
	}, nil
}
