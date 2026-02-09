package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type AddBook struct {
	bookImplementation    book.Interface
	storageImplementation storage.Interface
	aiImplementation      ai.Interface
	environmentVariables  *configs.EnvironmentVariables
}

type Parameter struct {
	File    string
	Chapter int
}

func (service *AddBook) Handle(parameter Parameter) error {
	defer os.Remove(parameter.File)
	format, err := utils.GetComicFormat(parameter.File)
	if err != nil {
		return err
	}

	tempDir, err := utils.GetDirectory("tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	outputDir, err := utils.GetDirectory("books")
	if err != nil {
		return err
	}

	err = utils.ExtractComicToDir(parameter.File, format, tempDir)
	if err != nil {
		return err
	}

	images, err := utils.SortImages(tempDir, outputDir)
	if err != nil {
		return err
	}

	if len(images) < 1 {
		return errors.New("no images extracted")
	}

	_, err = utils.GenerateBlurCover(outputDir, images[0])

	const BatchSize = 20

	for i := 0; i < len(images); i += BatchSize {
		end := i + BatchSize
		if end > len(images) {
			end = len(images)
		}
		batch := images[i:end]

		var wg sync.WaitGroup

		for _, imagePath := range batch {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				defer os.Remove(path)
				utils.DetectAndExtractPanels(path)
			}(imagePath)
		}

		wg.Wait()
	}

	overlayImages := make([]string, 0)
	for _, imagePath := range images {
		dir := filepath.Dir(imagePath)
		name := strings.TrimSuffix(filepath.Base(imagePath), filepath.Ext(imagePath))
		overlayPath := filepath.Join(dir, name+".png")

		if _, err := os.Stat(overlayPath); err == nil {
			overlayImages = append(overlayImages, overlayPath)
		}
	}

	var wg sync.WaitGroup
	errs := make([]error, len(overlayImages))

	for i, overlayPath := range overlayImages {
		wg.Add(1)
		go func(path string, pageNum int, idx int) {
			defer wg.Done()
			errs[idx] = utils.AddPageNumberToOverlay(path, pageNum)
		}(overlayPath, i+1, i)
	}

	wg.Wait()

	for _, err = range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func NewAddBook(bookImplementation book.Interface, storageImplementation storage.Interface, aiImplementation ai.Interface, environmentVariables *configs.EnvironmentVariables) *AddBook {
	return &AddBook{
		bookImplementation, storageImplementation, aiImplementation, environmentVariables,
	}
}
