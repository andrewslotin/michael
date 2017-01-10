package slack

import "sync"

type messageSender interface {
	OpenIMChannel(User) (string, error)
	PostMessage(string, Message) error
}

type InstantMessenger struct {
	api messageSender

	mu       sync.RWMutex
	channels map[string]string
}

func NewInstantMessenger(api messageSender) *InstantMessenger {
	return &InstantMessenger{
		api:      api,
		channels: make(map[string]string),
	}
}

func (im *InstantMessenger) SendMessage(user User, message Message) error {
	channelID, err := im.openChannel(user)
	if err != nil {
		return err
	}

	return im.api.PostMessage(channelID, message)
}

func (im *InstantMessenger) openChannel(user User) (string, error) {
	im.mu.RLock()
	channelID, ok := im.channels[user.Name]
	im.mu.RUnlock()

	if ok {
		return channelID, nil
	}

	im.mu.Lock()
	defer im.mu.Unlock()

	channelID, ok = im.channels[user.Name]
	if ok {
		return channelID, nil
	}

	channelID, err := im.api.OpenIMChannel(user)
	if err != nil {
		return "", err
	}

	im.channels[user.Name] = channelID

	return im.channels[user.Name], nil
}
