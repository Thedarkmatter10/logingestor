
## Log Ingestor System 

Project Overview
This project provides a solution for efficiently handling vast volumes of log data. It includes a log ingestor system that receives logs over HTTP and an interface for querying these logs using full-text search or specific field filters.

Getting Started
Prerequisites
Go (Golang) installed on your machine.
Elasticsearch running locally or on a server.
Installation
Clone the repository to your local machine:


git clone [repository URL]
cd [repository directory]
Running the Log Ingestor
To start the log ingestor on port 3000:

go run main.go

---

# Log Ingestor System

## Project Overview

This project provides a solution for efficiently handling vast volumes of log data. It includes a log ingestor system that receives logs over HTTP and an interface for querying these logs using full-text search or specific field filters.

## Getting Started

### Prerequisites

- Go (Golang) installed on your machine.
- Elasticsearch running locally or on a server.

### Installation

Clone the repository to your local machine:

```
git clone [repository URL]
cd [repository directory]
```

### Running the Log Ingestor

To start the log ingestor on port 3000:

```
go run main.go
```

## Usage

### Log Ingestion

Send a POST request to `http://localhost:3000/ingest` with a JSON payload (as mentioned in problem statement) of the log data.
Send a GET  request to `http://localhost:3000/search` with a Query


### Querying Logs

The system offers an interface to query the ingested logs. Use full-text search or specific field filters to retrieve relevant log data.

## Elasticsearch Configuration

Ensure that your Elasticsearch instance is correctly set up and accessible. Configure the connection details in the provided Elasticsearch configuration file.

### Elasticsearch Configuration and Client Setup

This section outlines the necessary steps to configure and use Elasticsearch in our project.

#### Configuration

##### ElasticsearchConfig

We use the `ElasticsearchConfig` struct to hold our Elasticsearch configuration. It includes the following fields:

- `UserName`: The username for Elasticsearch authentication.
- `Password`: The password for Elasticsearch authentication.
- `Addresses`: The URL of the Elasticsearch cluster.
- `CertificateKey`: The certificate fingerprint for secure connections.

##### Default Configuration

Our default configuration is as follows:

- **Address**: `https://127.0.0.1:9200`
- **Username**: `elastic`
- **Password**: `[REDACTED]` (Please replace with your actual password)
- **CertificateKey**: `[REDACTED]` (Please replace with your actual certificate key)

---

**Note**: Replace `[repository URL]` and `[repository directory]` with the actual details of your repository. Adjust the `[REDACTED]` parts in the Elasticsearch configuration with the actual credentials.
