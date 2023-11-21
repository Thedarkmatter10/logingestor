package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"logingestor/config"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response struct defines the structure for the Elasticsearch query response
type Response struct {
	Hits Hits `json:"hits"`
}

// Hits struct encapsulates the array of Hit items
type Hits struct {
	Hits []Hit `json:"hits"`
}

// Hit struct represents a single search result, containing the source log data
type Hit struct {
	Source LogData `json:"_source"`
}

// LogData struct models the log information. Customize fields as per your log format.
type LogData struct {
	Level      string            `json:"level"`
	Message    string            `json:"message"`
	ResourceId string            `json:"resourceId"`
	Timestamp  time.Time         `json:"timestamp"`
	TraceId    string            `json:"traceId"`
	SpanId     string            `json:"spanId"`
	Commit     string            `json:"commit"`
	Metadata   map[string]string `json:"metadata"`
}

// LogIngestHandler handles the ingestion of log data into Elasticsearch
func LogIngestHandler(c *gin.Context) {
	// Initialize Elasticsearch client
	esClient, err := config.GetElasticsearchClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Elasticsearch client"})
		return
	}

	// Bind the incoming JSON payload to the logsData structure
	var logsData []LogData
	if err := c.BindJSON(&logsData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Iterate over each log data item and index it into Elasticsearch
	for _, logData := range logsData {
		jsonData, err := json.Marshal(logData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling log data"})
			return
		}

		// Create a unique document ID for each log entry
		documentID := uuid.New().String()
		req := esapi.IndexRequest{
			Index:      "logs",
			DocumentID: documentID,
			Body:       bytes.NewReader(jsonData),
		}
		res, err := req.Do(context.Background(), esClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to Elasticsearch"})
			return
		}
		defer res.Body.Close()

		// Handle errors in indexing the document
		if res.IsError() {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to index document with ID=%s", documentID)})
			return
		}
	}

	// Return a success response after successful ingestion
	c.JSON(http.StatusOK, gin.H{"message": "Log ingested successfully"})
}

// SearchLogsHandler handles search queries against the Elasticsearch logs index
func SearchLogsHandler(c *gin.Context) {
	// Build the search query based on the request parameters
	query := buildQuery(c)

	// Execute the search query in Elasticsearch
	res, err := searchInElasticsearch(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing search query"})
		return
	}
	defer res.Body.Close()

	// Parse the Elasticsearch response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read search response"})
		return
	}
	var respi Response
	err = json.Unmarshal(body, &respi)
	if err != nil {
		log.Printf("[ERROR] Failed to parse response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse search response"})
		return
	}

	// Send the parsed Elasticsearch response back to the client
	c.JSON(http.StatusOK, respi)
}

// searchInElasticsearch executes the actual search query against Elasticsearch
func searchInElasticsearch(query map[string]interface{}) (*esapi.Response, error) {
	// Initialize Elasticsearch client
	esClient, err := config.GetElasticsearchClient()
	if err != nil {
		log.Printf("[ERROR] Failed to initialize Elasticsearch client: %v", err)
		return nil, err
	}

	// Encode the query into a buffer for the request body
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("[ERROR] Failed to encode query: %v", err)
		return nil, err
	}

	// Execute the search request against Elasticsearch
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("logs"),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("[ERROR] Search query failed: %v", err)
		return nil, err
	}

	return res, nil
}

// buildQuery constructs an Elasticsearch query based on the parameters received from the HTTP request
func buildQuery(c *gin.Context) map[string]interface{} {
	mustClauses := []map[string]interface{}{}

	// Add filter for 'level'
	// Checks if 'level' parameter is present in the request and adds it to the query
	if level := c.Query("level"); level != "" {
		levelQuery := map[string]interface{}{
			"match": map[string]string{"level": level},
		}
		mustClauses = append(mustClauses, levelQuery)
	}

	// Add filter for 'resourceId'
	// Checks if 'resourceId' parameter is present and includes it in the search criteria
	if resourceId := c.Query("resourceId"); resourceId != "" {
		resourceIdQuery := map[string]interface{}{
			"match": map[string]string{"resourceId": resourceId},
		}
		mustClauses = append(mustClauses, resourceIdQuery)
	}

	// Add filter for 'traceId'
	// Adds a match condition for 'traceId' if it's specified in the request
	if traceId := c.Query("traceId"); traceId != "" {
		traceIdQuery := map[string]interface{}{
			"match": map[string]string{"traceId": traceId},
		}
		mustClauses = append(mustClauses, traceIdQuery)
	}

	// Add filter for 'spanId'
	// Includes a condition to filter logs by 'spanId' if provided in the request
	if spanId := c.Query("spanId"); spanId != "" {
		spanIdQuery := map[string]interface{}{
			"match": map[string]string{"spanId": spanId},
		}
		mustClauses = append(mustClauses, spanIdQuery)
	}

	// Add filter for 'commit'
	// Appends a match condition for 'commit' in the query
	if commit := c.Query("commit"); commit != "" {
		commitQuery := map[string]interface{}{
			"match": map[string]string{"spanId": commit},
		}
		mustClauses = append(mustClauses, commitQuery)
	}

	// Add filter for 'message'
	// Filters the logs based on the 'message' field if specified
	if message := c.Query("message"); message != "" {
		messageQuery := map[string]interface{}{
			"match": map[string]string{"message": message},
		}
		mustClauses = append(mustClauses, messageQuery)
	}

	// Add date range filter
	// If 'startDate' and 'endDate' parameters are present, adds a range filter for the 'timestamp' field
	if startDate, endDate := c.Query("startDate"), c.Query("endDate"); startDate != "" && endDate != "" {
		dateRangeQuery := map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": startDate,
					"lte": endDate,
				},
			},
		}
		mustClauses = append(mustClauses, dateRangeQuery)
	}

	// Construct the final Elasticsearch query
	// The query uses a boolean 'must' condition to ensure all specified filters are applied
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
	}

	return query
}
