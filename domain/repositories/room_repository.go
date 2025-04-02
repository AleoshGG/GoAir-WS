package repositories

type MessageConsumer interface {
	ConsumeMessages(handler func(placeID string, data []byte)) error
}