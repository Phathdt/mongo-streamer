package watermillapp

type Publisher interface {
	Publish(topic string, data interface{}) error
	PublishRaw(topic string, data []byte) error
}
