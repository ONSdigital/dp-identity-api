package mock

func (m *CognitoIdentityProviderClientStub) AddUserPool(poolId string) {
	m.UserPools = append(m.UserPools, poolId)
}
