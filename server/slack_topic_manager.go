package server

import (
	"log"
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

func (mgr *SlackTopicManager) DeployStarted(channelID string) {
	topic, err := mgr.api.GetChannelTopic(channelID)
	if err == nil {
		err = mgr.api.SetChannelTopic(channelID, strings.Replace(topic, DeployDoneEmotion, DeployInProgressEmotion, -1))
	}

	if err != nil {
		log.Printf("slack-topic-manager: %s", err)
	}
}

func (mgr *SlackTopicManager) DeployCompleted(channelID string) {
	topic, err := mgr.api.GetChannelTopic(channelID)
	if err == nil {
		err = mgr.api.SetChannelTopic(channelID, strings.Replace(topic, DeployInProgressEmotion, DeployDoneEmotion, -1))
	}

	if err != nil {
		log.Printf("slack-topic-manager: %s", err)
	}
}
