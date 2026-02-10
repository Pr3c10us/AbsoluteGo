package ai

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
	"google.golang.org/genai"
	"time"
)

type ModelPricing struct {
	TextInputPerMillion   float64
	AudioInputPerMillion  float64
	TextOutputPerMillion  float64
	AudioOutputPerMillion float64
}

var modelPricing = map[string]ModelPricing{
	"gemini-3-pro-preview":                          {TextInputPerMillion: 2.00, TextOutputPerMillion: 12.00},
	"gemini-3-flash-preview":                        {TextInputPerMillion: 0.50, AudioInputPerMillion: 1.00, TextOutputPerMillion: 3.00},
	"gemini-2.5-pro":                                {TextInputPerMillion: 1.25, TextOutputPerMillion: 10.00},
	"gemini-2.5-flash":                              {TextInputPerMillion: 0.30, AudioInputPerMillion: 1.00, TextOutputPerMillion: 2.50},
	"gemini-2.5-flash-preview-tts":                  {TextInputPerMillion: 0.50, TextOutputPerMillion: 2.50, AudioOutputPerMillion: 10.00},
	"gemini-2.5-pro-preview-tts":                    {TextInputPerMillion: 1.00, TextOutputPerMillion: 10.00, AudioOutputPerMillion: 20.00},
	"gemini-2.5-flash-native-audio-preview-12-2025": {TextInputPerMillion: 0.50, AudioInputPerMillion: 3.00, TextOutputPerMillion: 2.00, AudioOutputPerMillion: 12.00},
	"gemini-2.0-flash":                              {TextInputPerMillion: 0.10, AudioInputPerMillion: 0.70, TextOutputPerMillion: 0.40},
	"gemini-2.0-flash-lite":                         {TextInputPerMillion: 0.075, TextOutputPerMillion: 0.30},
	"gemini-2.5-flash-lite":                         {TextInputPerMillion: 0.10, AudioInputPerMillion: 0.30, TextOutputPerMillion: 0.40},
	"gemini-2.5-flash-lite-preview-09-2025":         {TextInputPerMillion: 0.10, AudioInputPerMillion: 0.30, TextOutputPerMillion: 0.40},
	"gemini-2.5-flash-preview-09-2025":              {TextInputPerMillion: 0.30, AudioInputPerMillion: 1.00, TextOutputPerMillion: 2.50},
}

type GoogleAI struct {
	client *genai.Client
	config *configs.GeminiConfig
}

func NewGoogleAI(client *genai.Client, config *configs.GeminiConfig) ai.Interface {
	return &GoogleAI{client: client, config: config}
}

func (g *GoogleAI) UploadFiles(files []ai.File) ([]ai.UploadedFile, error) {
	ctx := context.Background()
	if len(files) == 0 {
		return nil, nil
	}

	type result struct {
		index int
		file  ai.UploadedFile
		err   error
	}

	ch := make(chan result, len(files))
	for i, f := range files {
		go func(i int, f ai.File) {
			uf, err := g.uploadFile(ctx, f)
			ch <- result{index: i, file: uf, err: err}
		}(i, f)
	}

	results := make([]result, len(files))
	for range files {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		results[r.index] = r
	}

	uploaded := make([]ai.UploadedFile, len(files))
	for i, r := range results {
		uploaded[i] = r.file
	}
	return uploaded, nil
}

func (g *GoogleAI) uploadFile(ctx context.Context, f ai.File) (ai.UploadedFile, error) {
	uploaded, err := utils.WithRetry(func() (*genai.File, error) {
		return g.client.Files.UploadFromPath(ctx, f.Path, &genai.UploadFileConfig{
			MIMEType: f.MIMEType,
		})
	}, 10, 300*time.Millisecond)
	if err != nil {
		return ai.UploadedFile{}, fmt.Errorf("upload file: %w", err)
	}
	if uploaded.URI == "" || uploaded.MIMEType == "" {
		return ai.UploadedFile{}, errors.New("uploaded file missing URI or MIME")
	}
	return ai.UploadedFile{URI: uploaded.URI, MIMEType: uploaded.MIMEType}, nil
}

