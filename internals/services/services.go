package services

import (
	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/services/book"
	"github.com/Pr3c10us/absolutego/internals/services/script"
)

type Services struct {
	BookServices   book.Services
	ScriptServices script.Services
}

func NewServices(adapters *adapters.Adapters) *Services {
	return &Services{
		BookServices:   book.NewBookServices(adapters.BookImplementation, adapters.StorageRepository, adapters.AiImplementation, adapters.EnvironmentVariables, adapters.ScriptImplementation),
		ScriptServices: script.NewScriptServices(adapters.ScriptImplementation, adapters.BookImplementation, adapters.AiImplementation),
	}
}
