package queries

import "github.com/Pr3c10us/absolutego/internals/domains/book"

type GetPages struct {
	bookImplementation book.Interface
}

func (service *GetPages) Handle(chapterIds []int64, withPanels bool) ([]book.Page, error) {
	return service.bookImplementation.GetPages(chapterIds, withPanels)
}

func NewGetPages(bookImplementation book.Interface) *GetPages {
	return &GetPages{
		bookImplementation,
	}
}
