package queues

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageQueue interface {
	Delete(ctx context.Context, queueName string) (result autorest.Response, err error)
	GetMetaData(ctx context.Context, queueName string) (result GetMetaDataResult, err error)
	SetMetaData(ctx context.Context, queueName string, metaData map[string]string) (result autorest.Response, err error)
	Create(ctx context.Context, queueName string, metaData map[string]string) (result autorest.Response, err error)
	GetResourceID(queueName string) string
	SetServiceProperties(ctx context.Context, properties StorageServiceProperties) (result autorest.Response, err error)
	GetServiceProperties(ctx context.Context) (result StorageServicePropertiesResponse, err error)
}
