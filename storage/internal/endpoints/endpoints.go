package endpoints

import (
	"fmt"
	"strings"
)

func GetAccountNameFromEndpoint(endpoint string) (*string, error) {
	segments := strings.Split(endpoint, ".")
	if len(segments) == 0 {
		return nil, fmt.Errorf("The Endpoint contained no segments")
	}
	return &segments[0], nil
}

// GetOrBuildBlobEndpoint returns the endpoint for Blob API Operations on this storage account
func GetOrBuildBlobEndpoint(endpoint, baseUri string, accountName string) string {
	if endpoint != "" {
		return endpoint
	}
	return fmt.Sprintf("https://%s.blob.%s", accountName, baseUri)
}

// GetOrBuildDataLakeStoreEndpoint returns the endpoint for Data Lake Store API Operations on this storage account
func GetOrBuildDataLakeStoreEndpoint(endpoint, baseUri string, accountName string) string {
	if endpoint != "" {
		return endpoint
	}
	return fmt.Sprintf("https://%s.dfs.%s", accountName, baseUri)
}

// GetOrBuildFileEndpoint returns the endpoint for File Share API Operations on this storage account
func GetOrBuildFileEndpoint(endpoint, baseUri string, accountName string) string {
	if endpoint != "" {
		return endpoint
	}
	return fmt.Sprintf("https://%s.file.%s", accountName, baseUri)
}

// GetOrBuildQueueEndpoint returns the endpoint for Queue API Operations on this storage account
func GetOrBuildQueueEndpoint(endpoint, baseUri string, accountName string) string {
	if endpoint != "" {
		return endpoint
	}
	return fmt.Sprintf("https://%s.queue.%s", accountName, baseUri)
}

// GetOrBuildTableEndpoint returns the endpoint for Table API Operations on this storage account
func GetOrBuildTableEndpoint(endpoint, baseUri string, accountName string) string {
	if endpoint != "" {
		return endpoint
	}
	return fmt.Sprintf("https://%s.table.%s", accountName, baseUri)
}
