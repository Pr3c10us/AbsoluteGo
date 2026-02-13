package queue

type Repository interface {
	Publish(params *MessageParams) error
}
