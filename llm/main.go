/*
Package llm provides a common interface for interacting with different Large Language Models (LLMs).

It defines a LanguageModel interface that abstracts the core functionality of generating text from a given prompt.
This allows for easy integration of different LLM providers without modifying the core application logic.

The package currently supports the following LLM providers:

- Anthropic: Uses the Anthropic API to access Claude models.
- Mistral: Uses the Mistral API to access Mistral models.
- Google Gemini: Uses the Google Gemini API to access Gemini models.

Each LLM provider has its own factory function for creating a new LanguageModel instance:

- NewAnthropicLLM: Creates a new Anthropic LLM instance.
- NewMistralLLM: Creates a new Mistral LLM instance.
- NewGeminiClient: Creates a new Google Gemini LLM instance.

These factory functions take a variable number of lLMOption arguments to customize the model's settings, such as:

- Model Name: Specifies the name of the LLM model to use.
- Temperature: Controls the randomness of the generated text.
- Max Tokens: Sets the maximum number of tokens allowed in the generated text.
- Top P: Sets the nucleus sampling threshold for the generated text.

The package also provides helper functions for creating common lLMOptions:

- WithMaxTokens: Creates an lLMOption that sets the maximum number of tokens.
- WithModelName: Creates an lLMOption that sets the model name.

Example Usage:

```go
// Create a new Anthropic LLM instance with default settings.
llm := llm.NewAnthropicLLM()

// Generate text using the LLM.
text, err := llm.GenerateText(context.Background(), "Hello, how are you?")

	if err != nil {
		// Handle error.
	}

// Print the generated text.
fmt.Println(text)
*/
package llm

import (
	"context"
	"log"
	"os"

	"github.com/gage-technologies/mistral-go"
	"github.com/google/generative-ai-go/genai"
	"github.com/liushuangls/go-anthropic/v2"
	"google.golang.org/api/option"
)

// LanguageModel defines a common interface for interacting with different Large Language Models (LLMs).
// It provides a single method, GenerateText, for generating text from a given prompt.
type LanguageModel interface {
	// GenerateText takes a context and a prompt string as input, and returns the generated text and an error.
	GenerateText(context.Context, string) (string, error)
}

/*
NewAnthropicLLM creates a new instance of a LanguageModel using Anthropic's API.
It takes a variable number of lLMOption arguments to customize the model's settings.

The function reads the CLAUDE_API_KEY environment variable to authenticate with the Anthropic API.

By default, the function initializes the Anthropic LLM with the following settings:
  - Model Name: "anthropic.ModelClaudeInstant1Dot2"
  - Temperature: 0.7
  - Max Tokens: 512
  - Top P: 1

These default settings can be overridden by passing in lLMOption arguments.
For example, to change the model name to "anthropic.ModelClaude2", you would use the following code:

	llm := NewAnthropicLLM(WithModelName("anthropic.ModelClaude2"))

The function returns a LanguageModel interface that can be used to generate text.
*/
func NewAnthropicLLM(opts ...lLMOption) LanguageModel {
	CLAUDE_API_KEY := os.Getenv("CLAUDE_API_KEY")

	llm := &anthropicLLM{
		modelName:   anthropic.ModelClaudeInstant1Dot2,
		temperature: 0.7,
		maxTokens:   512,
		topP:        1,
		client:      anthropic.NewClient(CLAUDE_API_KEY),
	}

	for _, opt := range opts {
		opt(llm)
	}

	return llm
}

/*
NewMistralLLM creates a new instance of a LanguageModel using the Mistral API.
It takes a variable number of lLMOption arguments to customize the model's settings.

The function initializes the Mistral LLM with the following default settings:
  - Model Name: "mistral-small-latest"
  - Temperature: 0.7
  - Max Tokens: 512
  - Top P: 1

It automatically retrieves the Mistral API key from the "MISTRAL_API_KEY" environment variable.

These default settings can be overridden by passing in lLMOption arguments.
For example, to change the model name to "mistral-large", you would use the following code:

	llm := NewMistralLLM(WithModelName("mistral-large"))

The function returns a LanguageModel interface that can be used to generate text.
*/
func NewMistralLLM(opts ...lLMOption) LanguageModel {
	llm := &mistralLLM{
		modelName:   "mistral-small-latest",
		temperature: 0.7,
		maxTokens:   512,
		topP:        1,
		// It will look for MISTRAL_API_KEY environment variable
		client: mistral.NewMistralClientDefault(""),
	}

	for _, opt := range opts {
		opt(llm)
	}

	return llm
}

/*
NewGeminiClient creates a new instance of a LanguageModel using Google's Gemini API.
It takes a variable number of lLMOption arguments to customize the model's settings.

The function reads the GEMINI_API_KEY environment variable to authenticate with the Gemini API.
If the environment variable is not set, the function will log a fatal error and exit.

By default, the function initializes the Gemini LLM with the following settings:
  - Model Name: "gemini-1.5-pro-exp-0801"
  - Temperature: 0.7
  - Max Tokens: 512
  - Top P: 1

These default settings can be overridden by passing in lLMOption arguments.
For example, to change the model name to "gemini-pro", you would use the following code:

	llm := NewGeminiClient(WithModelName("gemini-pro"))

The function returns a LanguageModel interface that can be used to generate text.
*/
func NewGeminiClient(opts ...lLMOption) LanguageModel {
	ctx := context.Background()

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	llm := &geminiLLM{
		modelName:   "gemini-1.5-pro-exp-0801",
		temperature: 0.7,
		maxTokens:   512,
		topP:        1,
		client:      client,
	}

	for _, opt := range opts {
		opt(llm)
	}

	return llm
}

/*
lLMOption is a function type that represents an option that can be applied
to a LanguageModel.

It takes an empty interface as input, which allows it to be used with
different LLM implementations. The actual implementation of the option
is responsible for type-asserting the input to the correct LLM type
and setting the desired option.
*/
type lLMOption func(interface{})

/*
WithMaxTokens creates an lLMOption that sets the maximum number of tokens
allowed in the generated text for the given LanguageModel.

It takes an integer maxTokens as input, representing the maximum number
of tokens allowed.

It returns an lLMOption function that takes an empty interface as input.
This function uses a type switch to determine the concrete type of the
LanguageModel passed to it and sets the maxTokens property accordingly.
*/
func WithMaxTokens(maxTokens int) lLMOption {
	return func(l interface{}) {
		switch v := l.(type) {
		case *mistralLLM:
			v.maxTokens = maxTokens
		case *anthropicLLM:
			v.maxTokens = maxTokens
		case *geminiLLM:
			v.maxTokens = maxTokens
		}
	}
}

/*
WithModelName creates an lLMOption that sets the model name for the given LanguageModel.

It takes a string modelName as input, representing the desired model name.

It returns an lLMOption function that takes an empty interface as input.
This function uses a type switch to determine the concrete type of the
LanguageModel passed to it and sets the modelName property accordingly.
*/
func WithModelName(modelName string) lLMOption {
	return func(l interface{}) {
		switch v := l.(type) {
		case *mistralLLM:
			v.modelName = modelName
		case *anthropicLLM:
			v.modelName = modelName
		case *geminiLLM:
			v.modelName = modelName
		}
	}
}
