package llm

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gage-technologies/mistral-go"
	"github.com/google/generative-ai-go/genai"
	"github.com/liushuangls/go-anthropic/v2"
	"google.golang.org/api/option"
)

// ToolType is a type constraint for the allowed tool types
type ToolType interface {
	mistral.Tool | anthropic.ToolDefinition | genai.Tool
}

// GenericTool is a generic struct that can hold specific types of tools
type GenericTool[T ToolType] struct {
	Tool T
}

// GenerateOptions hold common options for text generations
type GenerateOptions[T ToolType] struct {
	Tools            []GenericTool[T]
	ResponseMIMEType string
}

// LanguageModel defines a common interface for interacting with Large Language Models (LLMs)
type LanguageModel[T ToolType] interface {
	GenerateText(ctx context.Context, prompt string, opts *GenerateOptions[T]) (string, error)
}

type LLMConfig struct {
	ModelName   string
	Temperature float32
	MaxTokens   int
	TopP        float32
}

func NewLLM[T ToolType](llmType string, opts ...LLMOption) (LanguageModel[T], error) {
	config := LLMConfig{
		ModelName:   "mistral-small-latest",
		Temperature: 0.7,
		MaxTokens:   512,
		TopP:        1,
	}

	for _, opt := range opts {
		opt(config)
	}

	switch llmType {
	case "anthropic":
		return newAnthropicLLM[T](config)
	case "mistral":
		return newMistralLLM[T](config)
	case "gemini":
		return newGeminiLLM[T](config)
	default:
		return nil, fmt.Errorf("unsuported LLM type: %s", llmType)
	}
}

func newAnthropicLLM[T ToolType](config LLMConfig) (LanguageModel[T], error) {
	apiKey, ok := os.LookupEnv("CLAUDE_API_KEY")
	if !ok {
		return nil, fmt.Errorf("the CLAUDE_API_KEY was not set")
	}

	llm := &anthropicLLM[T]{
		config: config,
		client: anthropic.NewClient(apiKey),
	}

	return llm, nil
}

func newMistralLLM[T ToolType](config LLMConfig) (LanguageModel[T], error) {
	llm := &mistralLLM[T]{
		config: config,
		// It will look for MISTRAL_API_KEY environment variable
		client: mistral.NewMistralClientDefault(""),
	}

	return llm, nil
}

func newGeminiLLM[T ToolType](config LLMConfig) (LanguageModel[T], error) {
	ctx := context.Background()

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		return nil, fmt.Errorf("the Environment variable GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	llm := &geminiLLM[T]{
		config: config,
		client: client,
	}

	return llm, nil
}

type LLMOption func(LLMConfig)

func WithMaxTokens(maxTokens int) LLMOption {
	return func(c LLMConfig) {
		c.MaxTokens = maxTokens
	}
}

func WithModelName(modelName string) LLMOption {
	return func(c LLMConfig) {
		c.ModelName = modelName
	}
}

func WithTemperature(temperature float32) LLMOption {
	return func(c LLMConfig) {
		c.Temperature = temperature
	}
}

func WithTopP(topP float32) LLMOption {
	return func(c LLMConfig) {
		c.TopP = topP
	}
}

// NewTool creates a new GenericTool
func NewTool[T ToolType](tool T) GenericTool[T] {
	return GenericTool[T]{
		Tool: tool,
	}
}
