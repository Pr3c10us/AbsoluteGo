package vab

type Interface interface {
	Create(vab VAB) (int64, error)
	Delete(id int64) error
	DeleteByScript(scriptId int64) error
	Update(id int64, vab VAB) error
	GetVABs(name string, scriptId int64, bookId int64) ([]VAB, error)
	GetVAB(id int64) (*VAB, error)
}
