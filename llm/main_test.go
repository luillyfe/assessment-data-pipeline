package llm

import (
	"context"
	"testing"

	"github.com/gage-technologies/mistral-go"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/go-cmp/cmp"
	"github.com/liushuangls/go-anthropic/v2"
)

type mockMistralClient struct{}

func (m *mockMistralClient) Chat(model string, messages []mistral.ChatMessage, params *mistral.ChatRequestParams) (*mistral.ChatCompletionResponse, error) {
	return &mistral.ChatCompletionResponse{
		Choices: []mistral.ChatCompletionResponseChoice{{Message: mistral.ChatMessage{Content: "Mistral Response"}}}}, nil
}

type mockAnthropicClient struct{}

func (m *mockAnthropicClient) CreateMessages(ctx context.Context, request anthropic.MessagesRequest) (response anthropic.MessagesResponse, err error) {
	text := "Anthropic Response"
	return anthropic.MessagesResponse{
		Content: []anthropic.MessageContent{{Text: &text}},
	}, nil
}

func TestGenerateText(t *testing.T) {
	tests := []struct {
		name         string
		llm          LanguageModel
		prompt       string
		opts         *GenerateOptions
		want         string
		wantErr      bool
		setEnvVars   func()
		unsetEnvVars func()
	}{
		{
			name: "Mistral Success",
			llm: &mistralLLM{
				modelName:   "mistral-small-latest",
				temperature: 0.7,
				maxTokens:   512,
				topP:        1,
				client:      &mockMistralClient{},
			},
			prompt: "Hello, how are you?",
			opts: &GenerateOptions{
				Tools: []GenericTool{
					NewMistralTool(mistral.Tool{
						Type: mistral.ToolTypeFunction,
						Function: mistral.Function{
							Name:        "test_function",
							Description: "A test function",
							Parameters:  map[string]interface{}{},
						},
					}),
				},
			},
			want: "Mistral Response",
		},
		{
			name: "Anthropic Success",
			llm: &anthropicLLM{
				modelName:   anthropic.ModelClaudeInstant1Dot2,
				temperature: 0.7,
				maxTokens:   512,
				topP:        1,
				client:      &mockAnthropicClient{},
			},
			prompt: "Hello, how are you?",
			opts: &GenerateOptions{
				Tools: []GenericTool{
					NewAnthropicTool(anthropic.ToolDefinition{
						Name:        "test_function",
						Description: "A test function",
						InputSchema: map[string]interface{}{},
					}),
				},
			},
			want: "Anthropic Response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.llm.GenerateText(context.Background(), tt.prompt, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GenerateText() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateTextWithNoTools(t *testing.T) {
	llm := &anthropicLLM{
		modelName:   anthropic.ModelClaudeInstant1Dot2,
		temperature: 0.7,
		maxTokens:   512,
		topP:        1,
		client:      &mockAnthropicClient{},
	}

	opts := &GenerateOptions{} // No tools

	got, err := llm.GenerateText(context.Background(), "Test prompt", opts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	want := "Anthropic Response"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GenerateText() mismatch (-want +got):\n%s", diff)
	}
}

func TestGenerateTextWithInvalidTools(t *testing.T) {
	llm := &mistralLLM{
		modelName:   "gemini-1.5-pro-exp-0801",
		temperature: 0.7,
		maxTokens:   512,
		topP:        1,
		client:      &mockMistralClient{},
	}

	opts := &GenerateOptions{
		Tools: []GenericTool{
			NewGeminiTool(&genai.Tool{}), // Invalid tool type for Mistral
		},
	}

	_, err := llm.GenerateText(context.Background(), "Test prompt", opts)
	if err == nil {
		t.Errorf("Expected error for invalid tool type, got nil")
	}
}
