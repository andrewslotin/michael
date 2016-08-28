package authtest

import "github.com/stretchr/testify/mock"

// TokenAuthenticatorMock imprements auth.TokenAuthenticator and is intended for using in tests.
type TokenAuthenticatorMock struct {
	mock.Mock
}

// Authenticate is needed to conform auth.TokenAuthenticator interface and returns mocked value.
func (m TokenAuthenticatorMock) Authenticate(token string) bool {
	return m.Called(token).Get(0).(bool)
}
