package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
)

type GenerateAudios struct {
	bookImplementation   book.Interface
	scriptImplementation script.Interface
	generateAudio        *GenerateAudio
}

type GenerateAudiosParameters struct {
	ScriptId   int64
	Voice      ai.Voice
	VoiceStyle string
}

func (service *GenerateAudios) Handle(parameters GenerateAudiosParameters) (int64, error) {
	scr, err := service.scriptImplementation.GetScript(parameters.ScriptId)
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

	maxWorkers := 20
	err = utils.RunWorkerPool(splits, maxWorkers, func(j script.Split) error {
		_, _ = service.generateAudio.Handle(AudioParameter{
			Id:         j.Id,
			Voice:      parameters.Voice,
			VoiceStyle: parameters.VoiceStyle,
		})
		return nil
	})

	return scr.Id, nil
}

func NewGenerateAudios(bookImplementation book.Interface, scriptImplementation script.Interface, aiImplementation ai.Interface, storageImplementation storage.Interface, environmentVariables *configs.EnvironmentVariables) *GenerateAudios {
	return &GenerateAudios{
		bookImplementation, scriptImplementation, NewGenerateAudio(scriptImplementation, aiImplementation, storageImplementation, environmentVariables),
	}
}
