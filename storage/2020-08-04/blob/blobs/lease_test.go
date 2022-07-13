package blobs

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestLeaseLifecycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "ubuntu.iso"

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.KindBlobStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	containersClient := containers.NewWithEnvironment(client.Environment)
	containersClient.Client = client.PrepareWithStorageResourceManagerAuth(containersClient.Client)

	_, err = containersClient.Create(ctx, accountName, containerName, containers.CreateInput{})
	if err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}
	defer containersClient.Delete(ctx, accountName, containerName)

	storageAuth := auth.NewSharedKeyLiteAuthorizer(accountName, testData.StorageAccountKey)
	blobClient := NewWithEnvironment(client.Environment)
	blobClient.Client = client.PrepareWithAuthorizer(blobClient.Client, storageAuth)

	t.Logf("[DEBUG] Copying file to Blob Storage..")
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
	}

	refreshInterval := 5 * time.Second
	if err := blobClient.CopyAndWait(ctx, accountName, containerName, fileName, copyInput, refreshInterval); err != nil {
		t.Fatalf("Error copying: %s", err)
	}
	defer blobClient.Delete(ctx, accountName, containerName, fileName, DeleteInput{})

	// Test begins here
	t.Logf("[DEBUG] Acquiring Lease..")
	leaseInput := AcquireLeaseInput{
		LeaseDuration: -1,
	}
	leaseInfo, err := blobClient.AcquireLease(ctx, accountName, containerName, fileName, leaseInput)
	if err != nil {
		t.Fatalf("Error acquiring lease: %s", err)
	}
	t.Logf("[DEBUG] Lease ID: %q", leaseInfo.LeaseID)

	t.Logf("[DEBUG] Changing Lease..")
	changeLeaseInput := ChangeLeaseInput{
		ExistingLeaseID: leaseInfo.LeaseID,
		ProposedLeaseID: "31f5bb01-cdd9-4166-bcdc-95186076bde0",
	}
	changeLeaseResult, err := blobClient.ChangeLease(ctx, accountName, containerName, fileName, changeLeaseInput)
	if err != nil {
		t.Fatalf("Error changing lease: %s", err)
	}
	t.Logf("[DEBUG] New Lease ID: %q", changeLeaseResult.LeaseID)

	t.Logf("[DEBUG] Releasing Lease..")
	if _, err := blobClient.ReleaseLease(ctx, accountName, containerName, fileName, changeLeaseResult.LeaseID); err != nil {
		t.Fatalf("Error releasing lease: %s", err)
	}

	t.Logf("[DEBUG] Acquiring a new lease..")
	leaseInput = AcquireLeaseInput{
		LeaseDuration: 30,
	}
	leaseInfo, err = blobClient.AcquireLease(ctx, accountName, containerName, fileName, leaseInput)
	if err != nil {
		t.Fatalf("Error acquiring lease: %s", err)
	}
	t.Logf("[DEBUG] Lease ID: %q", leaseInfo.LeaseID)

	t.Logf("[DEBUG] Renewing lease..")
	if _, err := blobClient.RenewLease(ctx, accountName, containerName, fileName, leaseInfo.LeaseID); err != nil {
		t.Fatalf("Error renewing lease: %s", err)
	}

	t.Logf("[DEBUG] Breaking lease..")
	breakLeaseInput := BreakLeaseInput{
		LeaseID: leaseInfo.LeaseID,
	}
	if _, err := blobClient.BreakLease(ctx, accountName, containerName, fileName, breakLeaseInput); err != nil {
		t.Fatalf("Error breaking lease: %s", err)
	}
}
