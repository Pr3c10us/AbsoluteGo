package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type CreateScript struct {
	eventImplementation event.Interface
	bookImplementation  book.Interface
	queueImplementation queue.Interface
}

func (s *CreateScript) Handle(parameters GenerateScriptParameters) error {
	b, err := s.bookImplementation.GetBook(parameters.BookId)
	if err != nil {
		return err
	}
	if b == nil {
		return appError.BadRequest(errors.New("book does not exist"))
	}

	fetchedChapters, err := s.bookImplementation.GetChapters(b.Id, parameters.Chapters)
	if err != nil {
		return err
	}
	if len(fetchedChapters) < 1 {
		return appError.BadRequest(errors.New("chapters does not exist"))
	}

}

func NewCreateScript(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface) *CreateScript {
	return &CreateScript{
		eventImplementation, bookImplementation, queueImplementation,
	}
}
