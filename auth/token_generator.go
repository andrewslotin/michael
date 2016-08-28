package auth

// TokenGenerator is an interface that wraps Generate method.
//
// Generate is used to generate strings of given length and is used to issue tokens.
type TokenGenerator interface {
	Generate(tokenLen int) string
}
