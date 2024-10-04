package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/luillyfe/assessment-data-pipeline/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLanguageModel is a mock implementation of the llm.LanguageModel interface
type MockLanguageModel struct {
	mock.Mock
}

func (m *MockLanguageModel) GenerateText(ctx context.Context, prompt string, opts *llm.GenerateOptions) (string, error) {
	args := m.Called(ctx, prompt, opts)
	return args.String(0), args.Error(1)
}

func TestExtractInsights_ProcessElement(t *testing.T) {
	mockLLM := new(MockLanguageModel)
	ei := &ExtractInsights{
		model:      mockLLM,
		MaxRetries: 3,
		RetryDelay: time.Millisecond,
	}

	testCases := []struct {
		name           string
		assessment     Assessment
		mockResponse   string
		mockError      error
		expectedResult InsightsResult
		expectError    bool
	}{
		{
			name: "Successful extraction",
			assessment: Assessment{
				Result: "User performed well on data engineering questions.",
			},
			mockResponse: `{
				"overall_assessment": "Good performance",
				"questions_answered_correctly": 8,
				"strengths": ["Data modeling", "ETL processes"],
				"weaknesses": ["Cloud security"],
				"actionable_feedback": {"study": "Focus on cloud security concepts"},
				"business_case_impact_analysis": {"efficiency": "Improved data pipeline design"}
			}`,
			expectedResult: InsightsResult{
				OverallAssessment:  "Good performance",
				CorrectAnswers:     8,
				Strengths:          []string{"Data modeling", "ETL processes"},
				Weaknesses:         []string{"Cloud security"},
				ActionableFeedback: map[string]string{"study": "Focus on cloud security concepts"},
				BusinessImpact:     map[string]string{"efficiency": "Improved data pipeline design"},
			},
		},
		{
			name: "LLM error with retry success",
			assessment: Assessment{
				Result: "User struggled with big data concepts.",
			},
			mockResponse: `{
				"overall_assessment": "Needs improvement",
				"questions_answered_correctly": 5,
				"strengths": ["SQL queries"],
				"weaknesses": ["Big data processing", "Data warehousing"],
				"actionable_feedback": {"practice": "Work on Hadoop and Spark exercises"},
				"business_case_impact_analysis": {"cost": "Potential inefficiencies in data processing"}
			}`,
			mockError: errors.New("API error"),
			expectedResult: InsightsResult{
				OverallAssessment:  "Needs improvement",
				CorrectAnswers:     5,
				Strengths:          []string{"SQL queries"},
				Weaknesses:         []string{"Big data processing", "Data warehousing"},
				ActionableFeedback: map[string]string{"practice": "Work on Hadoop and Spark exercises"},
				BusinessImpact:     map[string]string{"cost": "Potential inefficiencies in data processing"},
			},
		},
		{
			name: "Persistent LLM error",
			assessment: Assessment{
				Result: "Assessment data unavailable.",
			},
			mockResponse: `{}`,
			mockError:    errors.New("Persistent API error"),
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockError != nil {
				mockLLM.On("GenerateText", mock.Anything, mock.Anything, mock.Anything).
					Return("", tc.mockError).Once()
				if tc.mockResponse != "" {
					mockLLM.On("GenerateText", mock.Anything, mock.Anything, mock.Anything).
						Return(tc.mockResponse, nil).Once()
				}
			} else {
				mockLLM.On("GenerateText", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.mockResponse, nil).Once()
			}

			var result InsightsResult
			emitFunc := func(insights InsightsResult) {
				result = insights
			}

			ei.ProcessElement(context.Background(), tc.assessment, emitFunc)

			if tc.expectError {
				assert.Equal(t, InsightsResult{}, result)
			} else {
				assert.Equal(t, tc.expectedResult, result)
			}

			mockLLM.AssertExpectations(t)
		})
	}
}

func TestNewExtractInsights(t *testing.T) {
	maxRetries := 5
	retryDelay := 2 * time.Second

	ei := NewExtractInsights(maxRetries, retryDelay)

	assert.Equal(t, maxRetries, ei.MaxRetries)
	assert.Equal(t, retryDelay, ei.RetryDelay)
}

func TestExtractInsights_extractInsights(t *testing.T) {
	mockLLM := new(MockLanguageModel)
	ei := &ExtractInsights{
		model:          mockLLM,
		InsightsSchema: `{"test": "schema"}`,
	}

	testCases := []struct {
		name           string
		assessment     Assessment
		mockResponse   string
		mockError      error
		expectedResult InsightsResult
		expectError    bool
	}{
		{
			name: "Successful extraction",
			assessment: Assessment{
				Result: "User showed proficiency in cloud architecture.",
			},
			mockResponse: `{
				"overall_assessment": "Excellent",
				"questions_answered_correctly": 10,
				"strengths": ["Cloud architecture", "Scalability"],
				"weaknesses": [],
				"actionable_feedback": {"advance": "Explore advanced cloud patterns"},
				"business_case_impact_analysis": {"innovation": "Can lead cloud migration projects"}
			}`,
			expectedResult: InsightsResult{
				OverallAssessment:  "Excellent",
				CorrectAnswers:     10,
				Strengths:          []string{"Cloud architecture", "Scalability"},
				Weaknesses:         []string{},
				ActionableFeedback: map[string]string{"advance": "Explore advanced cloud patterns"},
				BusinessImpact:     map[string]string{"innovation": "Can lead cloud migration projects"},
			},
		},
		{
			name: "LLM error",
			assessment: Assessment{
				Result: "Error occurred during assessment.",
			},
			mockError:   errors.New("LLM API error"),
			expectError: true,
		},
		{
			name: "Invalid JSON response",
			assessment: Assessment{
				Result: "User performance data.",
			},
			mockResponse: `{"invalid": "json"`,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLLM.On("GenerateText", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.mockResponse, tc.mockError).Once()

			result, err := ei.extractInsights(context.Background(), tc.assessment)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}

			mockLLM.AssertExpectations(t)
		})
	}
}
