# Assessment Data Pipeline

## Overview

This project implements a data pipeline using Apache Beam on Google Cloud Platform (GCP) to process assessment data.

**Key Features:**

- **Efficient Data Processing:** Leverages Apache Beam for scalable and distributed data processing.
- **Cloud-Native:** Designed for deployment on Google Cloud Platform, utilizing GCP services.
- **Configurable:** Utilizes OS environment variables for flexible pipeline configuration.

## Project Structure

└── main.go \
└── firestoreio/ \
└──── read.go \
└──── common.go

- **`main.go`**: Contains the main Go code for the data pipeline, including pipeline setup, data processing logic, and interaction with GCP services.
- **`firestoreio/`**:
  - **`read.go`**: Provides a way to read data from a Firestore collection as part of an Apache Beam pipeline. It handles the integration with Beam's parallel processing capabilities.
  - **`common.go`**: Provides a foundation for the firestoreio package to build upon. It abstracts away common setup, teardown, and configuration details, allowing other files to focus on specific Firestore operations like reading or writing data.

## Getting Started

### Prerequisites

- **Google Cloud Project:** An active GCP project with billing enabled.
- **Google Cloud SDK:** Install and configure the Google Cloud SDK on your local machine.
- **Go Programming Language:** Ensure you have Go installed and configured.
- **Apache Beam SDK for Go:** Install the Apache Beam Go SDK.

### Installation

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/your-username/assessment-data-pipeline.git
   cd assessment-data-pipeline
   ```

### Configuration

1. **Environment Variables:** The pipeline relies on the following environment variables:

   - `GOOGLE_CLOUD_PROJECT`: (Required) The ID of your Google Cloud Project.
   - `ASSESSMENT_COLLECTION`: (Required) The name of the Firestore collection containing the assessment data.

   **Example (Bash):**

   ```bash
   export GOOGLE_CLOUD_PROJECT="your-gcp-project-id"
   export ASSESSMENT_COLLECTION="your-assessment-collection-name"
   go run main.go
   ```

## Data Processing Logic

- Data Ingestion: The pipeline reads assessment data from a Firestore collection.
- Data Transformation: Implement your data processing logic here. This might include:
  - Extracts the "Result" property from each document.
- Data Output: The processed data is written to a text file.

## Acknowledgments

This project utilizes the `firestoreio` module, which is based on the excellent work of Johanna Ojelin. You can find her original repository here: [[Link to Johanna's Repository](https://github.com/johannaojeling/go-beam-pipeline/)]

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have any suggestions or improvements.

## License

This project is licensed under the MIT License.
