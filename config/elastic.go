package config

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

// ElasticsearchConfig represents Elasticsearch configuration
type ElasticsearchConfig struct {
	UserName       string
	Password       string
	Addresses      string
	CertificateKey string
}

// BuildElasticsearchConfig builds Elasticsearch config object from environment variables
func BuildElasticsearchConfig() ElasticsearchConfig {
	elasticsearchConfig := ElasticsearchConfig{
		Addresses:      "https://127.0.0.1:9200",
		UserName:       "elastic",
		Password:       "cAiuR5fpUvukS84Sbbbt",
		CertificateKey: "afc84c6422a7b160096a52219c6e16ee07fc43cd61d8e8dfce2c63f1a46fb14d",
	}

	return elasticsearchConfig
}

// GetElasticsearchClient creates an Elasticsearch client based on the configuration
func GetElasticsearchClient() (*elasticsearch.Client, error) {
	elasticsearchConfig := BuildElasticsearchConfig()

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:              []string{elasticsearchConfig.Addresses},
		Username:               elasticsearchConfig.UserName,
		Password:               elasticsearchConfig.Password,
		CertificateFingerprint: elasticsearchConfig.CertificateKey,
	})

	if err != nil {
		log.Println("GetElasticsearchClient: Failed to create Elasticsearch client with :", err)
		return nil, err
	}

	return esClient, nil
}
