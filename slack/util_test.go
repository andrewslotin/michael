package slack_test

import (
	"testing"

	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
)

func TestEscapeMessage(t *testing.T) {
	examples := map[string]struct{ Value, Expected string }{
		"common":   {`"Hello' & <<world>>!`, `"Hello' &amp; &lt;&lt;world&gt;&gt;!`},
		"user_ref": {"Hello <<@U123456|user1>>!", "Hello &lt;<@U123456|user1>&gt;!"},
		"link":     {"Check <<https://google.com?q=search+term&source=Chrome>>", "Check &lt;<https://google.com?q=search+term&source=Chrome>&gt;"},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, example.Expected, slack.EscapeMessage(example.Value))
		})
	}
}
