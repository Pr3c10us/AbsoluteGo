package commands

import (
	"encoding/json"
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type CreateScript struct {
	eventImplementation  event.Interface
	bookImplementation   book.Interface
	queueImplementation  queue.Interface
	scriptImplementation script.Interface
}
type CreateScriptParameters struct {
	BookId          int64
	Chapters        []int
	Name            string
	PreviousScripts []int64
}

func (s *CreateScript) Handle(parameters CreateScriptParameters) error {
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

	scriptId, err := s.scriptImplementation.CreateScript(&script.Script{
		Name:     parameters.Name,
		BookId:   b.Id,
		Chapters: parameters.Chapters,
	})
	if err != nil {
		return err
	}

	eventId, err := s.eventImplementation.Create(event.Event{
		Status:    event.StatusEnqueue,
		Operation: event.OpGenScript,
		ScriptId:  scriptId,
	})

	addChapterParameter := GenerateScriptParameters{
		ScriptId:        scriptId,
		PreviousScripts: parameters.PreviousScripts,
	}

	var dataByte []byte
	dataByte, err = json.Marshal(addChapterParameter)
	if err != nil {
		return err
	}

	qMsg := queue.Message{
		EventId: eventId,
		Data:    dataByte,
	}

	err = s.queueImplementation.Publish(&queue.MessageParams{
		Queue:   queue.QueueGenScript,
		Message: qMsg,
	})
	if err != nil {
		return err
	}

	return nil
}

func NewCreateScript(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface, scriptImplementation script.Interface) *CreateScript {
	return &CreateScript{
		eventImplementation, bookImplementation, queueImplementation, scriptImplementation,
	}
}
