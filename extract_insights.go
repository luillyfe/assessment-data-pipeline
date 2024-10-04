package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/luillyfe/assessment-data-pipeline/llm"
)

// ExtractInsights is a DoFn that extracts insights from user's performance.
type ExtractInsights struct {
	model          llm.LanguageModel
	InsightsSchema string
	MaxRetries     int
	RetryDelay     time.Duration
}

// InsightsResult represents the structure of the extracted insights.
type InsightsResult struct {
	OverallAssessment  string            `json:"overall_assessment"`
	CorrectAnswers     int               `json:"questions_answered_correctly"`
	Strengths          []string          `json:"strengths"`
	Weaknesses         []string          `json:"weaknesses"`
	ActionableFeedback map[string]string `json:"actionable_feedback"`
	BusinessImpact     map[string]string `json:"business_case_impact_analysis"`
}

// ProcessElement sends a request to the LLM to extract key insights from user performance.
func (ei *ExtractInsights) ProcessElement(ctx context.Context, assessment Assessment, emit func(InsightsResult)) {
	var (
		insights InsightsResult
		err      error
	)

	for attempt := 0; attempt < ei.MaxRetries; attempt++ {
		insights, err = ei.extractInsights(ctx, assessment)
		if err == nil {
			emit(insights)
			return
		}

		log.Printf("Attempt %d failed: %v. Retrying...", attempt+1, err)
		time.Sleep(ei.RetryDelay)
	}

	log.Printf("Failed to extract insights after %d attempts: %v", ei.MaxRetries, err)
}

func (ei *ExtractInsights) extractInsights(ctx context.Context, assessment Assessment) (InsightsResult, error) {
	prompt := fmt.Sprintf("Given the following assessment from a user's performance on the Professional Data Engineer Certification Prep:\n%s\nPlease extract key insights and respond in the following JSON schema:\n%s . Remove any ```json or ``` characters. Avoid any comments or explanations", assessment.Result, ei.InsightsSchema)

	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	text, err := ei.model.GenerateText(
		ctx,
		prompt,
		&llm.GenerateOptions{
			ResponseMIMEType: "application/json",
		},
	)
	if err != nil {
		return InsightsResult{}, fmt.Errorf("error generating text: %w", err)
	}

	var insights InsightsResult
	if err := json.Unmarshal([]byte(text), &insights); err != nil {
		return InsightsResult{}, fmt.Errorf("error unmarshaling insights: %w", err)
	}

	return insights, nil
}

func (ei *ExtractInsights) Setup() error {
	var err error
	ei.InsightsSchema, err = readFile("insights_schema.json")
	if err != nil {
		return fmt.Errorf("error reading insights schema: %w", err)
	}
	ei.model = llm.NewGeminiClient(llm.WithMaxTokens(8192))
	return nil
}

func init() {
	register.DoFn3x0[context.Context, Assessment, func(InsightsResult)](&ExtractInsights{})
	register.Function2x1(NewExtractInsights)
	beam.RegisterType(reflect.TypeOf((*InsightsResult)(nil)).Elem())
}

// NewExtractInsights creates a new ExtractInsights DoFn with custom retry settings.
func NewExtractInsights(maxRetries int, retryDelay time.Duration) *ExtractInsights {
	return &ExtractInsights{
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
	}
}
