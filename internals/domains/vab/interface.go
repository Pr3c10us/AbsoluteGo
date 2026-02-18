package vab

type Interface interface {
	Create(vab VAB) (int64, error)
	Delete(id int64) error
	DeleteByScript(scriptId int64) error
	Update(id int64, vab VAB) error
	GetVABs(name string, scriptId, bookId int64, page, limit int) ([]VAB, error)
	GetVAB(id int64) (*VAB, error)
}
