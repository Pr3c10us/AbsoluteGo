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

	fetchedChapters, err := s.bookImplementation.GetChapters(b.Id, parameters.Chapters, 0, 0)
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
		BookId:      b.Id,
	})
	if err != nil {
		return err
	}

	generateScriptParameters := GenerateScriptParameters{
		BookId:          b.Id,
		Name:            parameters.Name,
		Chapters:        parameters.Chapters,
		PreviousScripts: parameters.PreviousScripts,
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

	parts := make([]string, len(sorted))
	for i, ch := range sorted {
		parts[i] = fmt.Sprintf("%d", ch)
	}
	return fmt.Sprintf("Create script from %q ch. %s", bookName, formatChapterList(sorted))
}

func formatChapterList(sorted []int) string {
	var parts []string
	i := 0
	for i < len(sorted) {
		start := sorted[i]
		end := start
		for i+1 < len(sorted) && sorted[i+1] == end+1 {
			i++
			end = sorted[i]
		}
		if start == end {
			parts = append(parts, fmt.Sprintf("%d", start))
		} else {
			parts = append(parts, fmt.Sprintf("%d-%d", start, end))
		}
		i++
	}
	return strings.Join(parts, ", ")
}
func NewCreateScript(eventImplementation event.Interface, bookImplementation book.Interface, queueImplementation queue.Interface, scriptImplementation script.Interface) *CreateScript {
	return &CreateScript{
		eventImplementation, bookImplementation, queueImplementation, scriptImplementation,
	}
}
