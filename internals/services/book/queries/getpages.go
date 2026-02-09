package queries

import "github.com/Pr3c10us/absolutego/internals/domains/book"

type GetPages struct {
	bookImplementation book.Interface
}

func (service *GetPages) Handle(chapterId int64) ([]book.Page, error) {
	return service.bookImplementation.GetPages(chapterId)
}

func NewGetPages(bookImplementation book.Interface) *GetPages {
	return &GetPages{
		bookImplementation,
	}
}
