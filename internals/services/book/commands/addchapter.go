package commands

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
)

type AddChapter struct {
	book    book.Interface
	storage storage.Interface
	ai      ai.Interface
	env     *configs.EnvironmentVariables
}

type Parameter struct {
	File    string
	Chapter int
	BookId  int64
}

// uploadTracker collects all uploaded URLs so they can be rolled back on failure.
// uploadTracker collects all uploaded URLs so they can be rolled back on failure.
type uploadTracker struct {
	mu   sync.Mutex
	urls []string
}

func (t *uploadTracker) Add(urls ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.urls = append(t.urls, urls...)
}

func (t *uploadTracker) All() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	dst := make([]string, len(t.urls))
	copy(dst, t.urls)
	return dst
}

// runWorkerPool processes items concurrently with a bounded number of workers.
// It stops early on the first error encountered.
func runWorkerPool[T any](items []T, maxWorkers int, fn func(T) error) error {
	if len(items) == 0 {
		return nil
	}
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	if maxWorkers > len(items) {
		maxWorkers = len(items)
	}

	var (
		wg       sync.WaitGroup
		once     sync.Once
		firstErr error
		work     = make(chan T, len(items))
	)

	for _, item := range items {
		work <- item
	}
	close(work)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range work {
				if err := fn(item); err != nil {
					once.Do(func() { firstErr = err })
					return
				}
			}
		}()
	}

	wg.Wait()
	return firstErr
}

func (s *AddChapter) Handle(p Parameter) error {
	defer os.Remove(p.File)

	b, err := s.book.GetBook(p.BookId)
	if err != nil {
		return err
	}
	if b == nil {
		return appError.BadRequest(errors.New("book does not exist"))
	}

	chapters, _ := s.book.GetChapters(b.Id, p.Chapter)
	for _, ch := range chapters {
		if err = s.deleteChapter(ch.Id); err != nil {
			return err
		}
	}

	chapterId, err := s.book.CreateChapter(p.BookId, p.Chapter, "")
	if err != nil {
		return err
	}

	outputDir, err := utils.GetDirectory("books")
	if err != nil {
		return err
	}
	defer os.RemoveAll(outputDir)

	if err = s.processFile(outputDir, p.File); err != nil {
		return err
	}

	pagePaths, err := utils.GetFilesInDir(outputDir)
	if err != nil {
		return err
	}

	// Track every uploaded URL for rollback on failure.
	tracker := &uploadTracker{}
	rollback := func() { s.storage.DeleteMany(tracker.All()) }

	processedPages, cover, err := s.processPages(pagePaths, chapterId, tracker)
	if err != nil {
		rollback()
		return err
	}

	if cover == nil || cover.URL == nil {
		rollback()
		return appError.BadRequest(errors.New("no cover image found in the uploaded file"))
	}

	if err = s.book.UpdateChapter(chapterId, 0, *cover.URL); err != nil {
		rollback()
		return err
	}

	pages, err := s.book.CreateManyPage(processedPages)
	if err != nil {
		rollback()
		return err
	}

	subDir, err := utils.GetSubDirs(outputDir)
	if err != nil {
		rollback()
		return err
	}

	processedPanels, err := s.processPanels(subDir, pages, tracker)
	if err != nil {
		rollback()
		return err
	}

	if _, err = s.book.CreateManyPanel(processedPanels); err != nil {
		rollback()
		return err
	}

	return nil
}

func (s *AddChapter) processFile(outputDir, filePath string) error {
	format, err := utils.GetComicFormat(filePath)
	if err != nil {
		return err
	}

	tempDir, err := utils.GetDirectory("tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err = utils.ExtractComicToDir(filePath, format, tempDir); err != nil {
		return err
	}

	images, err := utils.SortImages(tempDir, outputDir)
	if err != nil {
		return err
	}
	if len(images) < 1 {
		return errors.New("no images extracted")
	}

	if _, err = utils.GenerateBlurCover(outputDir, images[0]); err != nil {
		return err
	}

	const maxWorkers = 100

	err = runWorkerPool(images, maxWorkers, func(path string) error {
		defer os.Remove(path)
		utils.DetectAndExtractPanels(path)
		return nil
	})
	if err != nil {
		return err
	}

	type overlayJob struct {
		path    string
		pageNum int
	}
	var jobs []overlayJob
	for i, img := range images {
		dir := filepath.Dir(img)
		name := strings.TrimSuffix(filepath.Base(img), filepath.Ext(img))
		overlayPath := filepath.Join(dir, name+".png")
		if _, err := os.Stat(overlayPath); err == nil {
			jobs = append(jobs, overlayJob{path: overlayPath, pageNum: i + 1})
		}
	}

	return runWorkerPool(jobs, maxWorkers, func(j overlayJob) error {
		return utils.AddPageNumberToOverlay(j.path, j.pageNum)
	})
}

func (s *AddChapter) deleteChapter(chapterId int64) error {
	ch, err := s.book.GetChapter(chapterId)
	if err != nil || ch == nil {
		return err
	}

	pages, err := s.book.GetPages(ch.Id)
	if err != nil || pages == nil {
		return err
	}

	var urls []string
	for _, page := range pages {
		urls = append(urls, *page.URL)

		panels, err := s.book.GetPanels(page.Id)
		if err != nil {
			return err
		}
		for _, panel := range panels {
			urls = append(urls, *panel.URL)
		}
		if err = s.book.DeletePanels(page.Id); err != nil {
			return err
		}
	}

	if err = s.book.DeletePages(ch.Id); err != nil {
		return err
	}
	if err = s.book.DeleteChapter(ch.Id); err != nil {
		return err
	}

	s.storage.DeleteMany(urls)
	return nil
}

