package llm

import (
	"context"
	"testing"

	"github.com/gage-technologies/mistral-go"
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
			want:   "Mistral Response",
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
			want:   "Anthropic Response",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnvVars != nil {
				tt.setEnvVars()
			}
			if tt.unsetEnvVars != nil {
				defer tt.unsetEnvVars()
			}

			got, err := tt.llm.GenerateText(context.Background(), tt.prompt)
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
