package script

type Interface interface {
	CreateScript(script *Script) (int64, error)
	UpdateScript(id int64, script *Script) error
	DeleteScript(id int64) error
	GetScripts(bookId int64, name string) ([]Script, error)
	GetScript(id int64) (*Script, error)

	CreateSplit(split *Split) (int64, error)
	CreateManySplit(splits []Split) ([]Split, error)
	UpdateSplit(id int64, split *Split) error
	DeleteSplits(scriptId int64) error
	DeleteSplit(id int64) error
	GetSplits(scriptId int64) ([]Split, error)
	GetSplit(id int64) (*Split, error)
}
