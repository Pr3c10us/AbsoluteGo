package queue

type Interface interface {
	Publish(params *MessageParams) error
}
