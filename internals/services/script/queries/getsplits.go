package queries

import "github.com/Pr3c10us/absolutego/internals/domains/script"

type GetSplits struct {
	scriptImplementation script.Interface
}

func (service *GetSplits) Handle(scriptId int64) ([]script.Split, error) {
	return service.scriptImplementation.GetSplits(scriptId)
}

func NewGetSplits(scriptImplementation script.Interface) *GetSplits {
	return &GetSplits{
		scriptImplementation,
	}
}
