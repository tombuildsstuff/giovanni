package shares

import (
	"context"
)

type StorageShare interface {
	GetACL(ctx context.Context, accountName, shareName string)
	SetACL(ctx context.Context, accountName, shareName string, acls []SignedIdentifier)
	Create(ctx context.Context, accountName, shareName string, input CreateInput)
	Delete(ctx context.Context, accountName, shareName string, deleteSnapshots bool)
	GetMetaData(ctx context.Context, accountName, shareName string)
	SetMetaData(ctx context.Context, accountName, shareName string, metaData map[string]string)
	GetProperties(ctx context.Context, accountName, shareName string)
	SetProperties(ctx context.Context, accountName, shareName string, newQuotaGB int)
	CreateSnapshot(ctx context.Context, accountName, shareName string, input CreateSnapshotInput)
	DeleteSnapshot(ctx context.Context, accountName, shareName string, shareSnapshot string)
	GetSnapshot(ctx context.Context, accountName, shareName, snapshotShare string)
}