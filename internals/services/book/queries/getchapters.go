package queries

import "github.com/Pr3c10us/absolutego/internals/domains/book"

type GetBooks struct {
	bookImplementation book.Interface
}

func (service *GetBooks) Handle(title string) ([]book.Book, error) {
	return service.bookImplementation.GetBooks(title)
}

func NewGetBooks(bookImplementation book.Interface) *GetBooks {
	return &GetBooks{
		bookImplementation,
	}
}
