package authtest

import "github.com/stretchr/testify/mock"

// TokenSourceMock implements auth.TokenGenerator and is intended for using in tests.
type TokenSourceMock struct {
	mock.Mock
}

// Generate is needed to conform auth.TokenGenerator interface and returns mocked value.
func (src *TokenSourceMock) Generate(tokenLen int) string {
	return src.Called(tokenLen).Get(0).(string)
}
