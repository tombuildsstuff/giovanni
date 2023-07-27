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

// BuildBlobEndpoint returns the endpoint for Blob API Operations on this storage account
func BuildBlobEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.blob.%s", accountName, baseUri)
}

// BuildDataLakeStoreEndpoint returns the endpoint for Data Lake Store API Operations on this storage account
func BuildDataLakeStoreEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.dfs.%s", accountName, baseUri)
}

// BuildFileEndpoint returns the endpoint for File Share API Operations on this storage account
func BuildFileEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.file.%s", accountName, baseUri)
}

// BuildQueueEndpoint returns the endpoint for Queue API Operations on this storage account
func BuildQueueEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.queue.%s", accountName, baseUri)
}

// BuildTableEndpoint returns the endpoint for Table API Operations on this storage account
func BuildTableEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.table.%s", accountName, baseUri)
}
