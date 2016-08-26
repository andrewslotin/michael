package authtest

// StaticTokenSource implements auth.TokenGenerator and always returns itself as a token.
type StaticTokenSource string

// Generate is needed to conform auth.TokenGenerator interface and returns StaticTokenSource
// as a string ignoring provided token length.
func (src StaticTokenSource) Generate(tokenLen int) string {
	return string(src)
}
