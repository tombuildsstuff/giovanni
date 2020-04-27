package containers

import "context"

type StorageContainer interface {
	Create(ctx context.Context, accountName, containerName string, input CreateInput)
	Delete(ctx context.Context, accountName, containerName string)
	GetProperties(ctx context.Context, accountName, containerName string)
	GetPropertiesWithLeaseID(ctx context.Context, accountName, containerName, leaseID string)
	AcquireLease(ctx context.Context, accountName, containerName string, input AcquireLeaseInput)
	BreakLease(ctx context.Context, accountName, containerName string, input BreakLeaseInput)
	ChangeLease(ctx context.Context, accountName, containerName string, input ChangeLeaseInput)
	ReleaseLease(ctx context.Context, accountName, containerName, leaseID string)
	RenewLease(ctx context.Context, accountName, containerName, leaseID string)
	ListBlobs(ctx context.Context, accountName, containerName string, input ListBlobsInput)
	SetAccessControl(ctx context.Context, accountName, containerName string, level AccessLevel)
	SetAccessControlWithLeaseID(ctx context.Context, accountName, containerName, leaseID string, level AccessLevel)
	SetMetaData(ctx context.Context, accountName, containerName string, metaData map[string]string)
	SetMetaDataWithLeaseID(ctx context.Context, accountName, containerName, leaseID string, metaData map[string]string)
}
