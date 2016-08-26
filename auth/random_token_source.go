package auth

import (
	"math/rand"
	"sync"
)

type RandomTokenSource struct {
	Src rand.Source
	mu  sync.Mutex
}

// Credits for this great solution go to Stack Overflow user icza.
// See his answer http://stackoverflow.com/a/31832326 for explanation.
func (gen *RandomTokenSource) Generate(tokenLen int) string {
	const (
		chars       = "abcdefghikjlmnopqurstuvwxyzABCDEFGHIKJLMNOPQURSTUVWXYZ0123456789"
		charIdxBits = 6 // number of bits needed to store an index: log2(len(chars)-1)
		charIdxMask = 1<<charIdxBits - 1
		charIdxMax  = 63 / charIdxBits
	)

	if tokenLen < 1 {
		return ""
	}

	token := make([]byte, tokenLen)
	for i, cache, remain := tokenLen, gen.getRandInt63(), charIdxMax; i > 0; {
		if remain == 0 {
			cache, remain = gen.getRandInt63(), charIdxMax
		}
		if idx := int(cache & charIdxMask); idx < len(chars) {
			token[i-1] = chars[idx]
			i--
		}
		cache >>= charIdxBits
		remain--
	}

	return string(token)
}

func (gen *RandomTokenSource) getRandInt63() int64 {
	gen.mu.Lock()
	defer gen.mu.Unlock()

	if gen.Src == nil {
		return rand.Int63()
	}

	return gen.Src.Int63()
}
