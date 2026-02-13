package query

import "github.com/Pr3c10us/absolutego/internals/domains/event"

type GetEvents struct {
	eventImplementation event.Interface
}

func (service *GetEvents) Handle(filter event.Filter) ([]event.Event, error) {
	return service.eventImplementation.GetEvents(filter)
}

func NewGetEvents(eventImplementation event.Interface) *GetEvents {
	return &GetEvents{
		eventImplementation,
	}
}
