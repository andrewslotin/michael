package slack_test

import (
	"testing"

	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
)

func TestEscapeMessage(t *testing.T) {
	s := `"Hello' & <<world>>!`

	assert.Equal(t, `"Hello' &amp; &lt;&lt;world&gt;&gt;!`, slack.EscapeMessage(s))
}
