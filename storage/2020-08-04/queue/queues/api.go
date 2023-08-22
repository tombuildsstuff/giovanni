package queues

import (
	"context"
)

type StorageQueue interface {
	Delete(ctx context.Context, queueName string) (resp DeleteResponse, err error)
	GetMetaData(ctx context.Context, queueName string) (resp GetMetaDataResponse, err error)
	SetMetaData(ctx context.Context, queueName string, input SetMetaDataInput) (resp SetMetaDataResponse, err error)
	Create(ctx context.Context, queueName string, input CreateInput) (resp CreateResponse, err error)
	GetResourceManagerID(subscriptionID, resourceGroup, accountName, queueName string) string
	SetServiceProperties(ctx context.Context, input SetStorageServicePropertiesInput) (resp SetStorageServicePropertiesResponse, err error)
	GetServiceProperties(ctx context.Context) (resp StorageServicePropertiesResponse, err error)
}
