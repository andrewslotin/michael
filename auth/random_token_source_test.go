package auth_test

import (
	"sync"
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/auth"
	"github.com/stretchr/testify/assert"
)

func TestRandomTokenSource_IssueToken_Uniqueness(t *testing.T) {
	const tokensNum = 256

	var src auth.RandomTokenSource

	tokens := make(map[string]int)
	for i := 0; i < tokensNum; i++ {
		token := src.Generate(10)
		dup, ok := tokens[token]

		assert.False(t, ok, "Token #%d duplicates token #%d", i, dup)
		tokens[token] = i
	}

	assert.Len(t, tokens, tokensNum)
}

func TestRandomTokenSource_IssueToken_Length(t *testing.T) {
	var src auth.RandomTokenSource

	for _, n := range [...]uint{0, 1, 2, 3, 4, 5, 6, 7, 8} {
		l := 1<<n - 1
		tok := src.Generate(l)
		assert.Len(t, tok, l)
	}
}

func TestRandomTokenSource_IssueToken_Concurrent(t *testing.T) {
	const tokensNum = 1000

	var src auth.RandomTokenSource

	done := make(chan struct{}, 1)
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < tokensNum; i++ {
			wg.Add(1)
			go func() {
				_ = src.Generate(10)
				wg.Done()
			}()
		}
		wg.Wait()

		done <- struct{}{}
	}()

	timer := time.NewTimer(3 * time.Second)
	select {
	case <-done:
	case <-timer.C:
		t.Errorf("timeout after 3s")
	}
}
