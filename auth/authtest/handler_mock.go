package authtest

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// HandlerMock implements http.Handler and is intended for use in tests.
type HandlerMock struct {
	mock.Mock
}

// ServeHTTP is needed to conform http.Handler interface and responds with HTTP 200 and an empty body.
func (m HandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.Write(nil)
}
