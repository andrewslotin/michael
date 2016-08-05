package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andrewslotin/slack-deploy-command/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/*      Test objects      */
type tokenAuthorizerMock struct {
	mock.Mock
}

func (m tokenAuthorizerMock) Authorize(token string) bool {
	return m.Called(token).Get(0).(bool)
}

type handlerMock struct {
	mock.Mock
}

func (m handlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.Write(nil)
}

/*          Tests         */
func TestTokenAuthMiddleware_ValidToken(t *testing.T) {
	token := "token1"

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/?token="+token, nil)
	require.NoError(t, err)

	var (
		handler    handlerMock
		authorizer tokenAuthorizerMock
	)
	handler.On("ServeHTTP", recorder, req).Return().Once()
	authorizer.On("Authorize", token).Return(true)

	auth.TokenAuthMiddleware(handler, authorizer).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	handler.AssertExpectations(t)
	authorizer.AssertExpectations(t)
}

func TestTokenAuthMiddleware_InvalidToken(t *testing.T) {
	token := "token1"

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/?token="+token, nil)
	require.NoError(t, err)

	var (
		handler    handlerMock
		authorizer tokenAuthorizerMock
	)
	authorizer.On("Authorize", token).Return(false)

	auth.TokenAuthMiddleware(handler, authorizer).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), strings.TrimSpace(recorder.Body.String()))

	handler.AssertExpectations(t)
	authorizer.AssertExpectations(t)
}

func TestTokenAuthMiddleware_MissingToken(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	var (
		handler    handlerMock
		authorizer tokenAuthorizerMock
	)

	auth.TokenAuthMiddleware(handler, authorizer).ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "Missing token", strings.TrimSpace(recorder.Body.String()))

	handler.AssertExpectations(t)
}
