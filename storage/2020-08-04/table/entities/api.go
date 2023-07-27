package entities

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageTableEntity interface {
	Delete(ctx context.Context, tableName string, input DeleteEntityInput) (result autorest.Response, err error)
	Insert(ctx context.Context, tableName string, input InsertEntityInput) (result autorest.Response, err error)
	InsertOrReplace(ctx context.Context, tableName string, input InsertOrReplaceEntityInput) (result autorest.Response, err error)
	InsertOrMerge(ctx context.Context, tableName string, input InsertOrMergeEntityInput) (result autorest.Response, err error)
	Query(ctx context.Context, tableName string, input QueryEntitiesInput) (result QueryEntitiesResult, err error)
	Get(ctx context.Context, tableName string, input GetEntityInput) (result GetEntityResult, err error)
	GetResourceID(tableName, partitionKey, rowKey string) string
}
