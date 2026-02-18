package queries

import "github.com/Pr3c10us/absolutego/internals/domains/book"

type GetBooks struct {
	bookImplementation book.Interface
}

func (service *GetBooks) Handle(title string, page, limit int) ([]book.Book, error) {
	return service.bookImplementation.GetBooks(title, page, limit)
}

func NewGetBooks(bookImplementation book.Interface) *GetBooks {
	return &GetBooks{
		bookImplementation,
	}
}
