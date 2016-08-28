package auth

// DefaultTokenLength is the default length for IssueToken
const DefaultTokenLength = 16

// TokenAuthenticator is an interface that wraps Authenticate method.
//
// Authenticate is used to check token authenticity.
type TokenAuthenticator interface {
	Authenticate(token string) bool
}

// TokenIssuer is an interface that wraps IssueToken method.
//
// IssueToken is used to generate a new token of given length.
type TokenIssuer interface {
	IssueToken(tokenLen int) (token string, err error)
}

// None is TokenAuthenticator and TokenIssuer that always issues an empty token and authenticates everything.
var None = noopAuth{}

type noopAuth struct{}

func (noopAuth) IssueToken(int) (string, error) {
	return "", nil
}

func (noopAuth) Authenticate(token string) bool {
	return true
}
