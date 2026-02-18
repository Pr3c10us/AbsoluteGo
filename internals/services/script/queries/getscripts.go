package queries

import "github.com/Pr3c10us/absolutego/internals/domains/script"

type GetScripts struct {
	scriptImplementation script.Interface
}

func (service *GetScripts) Handle(bookId int64, name string, ids []int64, page, limit int) ([]script.Script, error) {
	return service.scriptImplementation.GetScripts(script.Query{
		bookId, name, ids, 0, page, limit,
	})
}

func NewGetScripts(scriptImplementation script.Interface) *GetScripts {
	return &GetScripts{
		scriptImplementation,
	}
}
