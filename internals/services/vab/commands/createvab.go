package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/internals/domains/vab"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
	"os"
	"path/filepath"
)

type CreateVAB struct {
	vabImplementation     vab.Interface
	bookImplementation    book.Interface
	scriptImplementation  script.Interface
	storageImplementation storage.Interface
	environmentVariables  *configs.EnvironmentVariables
}

type CreateVABParameter struct {
	ScriptId int64
	Name     string
}

func (service *CreateVAB) Handle(parameter CreateVABParameter) (int64, error) {
	scr, err := service.scriptImplementation.GetScript(parameter.ScriptId)
	if err != nil {
		return 0, err
	}
	if scr == nil {
		return 0, errors.New("script does not exist")
	}

	b, err := service.bookImplementation.GetBook(scr.BookId)
	if err != nil {
		return 0, err
	}
	if b == nil {
		return 0, errors.New("book does not exist")
	}

	splits, err := service.scriptImplementation.GetSplits(scr.Id)
	if err != nil {
		return 0, err
	}
	if len(splits) < 1 {
		return 0, errors.New("no splits for script")
	}

	var videoURLs []string
	for _, split := range splits {
		if split.VideoURL != nil {
			videoURLs = append(videoURLs, *split.VideoURL)
		}
	}
	if len(videoURLs) < 1 {
		return 0, errors.New("no video generated for script")
	}

	var videoPaths []string
	for _, url := range videoURLs {
		var videoPath string
		videoPath, err = utils.DownloadPage(url)
		if err != nil {
			return 0, err
		}
		videoPaths = append(videoPaths, videoPath)
	}

	videoDir, err := utils.GetDirectory("tmp")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(videoDir)

	outputPath := filepath.Join(videoDir, "video.mp4")

	err = utils.MergeVideos(videoPaths, outputPath, nil)
	if err != nil {
		return 0, err
	}

	osFile, err := os.Open(outputPath)
	if err != nil {
		return 0, err
	}
	defer osFile.Close()

	url, err := service.storageImplementation.UploadFile(service.environmentVariables.Buckets.VABBucket, osFile)
	if err != nil {
		return 0, err
	}

	var id int64
	id, err = service.vabImplementation.Create(vab.VAB{
		Name:     parameter.Name,
		Url:      &url,
		ScriptId: scr.Id,
		BookId:   b.Id,
	})

	if err != nil {
		service.storageImplementation.DeleteFile(url)
		return 0, err
	}

	return id, nil
}

func NewCreateVAB(vabImplementation vab.Interface, bookImplementation book.Interface, scriptImplementation script.Interface, storageImplementation storage.Interface, environmentVariables *configs.EnvironmentVariables) *CreateVAB {
	return &CreateVAB{
		vabImplementation, bookImplementation, scriptImplementation, storageImplementation, environmentVariables,
	}
}
