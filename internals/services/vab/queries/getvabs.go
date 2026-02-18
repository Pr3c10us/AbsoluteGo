package queries

import "github.com/Pr3c10us/absolutego/internals/domains/vab"

type GetVABs struct {
	vabImplementation vab.Interface
}

func (service *GetVABs) Handle(scriptId, bookId int64, name string) ([]vab.VAB, error) {
	return service.vabImplementation.GetVABs(name, scriptId, bookId)
}

func NewGetVABs(vabImplementation vab.Interface) *GetVABs {
	return &GetVABs{
		vabImplementation,
	}
}
