package messages

import (
	"context"
)

type StorageQueueMessage interface {
	Delete(ctx context.Context, accountName, queueName, messageID, popReceipt string)
	Get(ctx context.Context, accountName, queueName string, numberOfMessages int, input GetInput)
	Peek(ctx context.Context, accountName, queueName string, numberOfMessages int)
	Put(ctx context.Context, accountName, queueName string, input PutInput)
	Update(ctx context.Context, accountName, queueName string, messageID string, input UpdateInput)
}
