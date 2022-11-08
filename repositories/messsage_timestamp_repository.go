package repositories

import "sync"

func NewMessageTimestampCache() *MessageTimestampCache {
	return &MessageTimestampCache{
		messageTimestamps: sync.Map{},
	}
}

type MessageTimestampCache struct {
	messageTimestamps sync.Map
}

func (r *MessageTimestampCache) Set(ts string) {
	r.messageTimestamps.Store(ts, ts)
}

func (r *MessageTimestampCache) IsSet(ts string) bool {
	_, ok := r.messageTimestamps.Load(ts)
	return ok
}
