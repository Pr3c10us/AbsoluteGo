package script

import (
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/internals/domains/vab"
	"github.com/Pr3c10us/absolutego/internals/services/script/commands"
	"github.com/Pr3c10us/absolutego/internals/services/script/queries"
	"github.com/Pr3c10us/absolutego/packages/configs"
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
	GenerateAudio  *commands.GenerateAudio
	GenerateAudios *commands.GenerateAudios
	CreateAudio    *commands.CreateAudio
	CreateAudios   *commands.CreateAudios
	GenerateVideo  *commands.GenerateVideo
	GenerateVideos *commands.GenerateVideos
	CreateVideo    *commands.CreateVideo
	CreateVideos   *commands.CreateVideos
}

type Queries struct {
	GetScripts *queries.GetScripts
	GetSplits  *queries.GetSplits
}

func NewScriptServices(script script.Interface, book book.Interface, ai ai.Interface, eventImplementation event.Interface, queueImplementation queue.Interface, storageImplementation storage.Interface, environmentVariables *configs.EnvironmentVariables, vabImplementation vab.Interface) Services {
	return Services{
		Commands: Commands{
			DeleteScript:   commands.NewDeleteScript(script, storageImplementation, vabImplementation),
			DeleteSplits:   commands.NewDeleteSplits(script, storageImplementation),
			GenerateScript: commands.NewGenerateScript(script, book, ai),
			CreateScript:   commands.NewCreateScript(eventImplementation, book, queueImplementation, script),
			CreateSplits:   commands.NewCreateSplits(eventImplementation, book, queueImplementation, script),
			GenerateSplits: commands.NewGenerateSplits(script, book, ai),
			GenerateAudio:  commands.NewGenerateAudio(script, ai, storageImplementation, environmentVariables),
			GenerateAudios: commands.NewGenerateAudios(book, script, ai, storageImplementation, environmentVariables),
			CreateAudio:    commands.NewCreateAudio(eventImplementation, book, queueImplementation, script),
			CreateAudios:   commands.NewCreateAudios(eventImplementation, book, queueImplementation, script),
			GenerateVideo:  commands.NewGenerateVideo(book, storageImplementation, environmentVariables, script),
			GenerateVideos: commands.NewGenerateVideos(book, script, storageImplementation, environmentVariables),
			CreateVideo:    commands.NewCreateVideo(eventImplementation, book, queueImplementation, script),
			CreateVideos:   commands.NewCreateVideos(eventImplementation, book, queueImplementation, script),
		},
		Queries: Queries{
			GetScripts: queries.NewGetScripts(script),
			GetSplits:  queries.NewGetSplits(script),
		},
	}
}
