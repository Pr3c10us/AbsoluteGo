package queries

import "github.com/Pr3c10us/absolutego/internals/domains/script"

type GetScripts struct {
	scriptImplementation script.Interface
}

func (service *GetScripts) Handle(bookId int64, name string) ([]script.Script, error) {
	return service.scriptImplementation.GetScripts(bookId, name)
}

func NewGetScripts(scriptImplementation script.Interface) *GetScripts {
	return &GetScripts{
		scriptImplementation,
	}
}
