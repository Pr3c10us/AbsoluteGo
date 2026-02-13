package event

type Interface interface {
	Create(event Event) (int64, error)
	UpdateStatus(id int64, status Status) error
	GetEvents(filter Filter) ([]Event, error)
	GetEvent(id int64) (*Event, error)
}
