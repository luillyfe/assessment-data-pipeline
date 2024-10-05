package llm

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

type geminiLLM[T ToolType] struct {
	config LLMConfig
	client *genai.Client
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
func (g *geminiLLM[T]) GenerateText(ctx context.Context, prompt string, opts *GenerateOptions[T]) (string, error) {
	// Model initialization
	model := g.client.GenerativeModel(g.config.ModelName)

	// Model configuration
	model.SetTemperature(g.config.Temperature)
	model.SetTopP(g.config.TopP)
	model.SetMaxOutputTokens(int32(g.config.MaxTokens))
	model.SetTopK(64)
	model.ResponseMIMEType = "text/plain" // Default MIME type

	// Tool handling
	if opts != nil && len(opts.Tools) > 0 {
		model.Tools = make([]*genai.Tool, 0)
		for _, opt := range opts.Tools {
			if geminiTool, ok := any(opt.Tool).(genai.Tool); ok {
				model.Tools = append(model.Tools, &geminiTool)
			} else {
				return "", fmt.Errorf("tool doesn't implement genai.Tool")
			}
		}

		// Update ResponseMIMEType if it was set by the caller
		if opts.ResponseMIMEType != "" {
			model.ResponseMIMEType = opts.ResponseMIMEType
		} else {
			model.ResponseMIMEType = "text/plain" // Default MIME type
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
