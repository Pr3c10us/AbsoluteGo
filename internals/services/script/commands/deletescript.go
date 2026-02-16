package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"

	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type DeleteScript struct {
	scriptImplementation  script.Interface
	storageImplementation storage.Interface
	deleteSplits          *DeleteSplits
}

func (s *DeleteScript) Handle(scriptId int64) error {
	sc, err := s.scriptImplementation.GetScript(scriptId)
	if err != nil {
		return err
	}
	if sc == nil {
		return appError.BadRequest(errors.New("script does not exist"))
	}

	err = s.deleteSplits.Handle(scriptId)
	if err != nil {
		return err
	}

	return s.scriptImplementation.DeleteScript(scriptId)
}

func NewDeleteScript(scriptImplementation script.Interface, storageImplementation storage.Interface) *DeleteScript {
	return &DeleteScript{
		scriptImplementation:  scriptImplementation,
		storageImplementation: storageImplementation,
		deleteSplits:          NewDeleteSplits(scriptImplementation, storageImplementation),
	}
}
