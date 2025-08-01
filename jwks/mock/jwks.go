// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"sync"

	"github.com/ONSdigital/dp-identity-api/v2/jwks"
)

// Ensure, that ManagerMock does implement jwks.Manager.
// If this is not the case, regenerate this file with moq.
var _ jwks.Manager = &ManagerMock{}

// ManagerMock is a mock implementation of jwks.Manager.
//
//	func TestSomethingThatUsesManager(t *testing.T) {
//
//		// make and configure a mocked jwks.Manager
//		mockedManager := &ManagerMock{
//			JWKSGetKeysetFunc: func(awsRegion string, poolID string) (*jwks.JWKS, error) {
//				panic("mock out the JWKSGetKeyset method")
//			},
//			JWKSToRSAJSONResponseFunc: func(jwksMoqParam *jwks.JWKS) ([]byte, error) {
//				panic("mock out the JWKSToRSAJSONResponse method")
//			},
//		}
//
//		// use mockedManager in code that requires jwks.Manager
//		// and then make assertions.
//
//	}
type ManagerMock struct {
	// JWKSGetKeysetFunc mocks the JWKSGetKeyset method.
	JWKSGetKeysetFunc func(awsRegion string, poolID string) (*jwks.JWKS, error)

	// JWKSToRSAJSONResponseFunc mocks the JWKSToRSAJSONResponse method.
	JWKSToRSAJSONResponseFunc func(jwksMoqParam *jwks.JWKS) ([]byte, error)

	// calls tracks calls to the methods.
	calls struct {
		// JWKSGetKeyset holds details about calls to the JWKSGetKeyset method.
		JWKSGetKeyset []struct {
			// AwsRegion is the awsRegion argument value.
			AwsRegion string
			// PoolID is the poolID argument value.
			PoolID string
		}
		// JWKSToRSAJSONResponse holds details about calls to the JWKSToRSAJSONResponse method.
		JWKSToRSAJSONResponse []struct {
			// JwksMoqParam is the jwksMoqParam argument value.
			JwksMoqParam *jwks.JWKS
		}
	}
	lockJWKSGetKeyset         sync.RWMutex
	lockJWKSToRSAJSONResponse sync.RWMutex
}

// JWKSGetKeyset calls JWKSGetKeysetFunc.
func (mock *ManagerMock) JWKSGetKeyset(awsRegion string, poolID string) (*jwks.JWKS, error) {
	if mock.JWKSGetKeysetFunc == nil {
		panic("ManagerMock.JWKSGetKeysetFunc: method is nil but Manager.JWKSGetKeyset was just called")
	}
	callInfo := struct {
		AwsRegion string
		PoolID    string
	}{
		AwsRegion: awsRegion,
		PoolID:    poolID,
	}
	mock.lockJWKSGetKeyset.Lock()
	mock.calls.JWKSGetKeyset = append(mock.calls.JWKSGetKeyset, callInfo)
	mock.lockJWKSGetKeyset.Unlock()
	return mock.JWKSGetKeysetFunc(awsRegion, poolID)
}

// JWKSGetKeysetCalls gets all the calls that were made to JWKSGetKeyset.
// Check the length with:
//
//	len(mockedManager.JWKSGetKeysetCalls())
func (mock *ManagerMock) JWKSGetKeysetCalls() []struct {
	AwsRegion string
	PoolID    string
} {
	var calls []struct {
		AwsRegion string
		PoolID    string
	}
	mock.lockJWKSGetKeyset.RLock()
	calls = mock.calls.JWKSGetKeyset
	mock.lockJWKSGetKeyset.RUnlock()
	return calls
}

// JWKSToRSAJSONResponse calls JWKSToRSAJSONResponseFunc.
func (mock *ManagerMock) JWKSToRSAJSONResponse(jwksMoqParam *jwks.JWKS) ([]byte, error) {
	if mock.JWKSToRSAJSONResponseFunc == nil {
		panic("ManagerMock.JWKSToRSAJSONResponseFunc: method is nil but Manager.JWKSToRSAJSONResponse was just called")
	}
	callInfo := struct {
		JwksMoqParam *jwks.JWKS
	}{
		JwksMoqParam: jwksMoqParam,
	}
	mock.lockJWKSToRSAJSONResponse.Lock()
	mock.calls.JWKSToRSAJSONResponse = append(mock.calls.JWKSToRSAJSONResponse, callInfo)
	mock.lockJWKSToRSAJSONResponse.Unlock()
	return mock.JWKSToRSAJSONResponseFunc(jwksMoqParam)
}

// JWKSToRSAJSONResponseCalls gets all the calls that were made to JWKSToRSAJSONResponse.
// Check the length with:
//
//	len(mockedManager.JWKSToRSAJSONResponseCalls())
func (mock *ManagerMock) JWKSToRSAJSONResponseCalls() []struct {
	JwksMoqParam *jwks.JWKS
} {
	var calls []struct {
		JwksMoqParam *jwks.JWKS
	}
	mock.lockJWKSToRSAJSONResponse.RLock()
	calls = mock.calls.JWKSToRSAJSONResponse
	mock.lockJWKSToRSAJSONResponse.RUnlock()
	return calls
}
