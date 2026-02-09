package commands

import "github.com/Pr3c10us/absolutego/internals/domains/book"

type AddBook struct {
	bookImplementation book.Interface
}

func (service *CreateBook) Handle(title string) error {
	_, err := service.bookImplementation.CreateBook(title)
	return err
}

func NewCreateBook(bookImplementation book.Interface) *CreateBook {
	return &CreateBook{
		bookImplementation,
	}
}
