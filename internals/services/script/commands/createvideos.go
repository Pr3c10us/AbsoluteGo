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

type CreateVideos struct {
	eventImplementation  event.Interface
	bookImplementation   book.Interface
	queueImplementation  queue.Interface
	scriptImplementation script.Interface
}

func (service *CreateVideos) Handle(scriptId int64) error {
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

	splits, err := service.scriptImplementation.GetSplits(scr.Id)
	if err != nil {
		return err
	}
	if len(splits) < 1 {
		return appError.BadRequest(errors.New("no splits for script"))
	}

	eventId, err := service.eventImplementation.Create(event.Event{
		Status:      event.StatusEnqueue,
		Operation:   event.OpGenVideo,
		Description: fmt.Sprintf("Video generation for script %q in %q with %d split(s)", scr.Name, b.Title, len(splits)),
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
		Queue:   queue.QueueGenVideos,
		Message: qMsg,
	})
	if err != nil {
		return err
	}

	return nil
}

func NewCreateVideos(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface, scriptImplementation script.Interface) *CreateVideos {
	return &CreateVideos{
		eventImplementation, bookImplementation, queueImplementation, scriptImplementation,
	}
}
