package vab

import (
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/internals/domains/vab"
	"github.com/Pr3c10us/absolutego/internals/services/vab/commands"
	"github.com/Pr3c10us/absolutego/internals/services/vab/queries"
	"github.com/Pr3c10us/absolutego/packages/configs"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
	CreateVAB *commands.CreateVAB
	QueueVAB  *commands.QueueVAB
	DeleteVAB *commands.DeleteVAB
}

type Queries struct {
	GetVABs *queries.GetVABs
}

func NewVABServices(vabImplementation vab.Interface, scriptImplementation script.Interface, bookImplementation book.Interface, eventImplementation event.Interface, queueImplementation queue.Interface, storageImplementation storage.Interface, environmentVariables *configs.EnvironmentVariables) Services {
	return Services{
		Commands: Commands{
			CreateVAB: commands.NewCreateVAB(vabImplementation, bookImplementation, scriptImplementation, storageImplementation, environmentVariables),
			QueueVAB:  commands.NewQueueVAB(vabImplementation, scriptImplementation, bookImplementation, eventImplementation, queueImplementation),
			DeleteVAB: commands.NewDeleteVAB(vabImplementation, storageImplementation),
		},
		Queries: Queries{
			GetVABs: queries.NewGetVABs(vabImplementation),
		},
	}
}
