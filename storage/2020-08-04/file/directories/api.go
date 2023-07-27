package directories

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageDirectory interface {
	Delete(ctx context.Context, shareName, path string) (result autorest.Response, err error)
	GetMetaData(ctx context.Context, shareName, path string) (result GetMetaDataResult, err error)
	SetMetaData(ctx context.Context, shareName, path string, metaData map[string]string) (result autorest.Response, err error)
	Create(ctx context.Context, shareName, path string, input CreateDirectoryInput) (result autorest.Response, err error)
	GetResourceID(shareName, directoryName string) string
	Get(ctx context.Context, shareName, path string) (result GetResult, err error)
}
