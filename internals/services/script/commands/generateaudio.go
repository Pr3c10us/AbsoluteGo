package commands

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/prompts"
	"github.com/Pr3c10us/absolutego/packages/utils"
	"log"
	"os"
	"path/filepath"
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

	buf, err := base64.StdEncoding.DecodeString(resp.Response)
	if err != nil {
		return 0, fmt.Errorf("failed to decode base64: %w", err)
	}
	silence := 1.0

	duration := utils.BufDuration(buf, 24000, 1, 2) + silence

	tempDir, err := utils.GetDirectory("tmp")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tempDir)
	audioPath := filepath.Join(tempDir, "audio.wav")

	err = utils.WriteWAV(audioPath, resp.Response, 24000, 1, 16, silence)
	if err != nil {
		log.Fatal(err)
	}

	osFile, err := os.Open(audioPath)
	if err != nil {
		return 0, err
	}
	defer osFile.Close()

	url, err := service.storageImplementation.UploadFile(service.environmentVariables.Buckets.AudioBucket, osFile)
	if err != nil {
		return 0, err
	}

	err = service.scriptImplementation.UpdateSplit(split.Id, &script.Split{
		AudioURL:      &url,
		AudioDuration: &duration,
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
