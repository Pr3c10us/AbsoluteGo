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
	"sort"
	"strings"
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

	eventId, err := s.eventImplementation.Create(event.Event{
		Status:      event.StatusEnqueue,
		Operation:   event.OpGenScript,
		Description: buildOperation(b.Title, parameters.Chapters),
	})

	addChapterParameter := GenerateScriptParameters{
		BookId:          b.Id,
		Name:            parameters.Name,
		Chapters:        parameters.Chapters,
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

func buildOperation(bookName string, chapters []int) string {
	if len(chapters) == 0 {
		return fmt.Sprintf("Create script from %q", bookName)
	}

	sorted := make([]int, len(chapters))
	copy(sorted, chapters)
	sort.Ints(sorted)

	if isContiguousRange(sorted) {
		if len(sorted) == 1 {
			return fmt.Sprintf("Create script from %q ch. %d", bookName, sorted[0])
		}
		return fmt.Sprintf("Create script from %q ch. %d to %d", bookName, sorted[0], sorted[len(sorted)-1])
	}

	parts := make([]string, len(sorted))
	for i, ch := range sorted {
		parts[i] = fmt.Sprintf("%d", ch)
	}
	return fmt.Sprintf("Create script from %q ch. %s", bookName, strings.Join(parts, ", "))
}

func isContiguousRange(sorted []int) bool {
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[i-1]+1 {
			return false
		}
	}
	return true
}

func NewCreateScript(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface, scriptImplementation script.Interface) *CreateScript {
	return &CreateScript{
		eventImplementation, bookImplementation, queueImplementation, scriptImplementation,
	}
}
