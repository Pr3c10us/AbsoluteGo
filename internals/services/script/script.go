package script

import (
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
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
	DeleteSplits   *commands.DeleteSplits
	GenerateScript *commands.GenerateScript
	CreateScript   *commands.CreateScript
	CreateSplits   *commands.CreateSplits
	GenerateSplits *commands.GenerateSplits
}

type Queries struct {
	GetScripts *queries.GetScripts
	GetSplits  *queries.GetSplits
}

func NewScriptServices(script script.Interface, book book.Interface, ai ai.Interface, eventImplementation event.Interface, queueImplementation queue.Interface) Services {
	return Services{
		Commands: Commands{
			DeleteScript:   commands.NewDeleteScript(script),
			DeleteSplits:   commands.NewDeleteSplits(script),
			GenerateScript: commands.NewGenerateScript(script, book, ai),
			CreateScript:   commands.NewCreateScript(eventImplementation, book, queueImplementation, script),
			GenerateSplits: commands.NewGenerateSplits(script, book, ai),
			CreateSplits:   commands.NewCreateSplits(eventImplementation, book, queueImplementation, script),
		},
		Queries: Queries{
			GetScripts: queries.NewGetScripts(script),
			GetSplits:  queries.NewGetSplits(script),
		},
	}
}
