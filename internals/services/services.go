package services

import (
	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/services/book"
	"github.com/Pr3c10us/absolutego/internals/services/event"
	"github.com/Pr3c10us/absolutego/internals/services/script"
)

type Services struct {
	BookServices   book.Services
	ScriptServices script.Services
	EventServices  event.Services
}

func NewServices(adapters *adapters.Adapters) *Services {
	return &Services{
		BookServices:   book.NewBookServices(adapters.BookImplementation, adapters.StorageImplementation, adapters.AiImplementation, adapters.EnvironmentVariables, adapters.ScriptImplementation, adapters.EventImplementation, adapters.QueueImplementation),
		ScriptServices: script.NewScriptServices(adapters.ScriptImplementation, adapters.BookImplementation, adapters.AiImplementation, adapters.EventImplementation, adapters.QueueImplementation, adapters.StorageImplementation, adapters.EnvironmentVariables),
		EventServices:  event.NewEventService(adapters.EventImplementation),
	}
}
