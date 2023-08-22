package messages

import (
	"context"
)

type StorageQueueMessage interface {
	Delete(ctx context.Context, queueName string, messageID string, input DeleteInput) (resp DeleteResponse, err error)
	Peek(ctx context.Context, queueName string, input PeekInput) (resp QueueMessagesListResponse, err error)
	GetResourceManagerID(subscriptionID, resourceGroup, accountName, queueName string) string
	Put(ctx context.Context, queueName string, input PutInput) (resp QueueMessagesListResponse, err error)
	Get(ctx context.Context, queueName string, input GetInput) (resp QueueMessagesListResponse, err error)
	Update(ctx context.Context, queueName string, messageID string, input UpdateInput) (resp UpdateResponse, err error)
}
