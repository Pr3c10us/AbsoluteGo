package commands

import (
	"errors"

	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type DeleteBook struct {
	bookImplementation    book.Interface
	storageImplementation storage.Interface
	deleteChapter         *DeleteChapter
}

func (s *DeleteBook) Handle(bookId int64) error {
	b, err := s.bookImplementation.GetBook(bookId)
	if err != nil {
		return err
	}
	if b == nil {
		return appError.BadRequest(errors.New("book does not exist"))
	}

	chapters, err := s.bookImplementation.GetChapters(b.Id, nil)
	if err != nil {
		return err
	}

	for _, chapter := range chapters {
		err = s.deleteChapter.Handle(chapter.Id)
		if err != nil {
			return err
		}
	}

	err = s.bookImplementation.DeleteBook(bookId)
	return err
}

func NewDeleteBook(bookImplementation book.Interface, storageImplementation storage.Interface) *DeleteBook {
	return &DeleteBook{
		bookImplementation:    bookImplementation,
		storageImplementation: storageImplementation,
		deleteChapter:         NewDeleteChapter(bookImplementation, storageImplementation),
	}
}
