package auth

import (
	"errors"
	"sync"
)

// OneTimeTokenAuthenticator issues tokens that can be used for authorization only once.
type OneTimeTokenAuthenticator struct {
	gen    TokenGenerator
	mu     sync.Mutex
	tokens map[string]struct{}
}

// NewOneTimeTokenAuthenticator returns an instance of *OneTimeTokenAuthenticator that uses src as a token source.
func NewOneTimeTokenAuthenticator(src TokenGenerator) *OneTimeTokenAuthenticator {
	return &OneTimeTokenAuthenticator{
		gen:    src,
		tokens: make(map[string]struct{}, 50),
	}
}

// IssueToken generates and stores a new unused token. This method returns an error if it failed to generate
// an unused token after 16777216 (2^24) attempts.
func (s *OneTimeTokenAuthenticator) IssueToken(tokenLen int) (token string, err error) {
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

// Authenticate checks if provided token has been issued by this instance of OneTimeTokenAuthenticator and annuls it.
func (s *OneTimeTokenAuthenticator) Authenticate(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.exists(token) {
		return false
	}

	delete(s.tokens, token)
	return true
}

func (s *OneTimeTokenAuthenticator) exists(token string) bool {
	_, ok := s.tokens[token]

	return ok
}
