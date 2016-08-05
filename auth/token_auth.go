package auth

// TokenAuthorizer is an interface that wraps Authorize method.
//
// Authorize is used to check token authenticity.
type TokenAuthorizer interface {
	Authorize(token string) bool
}

// TokenIssuer is an interface that wraps IssueToken method.
//
// IssueToken is used to generate a new token of given length.
type TokenIssuer interface {
	IssueToken(tokenLen int) (token string, err error)
}
