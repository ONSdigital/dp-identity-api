package models

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Precedence  string `json:"precedence"`
}

func (g *Group) BuildGetGroupRequest(userPoolId string) *cognitoidentityprovider.GetGroupInput {
	return &cognitoidentityprovider.GetGroupInput{
		GroupName:  &g.Name,
		UserPoolId: &userPoolId,
	}
}
