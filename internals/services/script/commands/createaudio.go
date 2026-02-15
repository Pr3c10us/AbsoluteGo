package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
)

type CreateAudio struct {
	eventImplementation  event.Interface
	bookImplementation   book.Interface
	queueImplementation  queue.Interface
	scriptImplementation script.Interface
}

func (service *CreateAudio) Handle(scriptId int64, voice ai.Voice, voiceStyle string) error {
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
		return errors.New("nno splits for script")
	}

	for _, split := range splits {
		eventId, err := service.eventImplementation.Create(event.Event{
			Status:      event.StatusEnqueue,
			Operation:   event.OpGenScript,
			Description: fmt.Sprintf("generate audio for split %d in script %s for %s", split.Id, scr.Name, b.Title),
			BookId:      b.Id,
		})
		if err != nil {
			return err
		}

		generateScriptParameters := AudioParameter{
			Id:         split.Id,
			Voice:      voice,
			VoiceStyle: voiceStyle,
		}

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
			Queue:   queue.QueueGenScript,
			Message: qMsg,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func NewCreateAudio(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface, scriptImplementation script.Interface) *CreateScript {
	return &CreateScript{
		eventImplementation, bookImplementation, queueImplementation, scriptImplementation,
	}
}
