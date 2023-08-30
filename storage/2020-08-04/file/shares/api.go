package shares

import (
	"context"
)

type StorageShare interface {
	SetACL(ctx context.Context, shareName string, input SetAclInput) (resp setAclResponse, err error)
	GetSnapshot(ctx context.Context, shareName string, input GetSnapshotPropertiesInput) (resp GetSnapshotPropertiesResponse, err error)
	GetStats(ctx context.Context, shareName string) (resp GetStatsResponse, err error)
	GetACL(ctx context.Context, shareName string) (resp GetACLResult, err error)
	SetMetaData(ctx context.Context, shareName string, input SetMetaDataInput) (resp SetMetaDataResponse, err error)
	GetMetaData(ctx context.Context, shareName string) (resp GetMetaDataResponse, err error)
	SetProperties(ctx context.Context, shareName string, properties ShareProperties) (resp SetPropertiesResponse, err error)
	DeleteSnapshot(ctx context.Context, accountName string, shareName string, shareSnapshot string) (resp DeleteSnapshotResponse, err error)
	CreateSnapshot(ctx context.Context, shareName string, input CreateSnapshotInput) (resp CreateSnapshotResponse, err error)
	GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, shareName string) string
	GetProperties(ctx context.Context, shareName string) (resp GetPropertiesResult, err error)
	Delete(ctx context.Context, shareName string, input DeleteInput) (resp DeleteResponse, err error)
	Create(ctx context.Context, shareName string, input CreateInput) (resp CreateResponse, err error)
}
