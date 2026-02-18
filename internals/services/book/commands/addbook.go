package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type AddBook struct {
	bookImplementation book.Interface
}

func (service *AddBook) Handle(title string) error {
	books, err := service.bookImplementation.GetBooks(title, 1, 1)
	if err != nil {
		return err
	}
	if books == nil || len(books) < 1 {
		_, err = service.bookImplementation.CreateBook(title)
		return err
	} else {
		return appError.BadRequest(errors.New("book already exist"))
	}

}

func NewAddBook(bookImplementation book.Interface) *AddBook {
	return &AddBook{
		bookImplementation,
	}
}
