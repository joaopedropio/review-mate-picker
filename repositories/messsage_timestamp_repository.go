package repositories

func NewMessageTimestampCache() *MessageTimestampCache {
	return &MessageTimestampCache{
		messageTimestamps: make(map[string]string),
	}
}

type MessageTimestampCache struct {
	messageTimestamps map[string]string
}

func (r *MessageTimestampCache) Set(ts string) {
	r.messageTimestamps[ts] = ts
}

func (r *MessageTimestampCache) IsSet(ts string) bool {
	_, ok := r.messageTimestamps[ts]
	return ok
}
