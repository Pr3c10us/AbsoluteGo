package services

import (
	"github.com/Pr3c10us/absolutego/internals/adapters"
	"github.com/Pr3c10us/absolutego/internals/services/book"
)

type Services struct {
	BookServices book.Services
}

func NewServices(adapters *adapters.Adapters) *Services {
	return &Services{
		BookServices: book.NewBookServices(adapters.BookImplementation, adapters.StorageRepository, adapters.AiImplementation, adapters.EnvironmentVariables),
	}
}
