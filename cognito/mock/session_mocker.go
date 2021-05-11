package mock

type Session struct {
	AccessToken		string
	IdToken			string
	RefreshToken	string
}

func (m *CognitoIdentityProviderClientStub) CreateSessionWithAccessToken(accessToken string) {
	idToken, refreshToken := "aaaa.bbbb.cccc", "dddd.eeee.ffff"
	m.Sessions = append(m.Sessions, m.GenerateSession(accessToken, idToken, refreshToken))
}

func (m *CognitoIdentityProviderClientStub) GenerateSession(accessToken string, idToken string, refreshToken string) Session {
	return Session{
		AccessToken: accessToken,
		IdToken: idToken,
		RefreshToken: refreshToken,
	}
}
