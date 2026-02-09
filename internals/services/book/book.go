package book

import (
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/internals/services/book/commands"
	"github.com/Pr3c10us/absolutego/packages/configs"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
	AddBook *commands.AddBook
}

type Queries struct {
}

func NewBookServices(bookImplementation book.Interface, storageImplementation storage.Interface, environmentVariables *configs.EnvironmentVariables) Services {
	return Services{
		Commands: Commands{
			AddBook: commands.NewAddBook(bookImplementation, storageImplementation, environmentVariables),
		},
		Queries: Queries{},
	}
}
