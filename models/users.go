package models

type CognitoUser struct {
	TemporaryPassword      *string
	UserAttributes         []*AttributeType
	UserPoolId             *string
	Username               *string
	DesiredDeliveryMediums []*string
}

type AttributeType struct {
	Name  *string
	Value *string
}
type UserOutput struct {
	UserAttributes []*AttributeType
	Username       *string
}
