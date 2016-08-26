package authtest

import "github.com/stretchr/testify/mock"

// TokenAuthorizerMock imprements auth.TokenAuthorizer and is intended for using in tests.
type TokenAuthorizerMock struct {
	mock.Mock
}

// Authorize is needed to conform auth.TokenAuthorizer interface and returns mocked value.
func (m TokenAuthorizerMock) Authorize(token string) bool {
	return m.Called(token).Get(0).(bool)
}
