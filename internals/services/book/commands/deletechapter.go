package commands

import (
	"errors"

	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type DeleteChapter struct {
	bookImplementation    book.Interface
	storageImplementation storage.Interface
}

func (s *DeleteChapter) Handle(chapterId int64) error {
	ch, err := s.bookImplementation.GetChapter(chapterId)
	if err != nil {
		return err
	}
	if ch == nil {
		return appError.BadRequest(errors.New("chapter does not exist"))
	}

	pages, err := s.bookImplementation.GetPages(ch.Id)
	if err != nil {
		return err
	}

	var urls []string

	for _, page := range pages {
		if page.URL != nil {
			urls = append(urls, *page.URL)
		}

		var panels []book.Panel
		panels, err = s.bookImplementation.GetPanels(page.Id)
		if err != nil {
			return err
		}

		for _, panel := range panels {
			if panel.URL != nil {
				urls = append(urls, *panel.URL)
			}
		}

		if err = s.bookImplementation.DeletePanels(page.Id); err != nil {
			return err
		}
	}

	if err = s.bookImplementation.DeletePages(ch.Id); err != nil {
		return err
	}

	if ch.BlurURL != "" {
		urls = append(urls, ch.BlurURL)
	}

	if err = s.bookImplementation.DeleteChapter(ch.Id); err != nil {
		return err
	}

	if len(urls) > 0 {
		s.storageImplementation.DeleteMany(urls)
	}

	return nil
}

func NewDeleteChapter(bookImplementation book.Interface, storageImplementation storage.Interface) *DeleteChapter {
	return &DeleteChapter{
		bookImplementation:    bookImplementation,
		storageImplementation: storageImplementation,
	}
}
