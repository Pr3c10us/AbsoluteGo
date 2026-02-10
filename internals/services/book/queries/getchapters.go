package queries

import "github.com/Pr3c10us/absolutego/internals/domains/book"

type GetChapters struct {
	bookImplementation book.Interface
}

func (service *GetChapters) Handle(bookId int64, number []int) ([]book.Chapter, error) {
	return service.bookImplementation.GetChapters(bookId, number)
}

func NewGetChapters(bookImplementation book.Interface) *GetChapters {
	return &GetChapters{
		bookImplementation,
	}
}
