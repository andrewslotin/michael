package deploy

import "time"

type ChannelDeploys struct {
	store Store
}

func NewChannelDeploys(store Store) *ChannelDeploys {
	return &ChannelDeploys{store: store}
}

func (repo *ChannelDeploys) Current(channelID string) (Deploy, bool) {
	return repo.store.Get(channelID)
}

func (repo *ChannelDeploys) Start(channelID string, d Deploy) (Deploy, bool) {
	if current, ok := repo.Current(channelID); ok {
		return current, false
	}

	d.StartedAt = time.Now()
	repo.store.Set(channelID, d)

	return d, true
}

func (repo *ChannelDeploys) Finish(channelID string) (Deploy, bool) {
	return repo.store.Del(channelID)
}
