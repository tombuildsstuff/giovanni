package messages

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageQueueMessage interface {
	Delete(ctx context.Context, accountName, queueName, messageID, popReceipt string) (result autorest.Response, err error)
	Peek(ctx context.Context, accountName, queueName string, numberOfMessages int) (result QueueMessagesListResult, err error)
	GetResourceID(accountName, queueName, messageID string) string
	Put(ctx context.Context, accountName, queueName string, input PutInput) (result QueueMessagesListResult, err error)
	Get(ctx context.Context, accountName, queueName string, numberOfMessages int, input GetInput) (result QueueMessagesListResult, err error)
	Update(ctx context.Context, accountName, queueName string, messageID string, input UpdateInput) (result autorest.Response, err error)
}
