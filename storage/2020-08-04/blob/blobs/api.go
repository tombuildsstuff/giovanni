package blobs

import (
	"context"
	"os"
	"time"
)

type StorageBlob interface {
	AppendBlock(ctx context.Context, containerName string, blobName string, input AppendBlockInput) (resp AppendBlockResponse, err error)
	Copy(ctx context.Context, containerName string, blobName string, input CopyInput) (resp CopyResponse, err error)
	AbortCopy(ctx context.Context, containerName string, blobName string, input AbortCopyInput) (resp CopyAbortResponse, err error)
	CopyAndWait(ctx context.Context, containerName string, blobName string, input CopyInput, pollingInterval time.Duration) error
	Delete(ctx context.Context, containerName string, blobName string, input DeleteInput) (resp DeleteResponse, err error)
	DeleteSnapshot(ctx context.Context, containerName string, blobName string, input DeleteSnapshotInput) (resp DeleteSnapshotResponse, err error)
	DeleteSnapshots(ctx context.Context, containerName string, blobName string, input DeleteSnapshotsInput) (resp DeleteSnapshotsResponse, err error)
	Get(ctx context.Context, containerName string, blobName string, input GetInput) (resp GetResponse, err error)
	GetBlockList(ctx context.Context, containerName string, blobName string, input GetBlockListInput) (resp GetBlockListResponse, err error)
	GetPageRanges(ctx context.Context, accountName, containerName, blobName string, input GetPageRangesInput) (result GetPageRangesResponse, err error)
	IncrementalCopyBlob(ctx context.Context, containerName string, blobName string, input IncrementalCopyBlobInput) (resp IncrementalCopyBlob, err error)
	AcquireLease(ctx context.Context, containerName string, blobName string, input AcquireLeaseInput) (resp AcquireLeaseResponse, err error)
	BreakLease(ctx context.Context, containerName string, blobName string, input BreakLeaseInput) (resp BreakLeaseResponse, err error)
	ChangeLease(ctx context.Context, containerName string, blobName string, input ChangeLeaseInput) (resp ChangeLeaseResponse, err error)
	ReleaseLease(ctx context.Context, containerName string, blobName string, input ReleaseLeaseInput) (resp ReleaseLeaseResponse, err error)
	RenewLease(ctx context.Context, containerName string, blobName string, input RenewLeaseInput) (resp RenewLeaseResponse, err error)
	SetMetaData(ctx context.Context, containerName string, blobName string, input SetMetaDataInput) (resp SetMetaDataResponse, err error)
	GetProperties(ctx context.Context, containerName string, blobName string, input GetPropertiesInput) (resp GetPropertiesResponse, err error)
	SetProperties(ctx context.Context, containerName string, blobName string, input SetPropertiesInput) (resp SetPropertiesResponse, err error)
	PutAppendBlob(ctx context.Context, containerName string, blobName string, input PutAppendBlobInput) (resp PutAppendBlobResponse, err error)
	PutBlock(ctx context.Context, containerName string, blobName string, input PutBlockInput) (resp PutBlockResponse, err error)
	PutBlockBlob(ctx context.Context, containerName string, blobName string, input PutBlockBlobInput) (resp PutBlockBlobResponse, err error)
	PutBlockBlobFromFile(ctx context.Context, containerName string, blobName string, file *os.File, input PutBlockBlobInput) error
	PutBlockList(ctx context.Context, containerName string, blobName string, input PutBlockListInput) (resp PutBlockListResponse, err error)
	PutBlockFromURL(ctx context.Context, containerName string, blobName string, input PutBlockFromURLInput) (resp PutBlockFromURLResponse, err error)
	PutPageBlob(ctx context.Context, containerName string, blobName string, input PutPageBlobInput) (resp PutPageBlobResponse, err error)
	PutPageClear(ctx context.Context, containerName string, blobName string, input PutPageClearInput) (resp PutPageClearResponse, err error)
	PutPageUpdate(ctx context.Context, containerName string, blobName string, input PutPageUpdateInput) (resp PutPageUpdateResponse, err error)
	SetTier(ctx context.Context, containerName string, blobName string, input SetTierInput) (resp SetTierResponse, err error)
	Snapshot(ctx context.Context, containerName string, blobName string, input SnapshotInput) (resp SnapshotResponse, err error)
	GetSnapshotProperties(ctx context.Context, containerName string, blobName string, input GetSnapshotPropertiesInput) (resp GetPropertiesResponse, err error)
	Undelete(ctx context.Context, containerName string, blobName string) (resp UndeleteResponse, err error)
}
