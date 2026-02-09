package commands

import (
	"errors"

	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type DeleteBook struct {
	bookImplementation    book.Interface
	storageImplementation storage.Interface
}

type DeleteBookParameter struct {
	BookId int64
}

func (s *DeleteBook) Handle(p DeleteBookParameter) error {
	b, err := s.bookImplementation.GetBook(p.BookId)
	if err != nil {
		return err
	}
	if b == nil {
		return appError.BadRequest(errors.New("book does not exist"))
	}

	// Get all chapters for this book (passing 0 for number to get all)
	chapters, err := s.bookImplementation.GetChapters(b.Id, 0)
	if err != nil {
		return err
	}

	// Collect all URLs to delete from storage
	var urls []string

	for _, ch := range chapters {
		// Get all pages for this chapter
		pages, err := s.bookImplementation.GetPages(ch.Id)
		if err != nil {
			return err
		}

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
	}

	// Delete the book from database
	if err = s.bookImplementation.DeleteBook(b.Id); err != nil {
		return err
	}

	// Delete all URLs from storage
	if len(urls) > 0 {
		s.storageImplementation.DeleteMany(urls)
	}

	return nil
}

func NewDeleteBook(bookImplementation book.Interface, storageImplementation storage.Interface) *DeleteBook {
	return &DeleteBook{
		bookImplementation:    bookImplementation,
		storageImplementation: storageImplementation,
	}
}
