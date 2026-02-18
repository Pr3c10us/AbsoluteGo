package event

type Interface interface {
	Create(event Event) (int64, error)
	Update(id int64, event Event) error
	GetEvents(filter Filter) ([]Event, error)
	GetEvent(id int64) (*Event, error)
}
