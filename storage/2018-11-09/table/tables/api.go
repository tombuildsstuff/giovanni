package tables

import (
	"context"
)

type StorageTable interface {
	GetACL(ctx context.Context, accountName, tableName string)
	SetACL(ctx context.Context, accountName, tableName string, acls []SignedIdentifier)
	Create(ctx context.Context, accountName, tableName string)
	Delete(ctx context.Context, accountName, tableName string)
	Exists(ctx context.Context, accountName, tableName string)
	Query(ctx context.Context, accountName string, metaDataLevel MetaDataLevel)
}
