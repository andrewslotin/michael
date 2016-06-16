package server

import (
	"strings"

	"github.com/andrewslotin/slack-deploy-command/slack"
)

const (
	DeployInProgressEmotion = ":no_entry:"
	DeployDoneEmotion       = ":white_check_mark:"
)

type SlackTopicManager struct {
	api *slack.WebAPI
}

func NewSlackTopicManager(webAPIClient *slack.WebAPI) *SlackTopicManager {
	return &SlackTopicManager{api: webAPIClient}
}

func (mgr *SlackTopicManager) DeployStarted(channelID string) error {
	topic, err := mgr.api.GetChannelTopic(channelID)
	if err != nil {
		return err
	}

	return mgr.api.SetChannelTopic(channelID, strings.Replace(topic, DeployDoneEmotion, DeployInProgressEmotion, -1))
}

func (mgr *SlackTopicManager) DeployCompleted(channelID string) error {
	topic, err := mgr.api.GetChannelTopic(channelID)
	if err != nil {
		return err
	}

	return mgr.api.SetChannelTopic(channelID, strings.Replace(topic, DeployInProgressEmotion, DeployDoneEmotion, -1))
}
