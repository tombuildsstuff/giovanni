package blobs

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-03-28/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestAppendBlobLifecycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "append-blob.txt"

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

	t.Logf("[DEBUG] Putting Append Blob..")
	if _, err := blobClient.PutAppendBlob(ctx, accountName, containerName, fileName, PutAppendBlobInput{}); err != nil {
		t.Fatalf("Error putting append blob: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties..")
	props, err := blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error retrieving properties: %s", err)
	}
	if props.ContentLength != 0 {
		t.Fatalf("Expected Content-Length to be 0 but it was %d", props.ContentLength)
	}

	t.Logf("[DEBUG] Appending First Block..")
	appendInput := AppendBlockInput{
		Content: &[]byte{
			12,
			48,
			93,
			76,
			29,
			10,
		},
	}
	if _, err := blobClient.AppendBlock(ctx, accountName, containerName, fileName, appendInput); err != nil {
		t.Fatalf("Error appending first block: %s", err)
	}

	t.Logf("[DEBUG] Re-Retrieving Properties..")
	props, err = blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error retrieving properties: %s", err)
	}
	if props.ContentLength != 6 {
		t.Fatalf("Expected Content-Length to be 6 but it was %d", props.ContentLength)
	}

	t.Logf("[DEBUG] Appending Second Block..")
	appendInput = AppendBlockInput{
		Content: &[]byte{
			92,
			62,
			64,
			47,
			83,
			77,
		},
	}
	if _, err := blobClient.AppendBlock(ctx, accountName, containerName, fileName, appendInput); err != nil {
		t.Fatalf("Error appending Second block: %s", err)
	}

	t.Logf("[DEBUG] Re-Retrieving Properties..")
	props, err = blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error retrieving properties: %s", err)
	}
	if props.ContentLength != 12 {
		t.Fatalf("Expected Content-Length to be 12 but it was %d", props.ContentLength)
	}

	t.Logf("[DEBUG] Acquiring Lease..")
	leaseDetails, err := blobClient.AcquireLease(ctx, accountName, containerName, fileName, AcquireLeaseInput{
		LeaseDuration: -1,
	})
	if err != nil {
		t.Fatalf("Error acquiring Lease: %s", err)
	}
	t.Logf("[DEBUG] Lease ID is %q", leaseDetails.LeaseID)

	t.Logf("[DEBUG] Appending Third Block..")
	appendInput = AppendBlockInput{
		Content: &[]byte{
			64,
			35,
			28,
			93,
			11,
			23,
		},
		LeaseID: &leaseDetails.LeaseID,
	}
	if _, err := blobClient.AppendBlock(ctx, accountName, containerName, fileName, appendInput); err != nil {
		t.Fatalf("Error appending Third block: %s", err)
	}

	t.Logf("[DEBUG] Re-Retrieving Properties..")
	props, err = blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{
		LeaseID: &leaseDetails.LeaseID,
	})
	if err != nil {
		t.Fatalf("Error retrieving properties: %s", err)
	}
	if props.ContentLength != 18 {
		t.Fatalf("Expected Content-Length to be 18 but it was %d", props.ContentLength)
	}

	t.Logf("[DEBUG] Breaking Lease..")
	breakLeaseInput := BreakLeaseInput{
		LeaseID: leaseDetails.LeaseID,
	}
	if _, err := blobClient.BreakLease(ctx, accountName, containerName, fileName, breakLeaseInput); err != nil {
		t.Fatalf("Error breaking lease: %s", err)
	}

	t.Logf("[DEBUG] Deleting Lease..")
	if _, err := blobClient.Delete(ctx, accountName, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
