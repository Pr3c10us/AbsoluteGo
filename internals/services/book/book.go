package book

import (
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/internals/domains/vab"
	"github.com/Pr3c10us/absolutego/internals/services/book/commands"
	"github.com/Pr3c10us/absolutego/internals/services/book/queries"
	"github.com/Pr3c10us/absolutego/packages/configs"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
	UploadChapter *commands.UploadChapter
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

func NewBookServices(
	bookImplementation book.Interface,
	storageImplementation storage.Interface,
	aiImplementation ai.Interface,
	environmentVariables *configs.EnvironmentVariables,
	scriptImplementation script.Interface,
	eventImplementation event.Interface,
	queueImplementation queue.Interface,
	vabImplementation vab.Interface,
) Services {
	return Services{
		Commands: Commands{
			UploadChapter: commands.NewUploadChapter(bookImplementation, storageImplementation, environmentVariables, eventImplementation, queueImplementation, scriptImplementation, vabImplementation),
			AddChapter:    commands.NewAddChapter(bookImplementation, storageImplementation, aiImplementation, environmentVariables, scriptImplementation, vabImplementation),
			AddBook:       commands.NewAddBook(bookImplementation),
			DeleteBook:    commands.NewDeleteBook(bookImplementation, storageImplementation, scriptImplementation, vabImplementation),
			DeleteChapter: commands.NewDeleteChapter(bookImplementation, storageImplementation, scriptImplementation, vabImplementation),
		},
		Queries: Queries{
			GetBooks:    queries.NewGetBooks(bookImplementation),
			GetChapters: queries.NewGetChapters(bookImplementation),
			GetPages:    queries.NewGetPages(bookImplementation),
			GetPanels:   queries.NewGetPanels(bookImplementation),
		},
	}
}
