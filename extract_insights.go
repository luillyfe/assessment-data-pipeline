package main

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/luillyfe/assessment-data-pipeline/llm"
)

// ExtractInsights is a DoFn that extract insights from user's performance.
type ExtractInsights struct {
	model llm.LanguageModel
}

// ProcessElement sends a request to the LLM to ask for extracting key insights from user performance.
func (a *ExtractInsights) ProcessElement(ctx context.Context, assessments Assessment, emit func(string)) {
	// Read file output schema
	insights_schema, err := readFile("insights_schema.json")
	if err != nil {
		log.Printf("Error reading insights schema: %v", err)
		emit("")
		return
	}

	prompt := fmt.Sprintf("Dear Gemini, given the following assessment from an user performance on the Professional Data Engineer Certification Prep:\n %s. Please extract key insights from it. And respond in the following json schema:\n %v. Remove any ```json and ```. Neither comments no explanations", assessments.Result, insights_schema)
	text, err := a.model.GenerateText(
		context.Background(),
		prompt,
		&llm.GenerateOptions{
			ResponseMIMEType: "application/json",
		},
	)
	// Handle error gracefully without terminating the program
	if err != nil {
		log.Printf("Error reading insights schema: %v", err)
		emit("")
		return
	}

	emit(text)
}

func init() {
	register.DoFn3x0[context.Context, Assessment, func(string)](&ExtractInsights{})
}

func (a *ExtractInsights) Setup() {
	a.model = llm.NewGeminiClient(llm.WithMaxTokens(8192))
}
