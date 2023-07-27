package containers

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Client is the base client for Blob Storage Containers.
type Client struct {
	autorest.Client
	BaseURI  string
	endpoint string
}

// New creates an instance of the Client client.
func New() Client {
	return NewWithEnvironment(azure.PublicCloud)
}

// NewWithBaseURI creates an instance of the Client client.
func NewWithEnvironment(environment azure.Environment) Client {
	return Client{
		Client:  autorest.NewClientWithUserAgent(UserAgent()),
		BaseURI: environment.StorageEndpointSuffix,
	}
}

// NewWithEndpoint creates an instance of the Client client with the endpoint specified.
func NewWithEndpoint(endpoint string) Client {
	return Client{
		Client:   autorest.NewClientWithUserAgent(UserAgent()),
		endpoint: endpoint,
	}
}

func (client Client) setAccessLevelIntoHeaders(headers map[string]interface{}, level AccessLevel) map[string]interface{} {
	// If this header is not included in the request, container data is private to the account owner.
	if level != Private {
		headers["x-ms-blob-public-access"] = string(level)
	}

	return headers
}
