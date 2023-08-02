package blobs

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestSnapshotLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "example.txt"

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.KindBlobStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}

	containersClient, err := containers.NewWithBaseUri(fmt.Sprintf("https://%s.blob.%s", testData.StorageAccountName, *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(containersClient.Client, testData); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	_, err = containersClient.Create(ctx, containerName, containers.CreateInput{})
	if err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}
	defer containersClient.Delete(ctx, containerName)

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLite)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
	blobClient := NewWithEnvironment(client.AutoRestEnvironment)
	blobClient.Client = client.PrepareWithAuthorizer(blobClient.Client, storageAuth)

	t.Logf("[DEBUG] Copying file to Blob Storage..")
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
	}

	refreshInterval := 5 * time.Second
	if err := blobClient.CopyAndWait(ctx, accountName, containerName, fileName, copyInput, refreshInterval); err != nil {
		t.Fatalf("Error copying: %s", err)
	}

	t.Logf("[DEBUG] First Snapshot..")
	firstSnapshot, err := blobClient.Snapshot(ctx, accountName, containerName, fileName, SnapshotInput{})
	if err != nil {
		t.Fatalf("Error taking first snapshot: %s", err)
	}
	t.Logf("[DEBUG] First Snapshot ID: %q", firstSnapshot.SnapshotDateTime)

	t.Log("[DEBUG] Waiting 2 seconds..")
	time.Sleep(2 * time.Second)

	t.Logf("[DEBUG] Second Snapshot..")
	secondSnapshot, err := blobClient.Snapshot(ctx, accountName, containerName, fileName, SnapshotInput{
		MetaData: map[string]string{
			"hello": "world",
		},
	})
	if err != nil {
		t.Fatalf("Error taking Second snapshot: %s", err)
	}
	t.Logf("[DEBUG] Second Snapshot ID: %q", secondSnapshot.SnapshotDateTime)

	t.Logf("[DEBUG] Leasing the Blob..")
	leaseDetails, err := blobClient.AcquireLease(ctx, accountName, containerName, fileName, AcquireLeaseInput{
		// infinite
		LeaseDuration: -1,
	})
	if err != nil {
		t.Fatalf("Error leasing Blob: %s", err)
	}
	t.Logf("[DEBUG] Lease ID: %q", leaseDetails.LeaseID)

	t.Logf("[DEBUG] Third Snapshot..")
	thirdSnapshot, err := blobClient.Snapshot(ctx, accountName, containerName, fileName, SnapshotInput{
		LeaseID: &leaseDetails.LeaseID,
	})
	if err != nil {
		t.Fatalf("Error taking Third snapshot: %s", err)
	}
	t.Logf("[DEBUG] Third Snapshot ID: %q", thirdSnapshot.SnapshotDateTime)

	t.Logf("[DEBUG] Releasing Lease..")
	if _, err := blobClient.ReleaseLease(ctx, accountName, containerName, fileName, leaseDetails.LeaseID); err != nil {
		t.Fatalf("Error releasing Lease: %s", err)
	}

	// get the properties from the blob, which should include the LastModifiedDate
	t.Logf("[DEBUG] Retrieving Properties for Blob")
	props, err := blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error getting properties: %s", err)
	}

	// confirm that the If-Modified-None returns an error
	t.Logf("[DEBUG] Third Snapshot..")
	fourthSnapshot, err := blobClient.Snapshot(ctx, accountName, containerName, fileName, SnapshotInput{
		LeaseID:         &leaseDetails.LeaseID,
		IfModifiedSince: &props.LastModified,
	})
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
	if fourthSnapshot.Response.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("Expected the status code to be Precondition Failed but got: %d", fourthSnapshot.Response.StatusCode)
	}

	t.Logf("[DEBUG] Retrieving the Second Snapshot Properties..")
	getSecondSnapshotInput := GetSnapshotPropertiesInput{
		SnapshotID: secondSnapshot.SnapshotDateTime,
	}
	if _, err := blobClient.GetSnapshotProperties(ctx, accountName, containerName, fileName, getSecondSnapshotInput); err != nil {
		t.Fatalf("Error retrieving properties for the second snapshot: %s", err)
	}

	t.Logf("[DEBUG] Deleting the Second Snapshot..")
	deleteSnapshotInput := DeleteSnapshotInput{
		SnapshotDateTime: secondSnapshot.SnapshotDateTime,
	}
	if _, err := blobClient.DeleteSnapshot(ctx, accountName, containerName, fileName, deleteSnapshotInput); err != nil {
		t.Fatalf("Error deleting snapshot: %s", err)
	}

	t.Logf("[DEBUG] Re-Retrieving the Second Snapshot Properties..")
	secondSnapshotProps, err := blobClient.GetSnapshotProperties(ctx, accountName, containerName, fileName, getSecondSnapshotInput)
	if err == nil {
		t.Fatalf("Expected an error retrieving the snapshot but got none")
	}
	if secondSnapshotProps.Response.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected the status code to be %d but got %q", http.StatusNoContent, secondSnapshotProps.Response.StatusCode)
	}

	t.Logf("[DEBUG] Deleting all the snapshots..")
	if _, err := blobClient.DeleteSnapshots(ctx, accountName, containerName, fileName, DeleteSnapshotsInput{}); err != nil {
		t.Fatalf("Error deleting snapshots: %s", err)
	}

	t.Logf("[DEBUG] Deleting the Blob..")
	deleteInput := DeleteInput{
		DeleteSnapshots: false,
	}
	if _, err := blobClient.Delete(ctx, accountName, containerName, fileName, deleteInput); err != nil {
		t.Fatalf("Error deleting Blob: %s", err)
	}
}
