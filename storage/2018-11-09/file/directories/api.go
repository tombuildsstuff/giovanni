package directories

import (
	"context"
)

type StorageDirectory interface {
	Create(ctx context.Context, accountName, shareName, path string, metaData map[string]string)
	Delete(ctx context.Context, accountName, shareName, path string)
	Get(ctx context.Context, accountName, shareName, path string)
	GetMetaData(ctx context.Context, accountName, shareName, path string)
	SetMetaData(ctx context.Context, accountName, shareName, path string, metaData map[string]string)
}