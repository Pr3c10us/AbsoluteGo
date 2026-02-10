package script

import (
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/services/script/commands"
	"github.com/Pr3c10us/absolutego/internals/services/script/queries"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
	DeleteScript *commands.DeleteScript
}

type Queries struct {
	GetScripts *queries.GetScripts
	GetSplits  *queries.GetSplits
}

func NewScriptServices(scriptImplementation script.Interface) Services {
	return Services{
		Commands: Commands{
			DeleteScript: commands.NewDeleteScript(scriptImplementation),
		},
		Queries: Queries{
			GetScripts: queries.NewGetScripts(scriptImplementation),
			GetSplits:  queries.NewGetSplits(scriptImplementation),
		},
	}
}
