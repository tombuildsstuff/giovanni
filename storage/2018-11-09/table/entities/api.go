package entities

import (
	"context"
)

type StorageTableEntity interface {
	Delete(ctx context.Context, accountName, tableName string, input DeleteEntityInput)
	Get(ctx context.Context, accountName, tableName string, input GetEntityInput)
	Insert(ctx context.Context, accountName, tableName string, input InsertEntityInput)
	InsertOrMerge(ctx context.Context, accountName, tableName string, input InsertOrMergeEntityInput)
	InsertOrReplace(ctx context.Context, accountName, tableName string, input InsertOrReplaceEntityInput)
	Query(ctx context.Context, accountName, tableName string, input QueryEntitiesInput)
}
