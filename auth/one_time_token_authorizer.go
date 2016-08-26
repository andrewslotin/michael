package auth

import (
	"errors"
	"sync"
)

// OneTimeTokenAuthorizer issues tokens that can be used for authorization only once.
type OneTimeTokenAuthorizer struct {
	gen    TokenGenerator
	mu     sync.Mutex
	tokens map[string]struct{}
}

// NewOneTimeTokenAuthorizer returns an instance of *OneTimeTokenAuthorizer that uses src as a token source.
func NewOneTimeTokenAuthorizer(src TokenGenerator) *OneTimeTokenAuthorizer {
	return &OneTimeTokenAuthorizer{
		gen:    src,
		tokens: make(map[string]struct{}, 50),
	}
}

// IssueToken generates and stores a new unused token. This method returns an error if it failed to generate
// an unused token after 16777216 (2^24) attempts.
func (s *OneTimeTokenAuthorizer) IssueToken(tokenLen int) (token string, err error) {
	const maxAttempts = 1 << 20

	s.mu.Lock()
	defer s.mu.Unlock()

	token = s.gen.Generate(tokenLen)
	for i := 0; s.exists(token) && i < maxAttempts; i++ {
		token = s.gen.Generate(tokenLen)
	}

	if s.exists(token) {
		return "", errors.New("failed to generate token")
	}

	s.tokens[token] = struct{}{}

	return token, nil
}

// Authorize checks if provided token has been issued by this instance of OneTimeTokenAuthorizer and annuls it.
func (s *OneTimeTokenAuthorizer) Authorize(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.exists(token) {
		return false
	}

	delete(s.tokens, token)
	return true
}

func (s *OneTimeTokenAuthorizer) exists(token string) bool {
	_, ok := s.tokens[token]

	return ok
}
