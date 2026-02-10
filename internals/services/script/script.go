package script

import (
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/services/script/commands"
	"github.com/Pr3c10us/absolutego/internals/services/script/queries"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
	DeleteScript   *commands.DeleteScript
	GenerateScript *commands.GenerateScript
}

type Queries struct {
	GetScripts *queries.GetScripts
	GetSplits  *queries.GetSplits
}

func NewScriptServices(script script.Interface, book book.Interface, ai ai.Interface) Services {
	return Services{
		Commands: Commands{
			DeleteScript:   commands.NewDeleteScript(script),
			GenerateScript: commands.NewGenerateScript(script, book, ai),
		},
		Queries: Queries{
			GetScripts: queries.NewGetScripts(script),
			GetSplits:  queries.NewGetSplits(script),
		},
	}
}
