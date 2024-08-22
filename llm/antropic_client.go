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
GenerateText generates text using the Anthropic LLM based on the provided prompt and optional generation options.

It takes a context.Context, a prompt string, and optional generation options as input.
It constructs an Anthropic MessagesRequest with the prompt and model parameters.
It sends the request to the Anthropic API using the client.
It handles potential errors, including Anthropic API errors.
It extracts and returns the generated text from the API response.

Args:

	ctx: The context for the request.
	prompt: The input prompt for text generation.
	opts: Optional generation options, such as tools.

Returns:

	A string containing the generated text and an error if any occurred.
*/
func (a *anthropicLLM) GenerateText(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	// Cast to float32
	temperature := float32(a.temperature)
	topP := float32(a.topP)

	// Tool handling
	var anthropicTools []anthropic.ToolDefinition
	if opts != nil && len(opts.Tools) > 0 {
		for _, genericTool := range opts.Tools {
			if genericTool.Type != AnthropicToolType {
				return "", fmt.Errorf("error: tool type mismatch for Anthropic LLM")
			}
			anthropicTool, ok := genericTool.Tool.(anthropic.ToolDefinition)
			if !ok {
				return "", fmt.Errorf("error: invalid tool type for Anthropic LLM")
			}
			anthropicTools = append(anthropicTools, anthropicTool)
		}
	}

	// Using chat completion
	resp, err := a.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: a.modelName,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(prompt),
		},
		MaxTokens:   a.maxTokens,
		Temperature: &temperature,
		TopP:        &topP,
		Tools:       anthropicTools,
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
