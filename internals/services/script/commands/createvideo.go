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

func (service *CreateAudio) Handle(id int64, voice ai.Voice, voiceStyle string) error {
	split, err := service.scriptImplementation.GetSplit(id)
	if err != nil {
		return err
	}
	if split == nil {
		return errors.New("split does not exist")
	}

	scr, err := service.scriptImplementation.GetScript(split.ScriptId)
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

	eventId, err := service.eventImplementation.Create(event.Event{
		Status:      event.StatusEnqueue,
		Operation:   event.OpGenAudio,
		Description: fmt.Sprintf("Audio generation for %s, script %q, split %d", b.Title, scr.Name, split.Id),
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
		Queue:   queue.QueueGenAudio,
		Message: qMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewCreateAudio(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface, scriptImplementation script.Interface) *CreateAudio {
	return &CreateAudio{
		eventImplementation, bookImplementation, queueImplementation, scriptImplementation,
	}
}
