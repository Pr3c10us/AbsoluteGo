package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
)

type GenerateVideos struct {
	bookImplementation   book.Interface
	scriptImplementation script.Interface
	generateVideo        *GenerateVideo
}

type GenerateVideosParameter struct {
	ScriptId int64
	Width    int
	Height   int
	FPS      int
}

func (service *GenerateVideos) Handle(parameter GenerateVideosParameter) (int64, error) {
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

	maxWorkers := 5
	err = utils.RunWorkerPool(splits, maxWorkers, func(j script.Split) error {
		_, _ = service.generateVideo.Handle(
			GenerateVideoParameter{
				Id:     j.Id,
				Width:  parameter.Width,
				Height: parameter.Height,
				FPS:    parameter.FPS,
			})
		return nil
	})

	return scr.Id, nil
}

func NewGenerateVideos(bookImplementation book.Interface, scriptImplementation script.Interface, storageImplementation storage.Interface, environmentVariables *configs.EnvironmentVariables) *GenerateVideos {
	return &GenerateVideos{
		bookImplementation, scriptImplementation, NewGenerateVideo(bookImplementation, storageImplementation, environmentVariables, scriptImplementation),
	}
}
