package script

import (
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/services/script/commands"
	"github.com/Pr3c10us/absolutego/internals/services/script/queries"
	"github.com/Pr3c10us/absolutego/packages/configs"
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

func NewScriptServices(scriptImplementation script.Interface, aiImplementation ai.Interface, environmentVariables *configs.EnvironmentVariables) Services {
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
