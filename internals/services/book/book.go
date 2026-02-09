package book

import (
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/internals/services/book/commands"
	"github.com/Pr3c10us/absolutego/internals/services/book/queries"
	"github.com/Pr3c10us/absolutego/packages/configs"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
	AddChapter    *commands.AddChapter
	AddBook       *commands.AddBook
	DeleteBook    *commands.DeleteBook
	DeleteChapter *commands.DeleteChapter
}

type Queries struct {
	GetBooks    *queries.GetBooks
	GetChapters *queries.GetChapters
	GetPages    *queries.GetPages
	GetPanels   *queries.GetPanels
}

func NewBookServices(bookImplementation book.Interface, storageImplementation storage.Interface, aiImplementation ai.Interface, environmentVariables *configs.EnvironmentVariables) Services {
	return Services{
		Commands: Commands{
			AddChapter:    commands.NewAddChapter(bookImplementation, storageImplementation, aiImplementation, environmentVariables),
			AddBook:       commands.NewAddBook(bookImplementation),
			DeleteBook:    commands.NewDeleteBook(bookImplementation, storageImplementation),
			DeleteChapter: commands.NewDeleteChapter(bookImplementation, storageImplementation),
		},
		Queries: Queries{
			GetBooks:    queries.NewGetBooks(bookImplementation),
			GetChapters: queries.NewGetChapters(bookImplementation),
			GetPages:    queries.NewGetPages(bookImplementation),
			GetPanels:   queries.NewGetPanels(bookImplementation),
		},
	}
}
