package files

import (
	"context"
	"os"
	"time"
)

type StorageFile interface {
	PutByteRange(ctx context.Context, shareName string, path string, fileName string, input PutByteRangeInput) (resp PutRangeResponse, err error)
	GetByteRange(ctx context.Context, shareName string, path string, fileName string, input GetByteRangeInput) (resp GetByteRangeResponse, err error)
	ClearByteRange(ctx context.Context, shareName string, path string, fileName string, input ClearByteRangeInput) (resp ClearByteRangeResponse, err error)
	SetProperties(ctx context.Context, shareName string, path string, fileName string, input SetPropertiesInput) (resp SetPropertiesResponse, err error)
	PutFile(ctx context.Context, shareName string, path string, fileName string, file *os.File, parallelism int) error
	Copy(ctx context.Context, shareName, path, fileName string, input CopyInput) (resp CopyResponse, err error)
	SetMetaData(ctx context.Context, shareName string, path string, fileName string, input SetMetaDataInput) (resp SetMetaDataResponse, err error)
	GetMetaData(ctx context.Context, shareName string, path string, fileName string) (resp GetMetaDataResponse, err error)
	AbortCopy(ctx context.Context, shareName string, path string, fileName string, input CopyAbortInput) (resp CopyAbortResponse, err error)
	GetFile(ctx context.Context, shareName string, path string, fileName string, input GetFileInput) (resp GetFileResponse, err error)
	ListRanges(ctx context.Context, shareName, path, fileName string) (result ListRangesResponse, err error)
	GetProperties(ctx context.Context, shareName string, path string, fileName string) (resp GetResponse, err error)
	Delete(ctx context.Context, shareName string, path string, fileName string) (resp DeleteResponse, err error)
	Create(ctx context.Context, shareName string, path string, fileName string, input CreateInput) (resp CreateResponse, err error)
	CopyAndWait(ctx context.Context, accountName, shareName, path, fileName string, input CopyInput, pollDuration time.Duration) (result CopyResponse, err error)
}
