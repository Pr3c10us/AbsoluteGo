package commands

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type CreateSplits struct {
	eventImplementation  event.Interface
	bookImplementation   book.Interface
	queueImplementation  queue.Interface
	scriptImplementation script.Interface
}

func (service *CreateSplits) Handle(scriptId int64) error {
	scr, err := service.scriptImplementation.GetScript(scriptId)
	if err != nil {
		return err
	}
	if scr == nil {
		return appError.BadRequest(errors.New("script does not exist"))
	}

	b, err := service.bookImplementation.GetBook(scr.BookId)
	if err != nil {
		return err
	}
	if b == nil {
		return appError.BadRequest(errors.New("book does not exist"))
	}

	eventId, err := service.eventImplementation.Create(event.Event{
		Status:      event.StatusEnqueue,
		Operation:   event.OpGenScriptSplit,
		Description: fmt.Sprintf("splitting %s script for %s", scr.Name, b.Title),
		BookId:      b.Id,
	})
	if err != nil {
		return err
	}

	dataByte := make([]byte, 8)
	binary.BigEndian.PutUint64(dataByte, uint64(scr.Id))

	qMsg := queue.Message{
		EventId: eventId,
		Data:    dataByte,
	}

	err = service.queueImplementation.Publish(&queue.MessageParams{
		Queue:   queue.QueueGenScript,
		Message: qMsg,
	})
	if err != nil {
		return err
	}

	return nil

}

func NewCreateSplits(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface, scriptImplementation script.Interface) *CreateSplits {
	return &CreateSplits{
		eventImplementation, bookImplementation, queueImplementation, scriptImplementation,
	}
}
