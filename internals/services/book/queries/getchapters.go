package queries

import "github.com/Pr3c10us/absolutego/internals/domains/book"

type GetChapters struct {
	bookImplementation book.Interface
}

func (service *GetChapters) Handle(bookId int64, number []int, page, limit int) ([]book.Chapter, error) {
	if limit <= 0 {
		limit = 20
	}
	return service.bookImplementation.GetChapters(bookId, number, page, limit)
}

func NewGetChapters(bookImplementation book.Interface) *GetChapters {
	return &GetChapters{
		bookImplementation,
	}
}
