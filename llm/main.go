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

type LanguageModel interface {
	GenerateText(context.Context, string) (string, error)
}

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

type lLMOption func(interface{})

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
