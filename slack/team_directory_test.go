package slack_test

import (
	"errors"
	"testing"

	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTeamDirectory_Fetch_ExistingUser(t *testing.T) {
	team := []slack.User{
		{ID: "U1", Name: "user1"},
		{ID: "U2", Name: "user2"},
		{ID: "U3", Name: "user3"},
	}

	api := new(apiMock)
	api.On("ListUsers").Return(team, nil)

	users := slack.NewTeamDirectory(api)

	for _, member := range team {
		user, err := users.Fetch(member.Name)
		if assert.NoError(t, err) {
			assert.Equal(t, member, user)
		}
	}

	api.AssertExpectations(t)
	api.AssertNumberOfCalls(t, "ListUsers", 1)
}

func TestTeamDirectory_Fetch_NonExistingUser(t *testing.T) {
	team := []slack.User{
		{ID: "U1", Name: "user1"},
	}

	api := new(apiMock)
	api.On("ListUsers").Return(team, nil)

	users := slack.NewTeamDirectory(api)

	_, err := users.Fetch("user2")
	assert.IsType(t, slack.NoSuchUserError{}, err)
	assert.EqualError(t, err, "there is no user with username 'user2' in the team")

	api.AssertExpectations(t)
	api.AssertNumberOfCalls(t, "ListUsers", 1)
}

func TestTeamDirectory_Fetch_WebAPIError(t *testing.T) {
	api := new(apiMock)
	api.On("ListUsers").Return(nil, errors.New("Slack Web API has returned an error"))

	users := slack.NewTeamDirectory(api)

	_, err := users.Fetch("user1")
	assert.EqualError(t, err, "failed to fetch team users: Slack Web API has returned an error")

	api.AssertExpectations(t)
	api.AssertNumberOfCalls(t, "ListUsers", 1)
}

type apiMock struct {
	mock.Mock
}

func (m *apiMock) ListUsers() ([]slack.User, error) {
	args := m.Called()

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).([]slack.User), nil
}
