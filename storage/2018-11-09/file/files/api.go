package files

import (
	"context"
	"os"
	"time"
)

type StorageFile interface {
	Copy(ctx context.Context, accountName, shareName, path, fileName string, input CopyInput)
	AbortCopy(ctx context.Context, accountName, shareName, path, fileName, copyID string)
	CopyAndWait(ctx context.Context, accountName, shareName, path, fileName string, input CopyInput, pollDuration time.Duration)
	Create(ctx context.Context, accountName, shareName, path, fileName string, input CreateInput)
	Delete(ctx context.Context, accountName, shareName, path, fileName string)
	GetMetaData(ctx context.Context, accountName, shareName, path, fileName string)
	SetMetaData(ctx context.Context, accountName, shareName, path, fileName string, metaData map[string]string)
	GetProperties(ctx context.Context, accountName, shareName, path, fileName string)
	SetProperties(ctx context.Context, accountName, shareName, path, fileName string, input SetPropertiesInput)
	ClearByteRange(ctx context.Context, accountName, shareName, path, fileName string, input ClearByteRangeInput)
	GetByteRange(ctx context.Context, accountName, shareName, path, fileName string, input GetByteRangeInput)
	GetFile(ctx context.Context, accountName, shareName, path, fileName string, parallelism int)
	PutByteRange(ctx context.Context, accountName, shareName, path, fileName string, input PutByteRangeInput)
	PutFile(ctx context.Context, accountName, shareName, path, fileName string, file *os.File, parallelism int)
	ListRanges(ctx context.Context, accountName, shareName, path, fileName string)
}