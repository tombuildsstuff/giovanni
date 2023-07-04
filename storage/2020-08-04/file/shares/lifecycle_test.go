package shares

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var _ StorageShare = Client{}

func TestSharesLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.KindStorageV2)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLite)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
	sharesClient := NewWithEnvironment(client.AutoRestEnvironment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	tier := CoolAccessTier
	input := CreateInput{
		QuotaInGB:  1,
		AccessTier: &tier,
	}
	_, err = sharesClient.Create(ctx, accountName, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}

	snapshot, err := sharesClient.CreateSnapshot(ctx, accountName, shareName, CreateSnapshotInput{})
	if err != nil {
		t.Fatalf("Error taking snapshot: %s", err)
	}
	t.Logf("Snapshot Date Time: %s", snapshot.SnapshotDateTime)

	snapshotDetails, err := sharesClient.GetSnapshot(ctx, accountName, shareName, snapshot.SnapshotDateTime)
	if err != nil {
		t.Fatalf("Error retrieving snapshot: %s", err)
	}

	t.Logf("MetaData: %s", snapshotDetails.MetaData)

	_, err = sharesClient.DeleteSnapshot(ctx, accountName, shareName, snapshot.SnapshotDateTime)
	if err != nil {
		t.Fatalf("Error deleting snapshot: %s", err)
	}

	stats, err := sharesClient.GetStats(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving stats: %s", err)
	}

	if stats.ShareUsageBytes != 0 {
		t.Fatalf("Expected `stats.ShareUsageBytes` to be 0 but got: %d", stats.ShareUsageBytes)
	}

	share, err := sharesClient.GetProperties(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving share: %s", err)
	}
	if share.QuotaInGB != 1 {
		t.Fatalf("Expected Quota to be 1 but got: %d", share.QuotaInGB)
	}
	if share.EnabledProtocol != SMB {
		t.Fatalf("Expected EnabledProtocol to SMB but got: %s", share.EnabledProtocol)
	}
	if share.AccessTier == nil || *share.AccessTier != CoolAccessTier {
		t.Fatalf("Expected AccessTier to be Cool but got: %v", share.AccessTier)
	}

	newTier := HotAccessTier
	quota := 5
	props := ShareProperties{
		AccessTier: &newTier,
		QuotaInGb:  &quota,
	}
	_, err = sharesClient.SetProperties(ctx, accountName, shareName, props)
	if err != nil {
		t.Fatalf("Error updating quota: %s", err)
	}

	share, err = sharesClient.GetProperties(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving share: %s", err)
	}
	if share.QuotaInGB != 5 {
		t.Fatalf("Expected Quota to be 5 but got: %d", share.QuotaInGB)
	}

	if share.AccessTier == nil || *share.AccessTier != HotAccessTier {
		t.Fatalf("Expected AccessTier to be Hot but got: %v", share.AccessTier)
	}

	updatedMetaData := map[string]string{
		"hello": "world",
	}
	_, err = sharesClient.SetMetaData(ctx, accountName, shareName, updatedMetaData)
	if err != nil {
		t.Fatalf("Erorr setting metadata: %s", err)
	}

	result, err := sharesClient.GetMetaData(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving metadata: %s", err)
	}

	if result.MetaData["hello"] != "world" {
		t.Fatalf("Expected metadata `hello` to be `world` but got: %q", result.MetaData["hello"])
	}
	if len(result.MetaData) != 1 {
		t.Fatalf("Expected metadata to be 1 item but got: %s", result.MetaData)
	}

	acls, err := sharesClient.GetACL(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving ACL's: %s", err)
	}
	if len(acls.SignedIdentifiers) != 0 {
		t.Fatalf("Expected 0 identifiers but got %d", len(acls.SignedIdentifiers))
	}

	updatedAcls := []SignedIdentifier{
		{
			Id: "abc123",
			AccessPolicy: AccessPolicy{
				Start:      "2020-07-01T08:49:37.0000000Z",
				Expiry:     "2020-07-01T09:49:37.0000000Z",
				Permission: "rwd",
			},
		},
		{
			Id: "bcd234",
			AccessPolicy: AccessPolicy{
				Start:      "2020-07-01T08:49:37.0000000Z",
				Expiry:     "2020-07-01T09:49:37.0000000Z",
				Permission: "rwd",
			},
		},
	}
	_, err = sharesClient.SetACL(ctx, accountName, shareName, updatedAcls)
	if err != nil {
		t.Fatalf("Error setting ACL's: %s", err)
	}

	acls, err = sharesClient.GetACL(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving ACL's: %s", err)
	}
	if len(acls.SignedIdentifiers) != 2 {
		t.Fatalf("Expected 2 identifiers but got %d", len(acls.SignedIdentifiers))
	}

	_, err = sharesClient.Delete(ctx, accountName, shareName, false)
	if err != nil {
		t.Fatalf("Error deleting Share: %s", err)
	}
}

