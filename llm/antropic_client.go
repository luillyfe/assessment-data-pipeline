package llm

import (
	"context"
	"errors"
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"
)

/*
AnthropicClient is an interface for interacting with the Anthropic API.

It defines a single method, CreateMessages, which sends a request to the Anthropic API
to generate text based on a given prompt and model parameters.
*/
type AnthropicClient interface {
	CreateMessages(ctx context.Context, request anthropic.MessagesRequest) (response anthropic.MessagesResponse, err error)
}

/*
anthropicLLM represents an Anthropic Large Language Model.

It implements the LanguageModel interface, providing text generation capabilities
using the Anthropic API.

Fields:

	modelName: The name of the Anthropic model to use for text generation.
	           e.g., "anthropic.ModelClaudeInstant1Dot2", "anthropic.ModelClaude2"

	temperature: Controls the randomness of the generated text.
	             Higher values (closer to 1) result in more random text,
	             while lower values (closer to 0) make the text more deterministic.

	maxTokens: The maximum number of tokens allowed in the generated text.

	topP: Sets the nucleus sampling threshold for the generated text.
	      This parameter controls the diversity of the generated text.

	client: An instance of the AnthropicClient interface, used to interact with the Anthropic API.
*/
type anthropicLLM struct {
	modelName   string
	temperature float64
	maxTokens   int
	topP        float64
	client      AnthropicClient
}

/*
GenerateText generates text using the Anthropic LLM based on the provided prompt.

It takes a context.Context and a prompt string as input.
It constructs an Anthropic MessagesRequest with the prompt and model parameters.
It sends the request to the Anthropic API using the client.
It handles potential errors, including Anthropic API errors.
It extracts and returns the generated text from the API response.

Args:

	ctx: The context for the request.
	prompt: The input prompt for text generation.

Returns:

	A string containing the generated text and an error if any occurred.
*/
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

	// Return generated text
	return *resp.Content[0].Text, nil
}
