package event

import (
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/services/event/query"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
}

type Queries struct {
	GetEvents *query.GetEvents
}

func NewEventService(eventImplementation event.Interface) Services {
	return Services{
		Commands: Commands{},
		Queries: Queries{
			GetEvents: query.NewGetEvents(eventImplementation),
		},
	}
}
