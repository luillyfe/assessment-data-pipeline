package llm

import (
	"context"
	"errors"
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"
)

type AnthropicClient interface {
	CreateMessages(ctx context.Context, request anthropic.MessagesRequest) (response anthropic.MessagesResponse, err error)
}

type anthropicLLM struct {
	modelName   string
	temperature float64
	maxTokens   int
	topP        float64
	client      AnthropicClient
}

func (a *anthropicLLM) GenerateText(ctx context.Context, prompt string) (string, error) {
	// Cast to float32
	temperature := float32(a.temperature)
	topP := float32(a.topP)

	// Using chat completion
	resp, err := a.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: a.modelName,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(prompt),
		},
		MaxTokens:   a.maxTokens,
		Temperature: &temperature,
		TopP:        &topP,
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			return "", fmt.Errorf("anthropic API error, type: %s, message: %s", e.Type, e.Message)
		}
		return "", fmt.Errorf("anthropic API error: %w", err)
	}

	// return generated text
	return *resp.Content[0].Text, nil
}
