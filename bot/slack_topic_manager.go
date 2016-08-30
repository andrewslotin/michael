package bot

import (
	"log"
	"strings"

	"github.com/andrewslotin/michael/slack"
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
	currentTopic, err := mgr.api.GetChannelTopic(channelID)
	if err == nil {
		newTopic := strings.Replace(currentTopic, DeployDoneEmotion, DeployInProgressEmotion, -1)
		if newTopic != currentTopic {
			err = mgr.api.SetChannelTopic(channelID, newTopic)
		}
	}

	if err != nil {
		log.Printf("slack-topic-manager: %s", err)
	}
}

func (mgr *SlackTopicManager) DeployCompleted(channelID string) {
	currentTopic, err := mgr.api.GetChannelTopic(channelID)
	if err == nil {
		newTopic := strings.Replace(currentTopic, DeployInProgressEmotion, DeployDoneEmotion, -1)
		if newTopic != currentTopic {
			err = mgr.api.SetChannelTopic(channelID, newTopic)
		}
	}

	if err != nil {
		log.Printf("slack-topic-manager: %s", err)
	}
}
