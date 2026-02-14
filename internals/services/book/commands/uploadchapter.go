package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"os"
)

type UploadChapterParameter struct {
	File    string
	Chapter int
	BookId  int64
}

type UploadChapter struct {
	book          book.Interface
	storage       storage.Interface
	env           *configs.EnvironmentVariables
	event         event.Interface
	queue         queue.Interface
	deleteChapter *DeleteChapter
}

func (s *UploadChapter) Handle(p UploadChapterParameter) error {
	defer os.Remove(p.File)

	b, err := s.book.GetBook(p.BookId)
	if err != nil {
		return err
	}
	if b == nil {
		return appError.BadRequest(errors.New("book does not exist"))
	}

	osFile, err := os.Open(p.File)
	if err != nil {
		return err
	}
	defer osFile.Close()

	url, err := s.storage.UploadFile(s.env.Buckets.ComicBucket, osFile)
	if err != nil {
		return errors.New("failed to upload chapter")
	}

	eventId, err := s.event.Create(event.Event{
		Status:      event.StatusEnqueue,
		Operation:   event.OpAddChapter,
		Description: fmt.Sprintf("adding chapter %d to %s", p.Chapter, b.Title),
		BookId:      b.Id,
	})
	if err != nil {
		return err
	}

	addChapterParameter := AddChapterParameter{
		FileUrl: url,
		Chapter: p.Chapter,
		BookId:  b.Id,
	}

	var dataByte []byte
	dataByte, err = json.Marshal(addChapterParameter)
	if err != nil {
		_ = s.storage.DeleteFile(url)
		return err
	}

	qMsg := queue.Message{
		EventId: eventId,
		Data:    dataByte,
	}

	err = s.queue.Publish(&queue.MessageParams{
		Queue:   queue.QueueAddChapter,
		Message: qMsg,
	})
	if err != nil {
		_ = s.storage.DeleteFile(url)
		return err
	}

	return nil
}

func NewUploadChapter(b book.Interface, st storage.Interface, env *configs.EnvironmentVariables, e event.Interface, q queue.Interface, scriptImplementation script.Interface) *UploadChapter {
	return &UploadChapter{b, st, env, e, q, NewDeleteChapter(b, st, scriptImplementation)}
}
