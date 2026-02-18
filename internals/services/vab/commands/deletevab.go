package commands

import (
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/internals/domains/vab"
	"github.com/Pr3c10us/absolutego/packages/appError"
)

type DeleteVAB struct {
	vabImplementation     vab.Interface
	storageImplementation storage.Interface
}

func (service *DeleteVAB) Handle(vabId, scriptId int64) error {
	if vabId != 0 {
		v, err := service.vabImplementation.GetVAB(vabId)
		if err != nil {
			return err
		}
		if v == nil {
			return appError.NotFound(fmt.Errorf("vab does not exist"))
		}

		err = service.vabImplementation.Delete(v.Id)
		if err != nil {
			return err
		}

		err = service.storageImplementation.DeleteFile(*v.Url)
		if err != nil {
			return err
		}
	} else if scriptId != 0 {
		vabs, err := service.vabImplementation.GetVABs("", scriptId, 0)
		if err != nil {
			return err
		}

		var urls []string
		for _, v := range vabs {
			urls = append(urls, *v.Url)
		}

		err = service.vabImplementation.DeleteByScript(scriptId)
		if err != nil {
			return err
		}

		_ = service.storageImplementation.DeleteMany(urls)
	}

	return nil
}

func NewDeleteVAB(vabImplementation vab.Interface, storageImplementation storage.Interface) *DeleteVAB {
	return &DeleteVAB{
		vabImplementation, storageImplementation,
	}
}
