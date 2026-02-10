package commands

import (
	"github.com/Pr3c10us/absolutego/internals/domains/script"
)

type DeleteSplits struct {
	scriptImplementation script.Interface
}

func (s *DeleteSplits) Handle(splitIds []int64) error {
	err := s.scriptImplementation.DeleteSplit(splitIds)
	return err
}

func NewDeleteSplits(scriptImplementation script.Interface) *DeleteSplits {
	return &DeleteSplits{
		scriptImplementation: scriptImplementation,
	}
}
