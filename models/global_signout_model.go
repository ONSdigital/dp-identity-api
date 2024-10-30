package models

import (
	"time"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type GlobalSignOut struct {
	ResultsChannel  chan string
	BackoffSchedule []time.Duration
	RetryAllowed    bool
}

// buildSignOutUserRequest - standalone request builder - builds a signout request array
//
//	this is required for concurrent global signout requests
func (g GlobalSignOut) BuildSignOutUserRequest(users *[]UserParams, userPoolId *string) []*cognitoidentityprovider.AdminUserGlobalSignOutInput {
	var usersDataArray []*cognitoidentityprovider.AdminUserGlobalSignOutInput
	userData := *users
	for i := 0; i < len(userData); i++ {
		userName := userData[i].ID
		usersDataArray = append(
			usersDataArray,
			&cognitoidentityprovider.AdminUserGlobalSignOutInput{
				UserPoolId: userPoolId,
				Username:   &userName,
			},
		)
	}
	return usersDataArray
}