func (s *AddChapter) processPages(pagePaths []string, chapterId int64, tracker *uploadTracker) ([]book.Page, *book.Page, error) {
	var (
		osFiles     []*os.File
		aiFiles     []ai.File
		pageNumbers []int
		cover       *book.Page
	)
	defer func() {
		for _, f := range osFiles {
			f.Close()
		}
	}()

	for _, p := range pagePaths {
		filename := filepath.Base(p)
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		ext := filepath.Ext(filename)

		f, err := os.Open(p)
		if err != nil {
			return nil, nil, err
		}

		if strings.HasPrefix(name, "cover") {
			defer f.Close()
			url, err := s.storage.UploadFile(s.env.Buckets.PageBucket, f)
			if err != nil {
				return nil, nil, err
			}
			tracker.Add(url)
			cover = &book.Page{URL: &url}
			continue
		}

		pageNum, err := strconv.Atoi(name)
		if err != nil {
			f.Close()
			continue
		}

		osFiles = append(osFiles, f)
		aiFiles = append(aiFiles, ai.File{Path: p, MIMEType: mime.TypeByExtension(ext)})
		pageNumbers = append(pageNumbers, pageNum)
	}

	var (
		llmFiles []ai.UploadedFile
		uploaded []storage.UploadResult
		llmErr   error
		wg       sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		llmFiles, llmErr = s.ai.UploadFiles(aiFiles)
	}()
	go func() {
		defer wg.Done()
		uploaded = s.storage.UploadMany(s.env.Buckets.PageBucket, osFiles)
	}()
	wg.Wait()

	if llmErr != nil {
		return nil, nil, llmErr
	}

	for _, f := range uploaded {
		if f.Err != nil {
			return nil, nil, errors.New("failed to upload all pages")
		}
		tracker.Add(f.URL)
	}

	pages := make([]book.Page, len(pageNumbers))
	for i := range pageNumbers {
		pages[i] = book.Page{
			PageNumber: pageNumbers[i],
			LLMURL:     &llmFiles[i].URI,
			MIME:       &llmFiles[i].MIMEType,
			URL:        &uploaded[i].URL,
			ChapterId:  chapterId,
		}
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].PageNumber < pages[j].PageNumber
	})

	return pages, cover, nil
}

func (s *AddChapter) processPanels(panelDir []string, pages []book.Page, tracker *uploadTracker) ([]book.Panel, error) {
	const maxWorkers = 100

	pageMap := make(map[int]int64, len(pages))
	for _, page := range pages {
		pageMap[page.PageNumber] = page.Id
	}

	var (
		mu        sync.Mutex
		allPanels []book.Panel
	)

	err := runWorkerPool(panelDir, maxWorkers, func(dir string) error {
		defer os.RemoveAll(dir)

		dirName := filepath.Base(dir)
		pageNumber, err := strconv.Atoi(dirName)
		if err != nil {
			return err
		}

		pageId, ok := pageMap[pageNumber]
		if !ok {
			return fmt.Errorf("page %d not found", pageNumber)
		}

		panels, err := s.processPanel(dir, pageId, tracker)
		if err != nil {
			return err
		}

		mu.Lock()
		allPanels = append(allPanels, panels...)
		mu.Unlock()
		return nil
	})

	return allPanels, err
}

func (s *AddChapter) processPanel(panelDir string, pageId int64, tracker *uploadTracker) ([]book.Panel, error) {
	panelPaths, err := utils.GetFilesInDir(panelDir)
	if err != nil {
		return nil, err
	}

	var (
		osFiles      []*os.File
		panelNumbers []int
	)
	defer func() {
		for _, f := range osFiles {
			f.Close()
		}
	}()

	for _, p := range panelPaths {
		filename := filepath.Base(p)
		name := strings.TrimSuffix(filename, filepath.Ext(filename))

		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}

		panelNum, err := strconv.Atoi(name)
		if err != nil {
			f.Close()
			continue
		}

		osFiles = append(osFiles, f)
		panelNumbers = append(panelNumbers, panelNum)
	}

	uploaded := s.storage.UploadMany(s.env.Buckets.PanelBucket, osFiles)

	for _, f := range uploaded {
		if f.Err != nil {
			return nil, errors.New("failed to upload all panels")
		}
		tracker.Add(f.URL)
	}

	panels := make([]book.Panel, len(panelNumbers))
	for i := range panelNumbers {
		panels[i] = book.Panel{
			PanelNumber: panelNumbers[i],
			URL:         &uploaded[i].URL,
			PageId:      pageId,
		}
	}

	sort.Slice(panels, func(i, j int) bool {
		return panels[i].PanelNumber < panels[j].PanelNumber
	})

	return panels, nil
}

func NewAddChapter(b book.Interface, st storage.Interface, a ai.Interface, env *configs.EnvironmentVariables) *AddChapter {
	return &AddChapter{b, st, a, env}
}
