package bot

import (
	"log"
	"strings"

	"github.com/andrewslotin/michael/deploy"
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

func (mgr *SlackTopicManager) DeployStarted(channelID string, _ deploy.Deploy) {
	err := mgr.channelTopicReplace(channelID, DeployDoneEmotion, DeployInProgressEmotion)
	if err != nil {
		log.Printf("slack-topic-manager: %s", err)
	}
}

func (mgr *SlackTopicManager) DeployCompleted(channelID string, _ deploy.Deploy) {
	err := mgr.channelTopicReplace(channelID, DeployInProgressEmotion, DeployDoneEmotion)
	if err != nil {
		log.Printf("slack-topic-manager: %s", err)
	}
}

func (mgr *SlackTopicManager) channelTopicReplace(channelID, old, new string) error {
	currentTopic, err := mgr.api.GetChannelTopic(channelID)
	if err != nil {
		return err
	}

	newTopic := strings.Replace(currentTopic, old, new, -1)
	if newTopic == currentTopic {
		return nil
	}

	return mgr.api.SetChannelTopic(channelID, newTopic)
}
