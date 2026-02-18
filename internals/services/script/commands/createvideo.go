package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type CreateVideo struct {
	eventImplementation  event.Interface
	bookImplementation   book.Interface
	queueImplementation  queue.Interface
	scriptImplementation script.Interface
}

func (service *CreateVideo) Handle(parameter GenerateVideoParameter) error {
	split, err := service.scriptImplementation.GetSplit(parameter.Id)
	if err != nil {
		return err
	}
	if split == nil {
		return appError.BadRequest(errors.New("split does not exist"))
	}
	if split.AudioURL == nil || split.AudioDuration == nil {
		return appError.BadRequest(errors.New("generate split audio first"))
	}

	scr, err := service.scriptImplementation.GetScript(split.ScriptId)
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
		Operation:   event.OpGenVideo,
		Description: fmt.Sprintf("Video generation for %s, script %q, split %d", b.Title, scr.Name, split.Id),
		BookId:      b.Id,
	})
	if err != nil {
		return err
	}

	var dataByte []byte
	dataByte, err = json.Marshal(parameter)
	if err != nil {
		return err
	}

	qMsg := queue.Message{
		EventId: eventId,
		Data:    dataByte,
	}

	err = service.queueImplementation.Publish(&queue.MessageParams{
		Queue:   queue.QueueGenVideo,
		Message: qMsg,
	})
	if err != nil {
		return err
	}

	return nil
}

func NewCreateVideo(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface, scriptImplementation script.Interface) *CreateVideo {
	return &CreateVideo{
		eventImplementation, bookImplementation, queueImplementation, scriptImplementation,
	}
}