func TestSharesLifecycleLargeQuota(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResourcesWithSku(ctx, resourceGroup, accountName, storage.KindFileStorage, storage.SkuNamePremiumLRS)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLite)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
	sharesClient := NewWithEnvironment(client.AutoRestEnvironment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	input := CreateInput{
		QuotaInGB: 1001,
	}
	_, err = sharesClient.Create(ctx, accountName, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}

	snapshot, err := sharesClient.CreateSnapshot(ctx, accountName, shareName, CreateSnapshotInput{})
	if err != nil {
		t.Fatalf("Error taking snapshot: %s", err)
	}
	t.Logf("Snapshot Date Time: %s", snapshot.SnapshotDateTime)

	snapshotDetails, err := sharesClient.GetSnapshot(ctx, accountName, shareName, snapshot.SnapshotDateTime)
	if err != nil {
		t.Fatalf("Error retrieving snapshot: %s", err)
	}

	t.Logf("MetaData: %s", snapshotDetails.MetaData)

	_, err = sharesClient.DeleteSnapshot(ctx, accountName, shareName, snapshot.SnapshotDateTime)
	if err != nil {
		t.Fatalf("Error deleting snapshot: %s", err)
	}

	stats, err := sharesClient.GetStats(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving stats: %s", err)
	}

	if stats.ShareUsageBytes != 0 {
		t.Fatalf("Expected `stats.ShareUsageBytes` to be 0 but got: %d", stats.ShareUsageBytes)
	}

	share, err := sharesClient.GetProperties(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving share: %s", err)
	}
	if share.QuotaInGB != 1001 {
		t.Fatalf("Expected Quota to be 1001 but got: %d", share.QuotaInGB)
	}

	newQuota := 6000
	props := ShareProperties{
		QuotaInGb: &newQuota,
	}
	_, err = sharesClient.SetProperties(ctx, accountName, shareName, props)
	if err != nil {
		t.Fatalf("Error updating quota: %s", err)
	}

	share, err = sharesClient.GetProperties(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving share: %s", err)
	}
	if share.QuotaInGB != 6000 {
		t.Fatalf("Expected Quota to be 6000 but got: %d", share.QuotaInGB)
	}

	updatedMetaData := map[string]string{
		"hello": "world",
	}
	_, err = sharesClient.SetMetaData(ctx, accountName, shareName, updatedMetaData)
	if err != nil {
		t.Fatalf("Erorr setting metadata: %s", err)
	}

	result, err := sharesClient.GetMetaData(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving metadata: %s", err)
	}

	if result.MetaData["hello"] != "world" {
		t.Fatalf("Expected metadata `hello` to be `world` but got: %q", result.MetaData["hello"])
	}
	if len(result.MetaData) != 1 {
		t.Fatalf("Expected metadata to be 1 item but got: %s", result.MetaData)
	}

	acls, err := sharesClient.GetACL(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving ACL's: %s", err)
	}
	if len(acls.SignedIdentifiers) != 0 {
		t.Fatalf("Expected 0 identifiers but got %d", len(acls.SignedIdentifiers))
	}

	updatedAcls := []SignedIdentifier{
		{
			Id: "abc123",
			AccessPolicy: AccessPolicy{
				Start:      "2020-07-01T08:49:37.0000000Z",
				Expiry:     "2020-07-01T09:49:37.0000000Z",
				Permission: "rwd",
			},
		},
		{
			Id: "bcd234",
			AccessPolicy: AccessPolicy{
				Start:      "2020-07-01T08:49:37.0000000Z",
				Expiry:     "2020-07-01T09:49:37.0000000Z",
				Permission: "rwd",
			},
		},
	}
	_, err = sharesClient.SetACL(ctx, accountName, shareName, updatedAcls)
	if err != nil {
		t.Fatalf("Error setting ACL's: %s", err)
	}

	acls, err = sharesClient.GetACL(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving ACL's: %s", err)
	}
	if len(acls.SignedIdentifiers) != 2 {
		t.Fatalf("Expected 2 identifiers but got %d", len(acls.SignedIdentifiers))
	}

	_, err = sharesClient.Delete(ctx, accountName, shareName, false)
	if err != nil {
		t.Fatalf("Error deleting Share: %s", err)
	}
}

func TestSharesLifecycleNFSProtocol(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	
	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResourcesWithSku(ctx, resourceGroup, accountName, storage.KindFileStorage, storage.SkuNamePremiumLRS)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLite)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
	sharesClient := NewWithEnvironment(client.AutoRestEnvironment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	input := CreateInput{
		QuotaInGB:       1000,
		EnabledProtocol: NFS,
	}
	_, err = sharesClient.Create(ctx, accountName, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}

	share, err := sharesClient.GetProperties(ctx, accountName, shareName)
	if err != nil {
		t.Fatalf("Error retrieving share: %s", err)
	}
	if share.EnabledProtocol != NFS {
		t.Fatalf(`Expected enabled protocol to be "NFS" but got: %q`, share.EnabledProtocol)
	}

	_, err = sharesClient.Delete(ctx, accountName, shareName, false)
	if err != nil {
		t.Fatalf("Error deleting Share: %s", err)
	}
}
