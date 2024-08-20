package llm

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

type geminiLLM struct {
	modelName   string
	temperature float64
	maxTokens   int
	topP        float64
	client      *genai.Client
}

func (g *geminiLLM) GenerateText(ctx context.Context, prompt string) (string, error) {
	model := g.client.GenerativeModel(g.modelName)

	model.SetTemperature(float32(g.temperature))
	model.SetTopP(float32(g.topP))
	model.SetMaxOutputTokens(int32(g.maxTokens))
	model.SetTopK(64)
	model.ResponseMIMEType = "text/plain"

	// model.SafetySettings = Adjust safety settings
	// See https://ai.google.dev/gemini-api/docs/safety-settings

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	output := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		output = fmt.Sprintf("%v\n", part)
	}

	// Return generated text
	return output, nil
}
