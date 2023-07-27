package filesystems

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// Client is the base client for Data Lake Storage FileSystem
type Client struct {
	autorest.Client
	endpoint string
}

// New creates an instance of the Data Lake Storage FileSystem client.
func New(accountName string) Client {
	return NewWithEnvironment(accountName, azure.PublicCloud)
}

// NewWithEnvironment creates an instance of the Data Lake Storage FileSystem client.
func NewWithEnvironment(accountName string, environment azure.Environment) Client {
	return Client{
		Client:   autorest.NewClientWithUserAgent(UserAgent()),
		endpoint: endpoints.BuildDataLakeStoreEndpoint(environment.StorageEndpointSuffix, accountName),
	}
}

// NewWithEndpoint creates an instance of the Client client with the endpoint specified.
func NewWithEndpoint(endpoint string) Client {
	return Client{
		Client:   autorest.NewClientWithUserAgent(UserAgent()),
		endpoint: endpoint,
	}
}