func (g *GoogleAI) GenerateAudioLive(text string, voice ai.Voice) (*ai.Response, error) {
	ctx := context.Background()

	model := g.config.LiveModel
	if voice == "" {
		voice = "Algieba"
	}

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityAudio},
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: string(voice),
				},
			},
		},
	}

	session, err := g.client.Live.Connect(ctx, model, config)
	if err != nil {
		return nil, fmt.Errorf("live connect: %w", err)
	}
	defer func(session *genai.Session) {
		err = session.Close()
		if err != nil {
		}
	}(session)

	err = session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{
			genai.NewContentFromText(text, genai.RoleUser),
		},
		TurnComplete: genai.Ptr(true),
	})
	if err != nil {
		return nil, fmt.Errorf("send client content: %w", err)
	}

	var audioChunks [][]byte
	var totalInputTokens, totalOutputTokens int32

	for {
		msg, err := session.Receive()
		if err != nil {
			return nil, fmt.Errorf("receive: %w", err)
		}

		if msg.ServerContent != nil && msg.ServerContent.ModelTurn != nil {
			for _, part := range msg.ServerContent.ModelTurn.Parts {
				if part.InlineData != nil && len(part.InlineData.Data) > 0 {
					audioChunks = append(audioChunks, part.InlineData.Data)
				}
			}
		}

		if msg.UsageMetadata != nil {
			totalInputTokens += msg.UsageMetadata.PromptTokenCount
			totalOutputTokens += msg.UsageMetadata.ResponseTokenCount
		}

		if msg.ServerContent != nil && msg.ServerContent.TurnComplete {
			break
		}
	}

	if len(audioChunks) == 0 {
		return nil, errors.New("no audio data received from Live API")
	}

	totalLen := 0
	for _, chunk := range audioChunks {
		totalLen += len(chunk)
	}
	combined := make([]byte, 0, totalLen)
	for _, chunk := range audioChunks {
		combined = append(combined, chunk...)
	}
	combinedBase64 := base64.StdEncoding.EncodeToString(combined)

	// Calculate cost
	dollars := 0.0
	if pricing, ok := modelPricing[model]; ok {
		inputCost := (float64(totalInputTokens) / 1_000_000) * pricing.TextInputPerMillion
		outputCostPerMillion := pricing.TextOutputPerMillion
		if pricing.AudioOutputPerMillion > 0 {
			outputCostPerMillion = pricing.AudioOutputPerMillion
		}
		outputCost := (float64(totalOutputTokens) / 1_000_000) * outputCostPerMillion
		dollars = inputCost + outputCost
	}

	return &ai.Response{Response: combinedBase64, Dollars: dollars}, nil
}

func (g *GoogleAI) GenerateAudio(text string, voice ai.Voice) (*ai.Response, error) {
	response, err := utils.WithRetry(func() (*ai.Response, error) {
		return g.GenerateAudioLive(text, voice)
	}, 10, 300*time.Millisecond)
	return response, err
}

func (g *GoogleAI) GenerateText(prompt string, useFastModel bool, uploadedFiles []ai.UploadedFile) (*ai.Response, error) {
	ctx := context.Background()

	var parts []*genai.Part

	for _, uf := range uploadedFiles {
		active, err := g.waitForFileActive(ctx, uf.URI, 30*time.Second)
		if err != nil {
			return nil, fmt.Errorf("wait for file active: %w", err)
		}
		if !active {
			return nil, errors.New("file failed to become active")
		}
		parts = append(parts, genai.NewPartFromURI(uf.URI, uf.MIMEType))
	}

	parts = append(parts, genai.NewPartFromText(prompt))

	model := g.config.Model
	if useFastModel {
		model = g.config.FastModel
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	genConfig := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: genai.ThinkingLevelHigh,
		},
		Tools: []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
		},
	}

	response, err := utils.WithRetry(func() (*genai.GenerateContentResponse, error) {
		return g.client.Models.GenerateContent(ctx, model, contents, genConfig)
	}, 10, 300*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("generate content: %w", err)
	}

	text := response.Text()
	if text == "" {
		return nil, errors.New("failed to generate text: empty response")
	}

	dollars := g.calculateCost(model, response, "text")
	return &ai.Response{Response: text, Dollars: dollars}, nil
}

func (g *GoogleAI) waitForFileActive(ctx context.Context, fileURI string, maxWait time.Duration) (bool, error) {

	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		file, err := g.client.Files.Get(ctx, fileURI, nil)
		if err != nil {
			return false, fmt.Errorf("get file: %w", err)
		}
		switch file.State {
		case genai.FileStateActive:
			return true, nil
		case genai.FileStateFailed:
			return false, errors.New("file processing failed")
		}

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}
	return false, nil
}

func (g *GoogleAI) calculateCost(model string, response *genai.GenerateContentResponse, costType string) float64 {
	pricing, ok := modelPricing[model]
	if !ok {
		return 0
	}
	if response.UsageMetadata == nil {
		return 0
	}

	inputTokens := float64(response.UsageMetadata.PromptTokenCount)
	outputTokens := float64(response.UsageMetadata.CandidatesTokenCount)

	inputCost := (inputTokens / 1_000_000) * pricing.TextInputPerMillion

	outputCostPerMillion := pricing.TextOutputPerMillion
	if costType == "audio" && pricing.AudioOutputPerMillion > 0 {
		outputCostPerMillion = pricing.AudioOutputPerMillion
	}
	outputCost := (outputTokens / 1_000_000) * outputCostPerMillion

	return inputCost + outputCost
}
