package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/prompts"
	"github.com/Pr3c10us/absolutego/packages/utils"
	"log"
	"os"
)

type GenerateAudio struct {
	scriptImplementation  script.Interface
	aiImplementation      ai.Interface
	storageImplementation storage.Interface
	environmentVariables  *configs.EnvironmentVariables
}

type AudioParameter struct {
	Id         int64
	Voice      ai.Voice
	VoiceStyle string
}

func (service *GenerateAudio) Handle(parameter AudioParameter) (int64, error) {
	split, err := service.scriptImplementation.GetSplit(parameter.Id)
	if err != nil {
		return 0, err
	}
	if split == nil {
		return 0, errors.New("split does not exist")
	}

	scr, err := service.scriptImplementation.GetScript(split.ScriptId)
	if err != nil {
		return 0, err
	}
	if scr == nil {
		return 0, errors.New("script does not exist")
	}

	prompt := prompts.AudioPrompt(*split.Content, parameter.VoiceStyle, split.PreviousContent)

	resp, err := service.aiImplementation.GenerateAudio(prompt, parameter.Voice)
	if err != nil {
		return 0, err
	}

	tempDir, err := utils.GetDirectory("tmp")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tempDir)

	err = utils.WriteWAV(tempDir, resp.Response, 24000, 1, 16)
	if err != nil {
		log.Fatal(err)
	}

	osFile, err := os.Open(tempDir)
	if err != nil {
		return 0, err
	}
	defer osFile.Close()

	url, err := service.storageImplementation.UploadFile(service.environmentVariables.Buckets.AudioBucket, osFile)
	if err != nil {
		return 0, err
	}

	err = service.scriptImplementation.UpdateSplit(split.Id, &script.Split{
		AudioURL: &url,
	})
	if err != nil {
		return 0, err
	}

	return scr.Id, nil
}

func NewGenerateAudio(scriptImplementation script.Interface, aiImplementation ai.Interface, storageImplementation storage.Interface, environmentVariables *configs.EnvironmentVariables) *GenerateAudio {
	return &GenerateAudio{
		scriptImplementation, aiImplementation, storageImplementation, environmentVariables,
	}
}
