package queues

import (
	"context"
)

type StorageQueue interface {
	Create(ctx context.Context, accountName, queueName string, metaData map[string]string)
	Delete(ctx context.Context, accountName, queueName string)
	GetMetaData(ctx context.Context, accountName, queueName string)
	SetMetaData(ctx context.Context, accountName, queueName string, metaData map[string]string)
	GetServiceProperties(ctx context.Context, accountName string)
	SetServiceProperties(ctx context.Context, accountName string, properties StorageServiceProperties)
}
