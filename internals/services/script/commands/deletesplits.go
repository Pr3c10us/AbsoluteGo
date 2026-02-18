package commands

import (
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
)

type DeleteSplits struct {
	scriptImplementation  script.Interface
	storageImplementation storage.Interface
}

func (service *DeleteSplits) Handle(scriptId int64) error {
	splits, err := service.scriptImplementation.GetSplits(scriptId)
	if err != nil {
		return err
	}
	if len(splits) < 1 {
		return nil
	}

	var urls []string
	for _, split := range splits {
		if split.AudioURL != nil {
			urls = append(urls, *split.AudioURL)
		}
		if split.VideoURL != nil {
			urls = append(urls, *split.VideoURL)
		}
	}

	_ = service.storageImplementation.DeleteMany(urls)

	err = service.scriptImplementation.DeleteSplits(scriptId)
	return err
}

func NewDeleteSplits(scriptImplementation script.Interface, storageImplementation storage.Interface) *DeleteSplits {
	return &DeleteSplits{
		scriptImplementation: scriptImplementation, storageImplementation: storageImplementation,
	}
}
