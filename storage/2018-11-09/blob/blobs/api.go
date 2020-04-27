package blobs

import (
	"context"
	"os"
	"time"
)

type StorageBlob interface {
	AppendBlock(ctx context.Context, accountName, containerName, blobName string, input AppendBlockInput)
	AbortCopy(ctx context.Context, accountName, containerName, blobName string, input AbortCopyInput)
	CopyAndWait(ctx context.Context, accountName, containerName, blobName string, input CopyInput, pollingInterval time.Duration)
	Delete(ctx context.Context, accountName, containerName, blobName string, input DeleteInput)
	DeleteSnapshot(ctx context.Context, accountName, containerName, blobName string, input DeleteSnapshotInput)
	DeleteSnapshots(ctx context.Context, accountName, containerName, blobName string, input DeleteSnapshotsInput)
	Get(ctx context.Context, accountName, containerName, blobName string, input GetInput)
	GetBlockList(ctx context.Context, accountName, containerName, blobName string, input GetBlockListInput)
	GetPageRanges(ctx context.Context, accountName, containerName, blobName string, input GetPageRangesInput)
	IncrementalCopyBlob(ctx context.Context, accountName, containerName, blobName string, input IncrementalCopyBlobInput)
	AcquireLease(ctx context.Context, accountName, containerName, blobName string, input AcquireLeaseInput)
	BreakLease(ctx context.Context, accountName, containerName, blobName string, input BreakLeaseInput)
	ChangeLease(ctx context.Context, accountName, containerName, blobName string, input ChangeLeaseInput)
	ReleaseLease(ctx context.Context, accountName, containerName, blobName, leaseID string)
	RenewLease(ctx context.Context, accountName, containerName, blobName, leaseID string)
	SetMetaData(ctx context.Context, accountName, containerName, blobName string, input SetMetaDataInput)
	GetProperties(ctx context.Context, accountName, containerName, blobName string, input GetPropertiesInput)
	SetProperties(ctx context.Context, accountName, containerName, blobName string, input SetPropertiesInput)
	PutAppendBlob(ctx context.Context, accountName, containerName, blobName string, input PutAppendBlobInput)
	PutBlock(ctx context.Context, accountName, containerName, blobName string, input PutBlockInput)
	PutBlockBlob(ctx context.Context, accountName, containerName, blobName string, input PutBlockBlobInput)
	PutBlockBlobFromFile(ctx context.Context, accountName, containerName, blobName string, file *os.File, input PutBlockBlobInput)
	PutBlockList(ctx context.Context, accountName, containerName, blobName string, input PutBlockListInput)
	PutBlockFromURL(ctx context.Context, accountName, containerName, blobName string, input PutBlockFromURLInput)
	PutPageBlob(ctx context.Context, accountName, containerName, blobName string, input PutPageBlobInput)
	PutPageClear(ctx context.Context, accountName, containerName, blobName string, input PutPageClearInput)
	PutPageUpdate(ctx context.Context, accountName, containerName, blobName string, input PutPageUpdateInput)
	SetTier(ctx context.Context, accountName, containerName, blobName string, tier AccessTier)
	Snapshot(ctx context.Context, accountName, containerName, blobName string, input SnapshotInput)
	GetSnapshotProperties(ctx context.Context, accountName, containerName, blobName string, input GetSnapshotPropertiesInput)
	Undelete(ctx context.Context, accountName, containerName, blobName string)
}