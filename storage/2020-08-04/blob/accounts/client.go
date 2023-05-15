package accounts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane/storage"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Client is the base client for Blob Storage Blobs.
type Client struct {
	Client *storage.BaseClient
}

func NewWithEnvironment(accountName string, environment environments.Environment) (*Client, error) {
	baseClient, err := storage.NewBaseClient(accountName, "blob", environment.Storage, serviceName, apiVersion)
	if err != nil {
		return nil, fmt.Errorf("building base client: %+v", err)
	}
	return &Client{
		Client: baseClient,
	}, nil
}
