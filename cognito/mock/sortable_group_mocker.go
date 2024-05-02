package mock

import (
	"time"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func (m *CognitoIdentityProviderClientStub) AddSortableGroupWithName(name string) error {
	newGroup, err := m.GenerateSortableGroup(name)
	if err != nil {
		return err
	}
	m.GroupsList = append(m.GroupsList, newGroup)
	return nil
}

func (m *CognitoIdentityProviderClientStub) GenerateSortableGroup(description string) (cognitoidentityprovider.ListGroupsOutput, error) {
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
