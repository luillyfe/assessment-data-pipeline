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
type mistralLLM struct {
	modelName   string
	temperature float64
	maxTokens   int
	topP        float64
	client      MistralClient
}

/*
GenerateText generates text using the Mistral LLM based on the provided prompt.

It takes a context.Context and a prompt string as input.
It constructs a Mistral ChatRequest with the prompt and model parameters.
It sends the request to the Mistral API using the client.
It handles potential errors from the Mistral API.
It extracts and returns the generated text from the API response.

Args:

	ctx: The context for the request.
	prompt: The input prompt for text generation.

Returns:

	A string containing the generated text and an error if any occurred.
*/
func (m *mistralLLM) GenerateText(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	// Tool handling
	var mistralTools []mistral.Tool
	if opts != nil && len(opts.Tools) > 0 {
		for _, genericTool := range opts.Tools {
			if genericTool.Type != MistralToolType {
				return "", fmt.Errorf("error: tool type mismatch for Mistral LLM")
			}
			mistralTool, ok := genericTool.Tool.(mistral.Tool)
			if !ok {
				return "", fmt.Errorf("error: invalid tool type for Mistral LLM")
			}
			mistralTools = append(mistralTools, mistralTool)
		}
	}

	// Using chat completion
	resp, err := m.client.Chat(m.modelName, []mistral.ChatMessage{{Content: prompt, Role: mistral.RoleUser}}, &mistral.ChatRequestParams{
		Temperature: m.temperature,
		MaxTokens:   m.maxTokens,
		TopP:        m.topP,
		Tools:       mistralTools,
	})
	if err != nil {
		return "", fmt.Errorf("error getting chat completion: %w", err)
	}

	// Return generated text
	return resp.Choices[0].Message.Content, nil
}
