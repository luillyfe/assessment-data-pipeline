package llm

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

type geminiLLM struct {
	modelName   string
	temperature float64
	maxTokens   int
	topP        float64
	client      *genai.Client
}

/*
GenerateText generates text using the Google Gemini LLM based on the provided prompt and optional generation options.

It takes a context.Context, a prompt string, and optional generation options as input.
It constructs a Gemini chat request with the prompt and model parameters.
It sends the request to the Gemini API using the client.
It handles potential errors from the Gemini API.
It extracts and returns the generated text from the API response.

Args:

	ctx: The context for the request.
	prompt: The input prompt for text generation.
	opts: Optional generation options, such as tools.

Returns:

	A string containing the generated text and an error if any occurred.
*/
func (g *geminiLLM) GenerateText(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	// Model initialization
	model := g.client.GenerativeModel(g.modelName)

	// Model configuration
	model.SetTemperature(float32(g.temperature))
	model.SetTopP(float32(g.topP))
	model.SetMaxOutputTokens(int32(g.maxTokens))
	model.SetTopK(64)
	model.ResponseMIMEType = "text/plain"

	// Tool handling
	if opts != nil && len(opts.Tools) > 0 {
		model.Tools = make([]*genai.Tool, 0)
		for _, genericTool := range opts.Tools {
			if genericTool.Type != GeminiToolType {
				return "", fmt.Errorf("error: tool type mismatch for Gemini LLM")
			}
			geminiTool, ok := genericTool.Tool.(*genai.Tool)
			if !ok {
				return "", fmt.Errorf("error: invalid tool type for Gemini LLM")
			}
			model.Tools = append(model.Tools, geminiTool)
		}
	}

	// Chat session
	session := model.StartChat()
	session.History = []*genai.Content{}

	// Message sending
	resp, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error sending message: %w", err)
	}

	output := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		output = fmt.Sprintf("%v\n", part)
	}

	// Return generated text
	return output, nil
}
