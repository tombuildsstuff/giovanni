package tables

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageTable interface {
	Delete(ctx context.Context, tableName string) (result autorest.Response, err error)
	Exists(ctx context.Context, tableName string) (result autorest.Response, err error)
	GetACL(ctx context.Context, tableName string) (result GetACLResult, err error)
	Create(ctx context.Context, tableName string) (result autorest.Response, err error)
	GetResourceID(tableName string) string
	Query(ctx context.Context, metaDataLevel MetaDataLevel) (result GetResult, err error)
	SetACL(ctx context.Context, tableName string, acls []SignedIdentifier) (result autorest.Response, err error)
}
