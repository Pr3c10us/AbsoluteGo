package commands

import (
	"github.com/Pr3c10us/absolutego/internals/domains/script"
)

type DeleteSplits struct {
	scriptImplementation script.Interface
}

func (s *DeleteSplits) Handle(scriptId int64) error {
	err := s.scriptImplementation.DeleteSplits(scriptId)
	return err
}

func NewDeleteSplits(scriptImplementation script.Interface) *DeleteSplits {
	return &DeleteSplits{
		scriptImplementation: scriptImplementation,
	}
}
