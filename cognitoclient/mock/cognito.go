package mock

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type MockCognitoIdentityProviderClient struct {
    cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

// mocked functions

func (m *MockCognitoIdentityProviderClient) DescribeUserPool(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	exsitingPoolID := "us-west-2_aaaaaaaaa"
    if *poolInputData.UserPoolId != exsitingPoolID {
		return nil, errors.New("Failed to load user pool data")
	} 
	return nil, nil
}


