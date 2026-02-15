package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
	"os"
)

type GenerateVideo struct {
	bookImplementation    book.Interface
	storageImplementation storage.Interface
	environmentVariable   *configs.EnvironmentVariables
	scriptImplementation  script.Interface
}

func (service *GenerateVideo) Handle(id int64) (int64, error) {
	split, err := service.scriptImplementation.GetSplit(id)
	if err != nil {
		return 0, err
	}
	if split == nil {
		return 0, errors.New("split does not exist")
	}

	scr, err := service.scriptImplementation.GetScript(split.ScriptId)
	if err != nil {
		return 0, err
	}
	if scr == nil {
		return 0, appError.BadRequest(errors.New("script does not exist"))
	}

	file, err := utils.DownloadPage(*split.Panel.URL)
	if err != nil {
		return 0, err
	}
	defer os.Remove(file)
	
	vidData := utils.VideoData{
		Panel:    file,
		Duration: split.AudioDuration,
		Effect:   ,
	}

}

func NewAddChapter(
	bookImplementation book.Interface,
	storageImplementation storage.Interface,
	environmentVariable *configs.EnvironmentVariables,
	scriptImplementation script.Interface,
) *GenerateVideo {
	return &GenerateVideo{
		bookImplementation,
		storageImplementation,
		environmentVariable,
		scriptImplementation,
	}
}
