package llm

import (
	"context"
	"fmt"

	"github.com/gage-technologies/mistral-go"
)

/*
MistralClient is an interface for interacting with the Mistral API.

It defines a single method, Chat, which sends a chat request to the Mistral API
to generate text based on a given model, messages, and parameters.
*/
type MistralClient interface {
	Chat(model string, messages []mistral.ChatMessage, params *mistral.ChatRequestParams) (*mistral.ChatCompletionResponse, error)
}

/*
mistralLLM represents a Mistral Large Language Model.

It implements the LanguageModel interface, providing text generation capabilities
using the Mistral API.

Fields:

	modelName: The name of the Mistral model to use for text generation.
	           e.g., "mistral-small-latest", "mistral-large"

	temperature: Controls the randomness of the generated text.
	             Higher values (closer to 1) result in more random text,
	             while lower values (closer to 0) make the text more deterministic.

	maxTokens: The maximum number of tokens allowed in the generated text.

	topP: Sets the nucleus sampling threshold for the generated text.
	      This parameter controls the diversity of the generated text.

	client: An instance of the MistralClient interface, used to interact with the Mistral API.
*/
type mistralLLM[T ToolType] struct {
	config LLMConfig
	client MistralClient
}

/*
GenerateText generates text using the Mistral LLM based on the provided prompt and optional generation options.

It takes a context.Context, a prompt string, and optional generation options as input.
It constructs a Mistral ChatRequest with the prompt and model parameters.
It sends the request to the Mistral API using the client.
It handles potential errors from the Mistral API.
It extracts and returns the generated text from the API response.

Args:

	ctx: The context for the request.
	prompt: The input prompt for text generation.
	opts: Optional generation options, such as tools.

Returns:

	A string containing the generated text and an error if any occurred.
*/
func (m *mistralLLM[T]) GenerateText(ctx context.Context, prompt string, opts *GenerateOptions[T]) (string, error) {
	// Tool handling
	var mistralTools []mistral.Tool
	if opts != nil && len(opts.Tools) > 0 {
		for _, opt := range opts.Tools {
			if mistralTool, ok := any(opt.Tool).(mistral.Tool); ok {
				mistralTools = append(mistralTools, mistralTool)
			} else {
				return "", fmt.Errorf("tool doesn't implement mistral.Tool")
			}
		}
	}

	// Using chat completion
	resp, err := m.client.Chat(m.config.ModelName, []mistral.ChatMessage{{Content: prompt, Role: mistral.RoleUser}}, &mistral.ChatRequestParams{
		Temperature: float64(m.config.Temperature),
		MaxTokens:   m.config.MaxTokens,
		TopP:        float64(m.config.TopP),
		Tools:       mistralTools,
	})
	if err != nil {
		return "", fmt.Errorf("error getting chat completion: %w", err)
	}

	// Return generated text
	return resp.Choices[0].Message.Content, nil
}
