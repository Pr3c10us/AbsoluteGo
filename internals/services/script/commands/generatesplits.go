package commands

import (
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/prompts"
)

type GenerateSplits struct {
	script script.Interface
	book   book.Interface
	ai     ai.Interface
}

func (s *GenerateSplits) Handle(scriptId int64) error {
	scr, err := s.script.GetScript(scriptId)
	if err != nil {
		return err
	}
	if scr == nil {
		return appError.BadRequest(errors.New("script does not exist"))
	}

	b, err := s.book.GetBook(scr.BookId)
	if err != nil {
		return err
	}
	if b == nil {
		return appError.BadRequest(errors.New("book does not exist"))
	}

	fetchedChapters, err := s.book.GetChapters(b.Id, scr.Chapters)
	if err != nil {
		return err
	}
	if len(fetchedChapters) < 1 {
		return appError.BadRequest(errors.New("chapters do not exist"))
	}

	chapterIds := make([]int64, len(fetchedChapters))
	for i, chapter := range fetchedChapters {
		chapterIds[i] = chapter.Id
	}

	pages, err := s.book.GetPages(chapterIds, true)
	if err != nil {
		return err
	}
	if len(pages) < 1 {
		return appError.BadRequest(errors.New("no pages across chapters"))
	}

	uploadedFiles, err := getUploads(chapterIds, s.book, s.ai)
	if err != nil {
		return err
	}

	splitResponse, err := s.ai.GenerateText(prompts.SplitScriptPrompt(*scr.Content), false, uploadedFiles)
	if err != nil {
		return err
	}

	splitResult, err := prompts.ParseSplitScriptResponse(splitResponse.Response)
	if err != nil {
		return err
	}

	pageMap := make(map[string]*book.Page, len(pages))
	for i := range pages {
		key := fmt.Sprintf("%d-%d", pages[i].Chapter.Number, pages[i].PageNumber)
		pageMap[key] = &pages[i]
	}

	type splitGroup struct {
		index   int
		content string
		effect  string
		panelId *int64
	}

	splitMap := make(map[string]*splitGroup)
	orderedKeys := make([]string, 0, len(splitResult))

	for i, result := range splitResult {
		key := fmt.Sprintf("%d-%d", result.Chapter, result.Page)
		p := pageMap[key]
		if p == nil || len(p.Panels) == 0 {
			continue
		}

		if existing, exists := splitMap[key]; exists {
			existing.content = existing.content + "\n" + result.Script
		} else {
			panelId := resolvePanelId(p.Panels, result.Panel)
			splitMap[key] = &splitGroup{
				index:   i,
				content: result.Script,
				effect:  result.Effect,
				panelId: panelId,
			}
			orderedKeys = append(orderedKeys, key)
		}
	}

	splits := make([]script.Split, 0, len(splitMap))
	for _, key := range orderedKeys {
		group := splitMap[key]
		split := script.Split{
			ScriptId: scr.Id,
			Content:  &group.content,
			Effect:   &group.effect,
			PanelId:  group.panelId,
		}
		splits = append(splits, split)
	}

	err = s.script.DeleteSplits(scriptId)
	if err != nil {
		return err
	}

	if len(splits) > 0 {
		_, err = s.script.CreateManySplit(splits)
		if err != nil {
			return err
		}
	}

	return nil
}

func resolvePanelId(panels []book.Panel, panelNum int) *int64 {
	if panelNum > len(panels) {
		return &panels[len(panels)-1].Id
	}
	if panelNum <= 0 {
		return &panels[0].Id
	}
	for i := range panels {
		if panels[i].PanelNumber == panelNum {
			return &panels[i].Id
		}
	}
	return &panels[0].Id
}

func NewGenerateSplits(script script.Interface, book book.Interface, ai ai.Interface) *GenerateSplits {
	return &GenerateSplits{script, book, ai}
}
