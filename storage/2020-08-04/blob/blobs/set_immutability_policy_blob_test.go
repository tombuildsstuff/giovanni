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

var _ StorageBlob = Client{}

func TestBlobPolicyLifecycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "example.txt"

	testData, err := client.BuildTestResourcesWithImmutability(ctx, resourceGroup, accountName, storage.KindBlobStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth := auth.NewSharedKeyLiteAuthorizer(accountName, testData.StorageAccountKey)
	containersClient := containers.NewWithEnvironment(client.Environment)
	containersClient.Client = client.PrepareWithAuthorizer(containersClient.Client, storageAuth)

	_, err = containersClient.Create(ctx, accountName, containerName, containers.CreateInput{})
	if err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}
	defer containersClient.Delete(ctx, accountName, containerName)

	blobClient := NewWithEnvironment(client.Environment)
	blobClient.Client = client.PrepareWithAuthorizer(blobClient.Client, storageAuth)

	t.Logf("[DEBUG] Copying file to Blob Storage..")
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-server-i386.iso.torrent",
	}

	refreshInterval := 5 * time.Second
	if err := blobClient.CopyAndWait(ctx, accountName, containerName, fileName, copyInput, refreshInterval); err != nil {
		t.Fatalf("Error copying: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Blob Properties..")
	details, err := blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error retrieving properties: %s", err)
	}

	t.Logf("[DEBUG] Checking it's returned in the List API..")
	listInput := containers.ListBlobsInput{}
	listResult, err := containersClient.ListBlobs(ctx, accountName, containerName, listInput)
	if err != nil {
		t.Fatalf("Error listing blobs: %s", err)
	}

	if len(listResult.Blobs.Blobs) != 1 {
		t.Fatalf("Expected there to be 1 blob in the container but got %d", len(listResult.Blobs.Blobs))
	}

	gmtTimeLoc := time.FixedZone("GMT", 0)
	untilDate := time.Now().Add(24 * time.Hour).In(gmtTimeLoc).Format(time.RFC1123)

	t.Logf("[DEBUG]: Setting immutability policy")
	immutabilityPolicy := ImmutabilityPolicyBlobInput{
		PolicyMode: ImmutabilityPolicyModeUnlocked,
		UntilDate:  untilDate,
	}
	_, err = blobClient.SetImmutabilityPolicyBlob(ctx, accountName, containerName, fileName, immutabilityPolicy)
	if err != nil {
		t.Fatalf("Error setting immutability policy: %s", err)
	}

	t.Logf("[DEBUG] Re-retrieving Blob Properties..")
	details, err = blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error re-retrieving properties: %s", err)
	}

	if details.ImmutabilityPolicyMode != "unlocked" {
		t.Fatalf("Expected immutability policy mode to be `unlocked`, but got %q", details.ImmutabilityPolicyMode)
	}

	if details.ImmutabilityPolicyUntilDate != untilDate {
		t.Fatalf("Expected immutability policy untilDate to be `%s`, but got %q", untilDate, details.ImmutabilityPolicyUntilDate)
	}

	t.Logf("[DEBUG] Deleting immutability policy")
	_, err = blobClient.DeleteImmutabilityPolicyBlob(ctx, accountName, containerName, fileName)
	if err != nil {
		t.Fatalf("Error deleting immutability policy: %s", err)
	}

	t.Logf("[DEBUG] Deleting Blob")
	if _, err := blobClient.Delete(ctx, accountName, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting Blob: %s", err)
	}
}
