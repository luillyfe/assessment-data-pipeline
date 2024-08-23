package main

import (
	"context"
	"log"
	"os"
	"reflect"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
	"github.com/luillyfe/assessment-data-pipeline/firestoreio"
)

type Assessment struct {
	Result string `firestore:"assessment_result"`
}

// From each assessment document extract its Result property
func ProcessElement(doc Assessment, emit func(string)) {
	emit(doc.Result)
}

func init() {
	beam.RegisterType(reflect.TypeOf((*Assessment)(nil)).Elem())
	beam.RegisterFunction(ProcessElement)
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

func transformData(scope beam.Scope, documents beam.PCollection) beam.PCollection {
	// Process the Firestore documents
	return beam.ParDo(scope, ProcessElement, documents)
}

func loadDataIntoDestination(scope beam.Scope, processed beam.PCollection) {
	// Write the processed data to the destination
	textio.Write(scope, "processed.txt", processed)
}
