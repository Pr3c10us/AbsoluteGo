package queries

import "github.com/Pr3c10us/absolutego/internals/domains/book"

type GetPanels struct {
	bookImplementation book.Interface
}

func (service *GetPanels) Handle(pageId int64) ([]book.Panel, error) {
	return service.bookImplementation.GetPanels(pageId)
}

func NewGetPanels(bookImplementation book.Interface) *GetPanels {
	return &GetPanels{
		bookImplementation,
	}
}
