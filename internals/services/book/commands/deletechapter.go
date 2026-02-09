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

type DeleteChapterParameter struct {
	ChapterId int64
}

func (s *DeleteChapter) Handle(p DeleteChapterParameter) error {
	ch, err := s.bookImplementation.GetChapter(p.ChapterId)
	if err != nil {
		return err
	}
	if ch == nil {
		return appError.BadRequest(errors.New("chapter does not exist"))
	}

	// Get all pages for this chapter
	pages, err := s.bookImplementation.GetPages(ch.Id)
	if err != nil {
		return err
	}

	// Collect all URLs to delete from storage
	var urls []string

	for _, page := range pages {
		if page.URL != nil {
			urls = append(urls, *page.URL)
		}

		// Get all panels for this page
		panels, err := s.bookImplementation.GetPanels(page.Id)
		if err != nil {
			return err
		}

		for _, panel := range panels {
			if panel.URL != nil {
				urls = append(urls, *panel.URL)
			}
		}

		// Delete panels for this page from database
		if err = s.bookImplementation.DeletePanels(page.Id); err != nil {
			return err
		}
	}

	// Delete pages for this chapter from database
	if err = s.bookImplementation.DeletePages(ch.Id); err != nil {
		return err
	}

	// Add chapter blur URL if exists
	if ch.BlurURL != "" {
		urls = append(urls, ch.BlurURL)
	}

	// Delete the chapter from database
	if err = s.bookImplementation.DeleteChapter(ch.Id); err != nil {
		return err
	}

	// Delete all URLs from storage
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
