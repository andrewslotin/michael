package slack_test

import (
	"errors"
	"testing"

	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInstantMessenger_SendMessage(t *testing.T) {
	user := slack.User{ID: "123", Name: "user1"}
	message := slack.Message{Text: "text message"}

	api := new(messageSenderMock)
	api.On("OpenIMChannel", user).Return("channel1", nil)
	api.On("PostMessage", "channel1", message).Return(nil)

	im := slack.NewInstantMessenger(api)
	require.NoError(t, im.SendMessage(user, message))

	api.AssertExpectations(t)
	api.AssertNumberOfCalls(t, "OpenIMChannel", 1)
	api.AssertNumberOfCalls(t, "PostMessage", 1)

	require.NoError(t, im.SendMessage(user, message))
	api.AssertNumberOfCalls(t, "OpenIMChannel", 1)
	api.AssertNumberOfCalls(t, "PostMessage", 2)
}

func TestInstantMessenger_SendMessage_OpenIMChannelError(t *testing.T) {
	user := slack.User{ID: "123", Name: "user1"}
	message := slack.Message{Text: "text message"}

	api := new(messageSenderMock)
	api.On("OpenIMChannel", user).Return("", errors.New("open_im_error"))

	im := slack.NewInstantMessenger(api)
	assert.Error(t, im.SendMessage(user, message))

	api.AssertExpectations(t)
	api.AssertNumberOfCalls(t, "OpenIMChannel", 1)

	assert.Error(t, im.SendMessage(user, message))
	api.AssertNumberOfCalls(t, "OpenIMChannel", 2)
}

func TestInstantMessenger_SendMessage_PostMessageError(t *testing.T) {
	user := slack.User{ID: "123", Name: "user1"}
	message := slack.Message{Text: "text message"}

	api := new(messageSenderMock)
	api.On("OpenIMChannel", user).Return("channel1", nil)
	api.On("PostMessage", "channel1", message).Return(errors.New("post_message_error"))

	im := slack.NewInstantMessenger(api)
	assert.Error(t, im.SendMessage(user, message))

	api.AssertExpectations(t)
	api.AssertNumberOfCalls(t, "OpenIMChannel", 1)
	api.AssertNumberOfCalls(t, "PostMessage", 1)

	assert.Error(t, im.SendMessage(user, message))
	api.AssertNumberOfCalls(t, "OpenIMChannel", 1)
	api.AssertNumberOfCalls(t, "PostMessage", 2)
}

type messageSenderMock struct {
	mock.Mock
}

func (m *messageSenderMock) OpenIMChannel(user slack.User) (string, error) {
	args := m.Called(user)

	return args.String(0), args.Error(1)
}

func (m *messageSenderMock) PostMessage(channelID string, message slack.Message) error {
	return m.Called(channelID, message).Error(0)
}
