package bot

import (
	"fmt"
	"log"

	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
)

type SlackIMNotifier struct {
	im    *slack.InstantMessenger
	users *slack.TeamDirectory
}

func NewSlackIMNotifier(api *slack.WebAPI) *SlackIMNotifier {
	return &SlackIMNotifier{
		im:    slack.NewInstantMessenger(api),
		users: slack.NewTeamDirectory(api),
	}
}

func (notifier *SlackIMNotifier) DeployStarted(_ string, _ deploy.Deploy) {}

func (notifier *SlackIMNotifier) DeployCompleted(_ string, d deploy.Deploy) {
	for _, userRef := range d.InterestedUsers {
		user, err := notifier.users.Fetch(userRef.Name)
		if err != nil {
			if _, ok := err.(slack.NoSuchUserError); !ok {
				log.Printf("cannot notify %s about completed deploy of %s: %s", user.Name, d.Subject, err)
			}

			continue
		}

		message := slack.Message{
			Text: fmt.Sprintf("%s just deployed %s", d.User, d.Subject),
		}

		err = notifier.im.SendMessage(user, message)
		if err != nil {
			log.Printf("failed to send an instant message to %s: %s", user.Name, err)
			continue
		}
	}
}
