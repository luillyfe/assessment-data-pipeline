package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
	"github.com/luillyfe/assessment-data-pipeline/firestoreio"
)

type Assessment struct {
	Result string `firestore:"assessment_result"`
}

func init() {
	beam.RegisterType(reflect.TypeOf((*Assessment)(nil)).Elem())
	beam.RegisterFunction(insightsToJSON)
}

func main() {
	// Handling os-environment variables
	projectID, assessmentCollection := handleOSEnvironmentVariables()

	// Initialize Beam
	beam.Init()

	// Create a new Beam pipeline
	pipeline, scope := beam.NewPipelineWithRoot()

	// Reading data from the source
	documents := readDataFromSource(scope, projectID, assessmentCollection)

	// Transforming the data
	processed := transformData(scope, documents)

	// Loading the data into the destination
	loadDataIntoDestination(scope, processed)

	// Run the pipeline
	if err := beamx.Run(context.Background(), pipeline); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}
}

func handleOSEnvironmentVariables() (string, string) {
	// Parse os-environment variables
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("Please set the GOOGLE_CLOUD_PROJECT environment variable.")
	}

	assessmentCollection := os.Getenv("ASSESSMENT_COLLECTION")
	if assessmentCollection == "" {
		log.Fatal("Please set the ASSESSMENT_COLLECTION environment variable.")
	}

	// Return the values of the flags
	return projectID, assessmentCollection
}

func readDataFromSource(scope beam.Scope, project, assessmentCollection string) beam.PCollection {
	// Define the ReadConfig
	cfg := firestoreio.ReadConfig{
		Project:    project,
		Collection: assessmentCollection,
	}

	// Define the element type
	elemType := reflect.TypeOf(Assessment{})

	// Read data from the source using firestoreio.Read
	return firestoreio.Read(scope, cfg, elemType)
}

func transformData(scope beam.Scope, assessments beam.PCollection) beam.PCollection {
	extractInsights := NewExtractInsights(3, 10*time.Second)
	// Process the Firestore documents
	return beam.ParDo(scope, extractInsights, assessments)
}

// insightsToJSON converts InsightsResult to JSON string
func insightsToJSON(insight InsightsResult) string {
	jsonBytes, err := json.Marshal(insight)
	if err != nil {
		log.Printf("Error marshaling insight to JSON: %v", err)
		return ""
	}
	return string(jsonBytes)
}

func loadDataIntoDestination(scope beam.Scope, processed beam.PCollection) {
	// Convert insights to JSON strings
	jsonInsights := beam.ParDo(scope, insightsToJSON, processed)
	// Write the processed data to the destination
	textio.Write(scope, "processed.jsonl", jsonInsights)
}
