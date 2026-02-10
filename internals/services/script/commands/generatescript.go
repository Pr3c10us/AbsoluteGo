package commands

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/prompts"
	"github.com/Pr3c10us/absolutego/packages/utils"
)

type GenerateScript struct {
	script script.Interface
	book   book.Interface
	ai     ai.Interface
}

type Parameters struct {
	BookId          int64
	Name            string
	Chapters        []int
	PreviousScripts []int64
}

func (s *GenerateScript) Handle(parameters Parameters) (string, int64, error) {
	b, err := s.book.GetBook(parameters.BookId)
	if err != nil {
		return "", 0, err
	}
	if b == nil {
		return "", 0, appError.BadRequest(errors.New("book does not exist"))
	}

	fetchedChapters, err := s.book.GetChapters(b.Id, parameters.Chapters)
	if err != nil {
		return "", 0, err
	}
	if len(fetchedChapters) < 1 {
		return "", 0, appError.BadRequest(errors.New("chapters does not exist"))
	}

	var chapterIds []int64
	for _, chapter := range fetchedChapters {
		chapterIds = append(chapterIds, chapter.Id)
	}
	
	uploadedFiles, err := getUploads(chapterIds, s.book, s.ai)
	if err != nil {
		return "", 0, err
	}

	previousScripts, err := s.script.GetScripts(b.Id, "", parameters.PreviousScripts)
	if err != nil {
		return "", 0, err
	}
	concatenatedScript := s.concatScripts(previousScripts, b.Title)
	scriptPrompt := prompts.ScriptPrompt(b.Title, parameters.Chapters, &concatenatedScript)

	scriptResponse, err := s.ai.GenerateText(scriptPrompt, false, uploadedFiles)
	if err != nil {
		return "", 0, err
	}

	scriptId, err := s.script.CreateScript(&script.Script{
		Name:     parameters.Name,
		Content:  &scriptResponse.Response,
		BookId:   b.Id,
		Chapters: parameters.Chapters,
	})

	return scriptResponse.Response, scriptId, err
}

func getUploads(chapterIds []int64, bookImplementation book.Interface, aiImplementation ai.Interface) ([]ai.UploadedFile, error) {
	const maxWorkers = 5

	pages, err := bookImplementation.GetPages(chapterIds, false)
	if err != nil {
		return nil, err
	}

	uploads := make([]ai.UploadedFile, len(pages))

	type downloadJob struct {
		page  book.Page
		index int
	}

	var jobs []downloadJob
	for i, page := range pages {
		if time.Since(page.UpdatedAt) < 24*time.Hour {
			uploads[i] = ai.UploadedFile{
				URI:      *page.LLMURL,
				MIMEType: *page.MIME,
			}
			continue
		}
		jobs = append(jobs, downloadJob{page: page, index: i})
	}

	if len(jobs) == 0 {
		return uploads, nil
	}

	type downloadResult struct {
		tmpPath  string
		mimeType string
		index    int
	}

	results := make([]downloadResult, len(jobs))

	err = utils.RunWorkerPool(jobs, maxWorkers, func(j downloadJob) error {
		jobIndex := -1
		for i, job := range jobs {
			if job.index == j.index {
				jobIndex = i
				break
			}
		}

		tmpPath, err := downloadPage(*j.page.URL)
		if err != nil {
			return err
		}

		results[jobIndex] = downloadResult{
			tmpPath:  tmpPath,
			mimeType: *j.page.MIME,
			index:    j.index,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	files := make([]ai.File, len(results))
	for i, r := range results {
		files[i] = ai.File{
			Path:     r.tmpPath,
			MIMEType: r.mimeType,
		}
	}

	uploadedFiles, err := aiImplementation.UploadFiles(files)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, r := range results {
			os.Remove(r.tmpPath)
		}
	}()

	for i, r := range results {
		newPage := pages[r.index]
		newPage.LLMURL = &uploadedFiles[i].URI
		err = bookImplementation.UpdatePage(newPage.Id, &newPage)
		if err != nil {
			return nil, err
		}
		uploads[r.index] = ai.UploadedFile{
			URI:      uploadedFiles[i].URI,
			MIMEType: *newPage.MIME,
		}
	}

	return uploads, nil
}

func downloadPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpDir, err := utils.GetDirectory("tmp")
	if err != nil {
		return "", err
	}

	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%d-%d", time.Now().UnixMilli(), rand.Intn(1000000000)))

	out, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	return tmpFile, nil
}

func (s *GenerateScript) concatScripts(scripts []script.Script, bookTitle string) string {
	var sb strings.Builder
	for i, s := range scripts {
		if s.Content == nil {
			continue
		}

		chapters := make([]string, len(s.Chapters))
		for j, ch := range s.Chapters {
			chapters[j] = fmt.Sprintf("Chapter %d", ch)
		}

		var chapterStr string
		switch len(chapters) {
		case 0:
			chapterStr = ""
		case 1:
			chapterStr = chapters[0]
		case 2:
			chapterStr = chapters[0] + " and " + chapters[1]
		default:
			chapterStr = strings.Join(chapters[:len(chapters)-1], ", ") + ", and " + chapters[len(chapters)-1]
		}

		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("# %s Recap Script for %s\n%s", bookTitle, chapterStr, *s.Content))
	}
	return sb.String()
}

func NewGenerateScript(script script.Interface, book book.Interface, ai ai.Interface) *GenerateScript {
	return &GenerateScript{
		script, book, ai,
	}
}
