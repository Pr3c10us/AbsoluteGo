package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/vab"
)

type QueueVAB struct {
	vabImplementation    vab.Interface
	scriptImplementation script.Interface
	bookImplementation   book.Interface
	eventImplementation  event.Interface
	queueImplementation  queue.Interface
}

func (service *QueueVAB) Handle(scriptId int64, name string) error {
	scr, err := service.scriptImplementation.GetScript(scriptId)
	if err != nil {
		return err
	}
	if scr == nil {
		return errors.New("script does not exist")
	}

	b, err := service.bookImplementation.GetBook(scr.BookId)
	if err != nil {
		return err
	}
	if b == nil {
		return errors.New("book does not exist")
	}

	splits, err := service.scriptImplementation.GetSplits(scr.Id)
	if err != nil {
		return err
	}
	if len(splits) < 1 {
		return errors.New("no splits for script")
	}

	eventId, err := service.eventImplementation.Create(event.Event{
		Status:      event.StatusEnqueue,
		Operation:   event.OpMergeVideo,
		Description: fmt.Sprintf("merging available videos for script %s", scr.Name),
		BookId:      b.Id,
	})
	if err != nil {
		return err
	}

	generateScriptParameters := CreateVABParameter{scr.Id, name}

	var dataByte []byte
	dataByte, err = json.Marshal(generateScriptParameters)
	if err != nil {
		return err
	}

	qMsg := queue.Message{
		EventId: eventId,
		Data:    dataByte,
	}

	err = service.queueImplementation.Publish(&queue.MessageParams{
		Queue:   queue.QueueMergeVideo,
		Message: qMsg,
	})
	if err != nil {
		return err
	}

	return nil

}

func NewQueueVAB(vabImplementation vab.Interface, scriptImplementation script.Interface, bookImplementation book.Interface, eventImplementation event.Interface, queueImplementation queue.Interface) *QueueVAB {
	return &QueueVAB{
		vabImplementation, scriptImplementation, bookImplementation, eventImplementation, queueImplementation,
	}
}
