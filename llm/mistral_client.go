package llm

import (
	"context"
	"fmt"

	"github.com/gage-technologies/mistral-go"
)

type mistralLLM struct {
	modelName   string
	temperature float64
	maxTokens   int
	topP        float64
	client      *mistral.MistralClient
}

func (m *mistralLLM) GenerateText(ctx context.Context, prompt string) (string, error) {
	// Using chat completion
	resp, err := m.client.Chat(m.modelName, []mistral.ChatMessage{{Content: prompt, Role: mistral.RoleUser}}, &mistral.ChatRequestParams{
		Temperature: m.temperature,
		MaxTokens:   m.maxTokens,
		TopP:        m.topP,
	})
	if err != nil {
		return "", fmt.Errorf("error getting chat completion: %w", err)
	}

	// return generated text
	return resp.Choices[0].Message.Content, nil
}
