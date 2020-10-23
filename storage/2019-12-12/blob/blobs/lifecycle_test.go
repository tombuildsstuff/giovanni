package blobs

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

var _ StorageBlob = Client{}

func TestLifecycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "example.txt"

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.BlobStorage)
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
		CopySource: "http://releases.ubuntu.com/18.04.2/ubuntu-18.04.2-desktop-amd64.iso",
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

	// default value
	if details.AccessTier != Hot {
		t.Fatalf("Expected the AccessTier to be %q but got %q", Hot, details.AccessTier)
	}
	if details.BlobType != BlockBlob {
		t.Fatalf("Expected BlobType to be %q but got %q", BlockBlob, details.BlobType)
	}
	if len(details.MetaData) != 0 {
		t.Fatalf("Expected there to be no items of metadata but got %d", len(details.MetaData))
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

	t.Logf("[DEBUG] Setting MetaData..")
	metaDataInput := SetMetaDataInput{
		MetaData: map[string]string{
			"hello": "there",
		},
	}
	if _, err := blobClient.SetMetaData(ctx, accountName, containerName, fileName, metaDataInput); err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	t.Logf("[DEBUG] Re-retrieving Blob Properties..")
	details, err = blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error re-retrieving properties: %s", err)
	}

	// default value
	if details.AccessTier != Hot {
		t.Fatalf("Expected the AccessTier to be %q but got %q", Hot, details.AccessTier)
	}
	if details.BlobType != BlockBlob {
		t.Fatalf("Expected BlobType to be %q but got %q", BlockBlob, details.BlobType)
	}
	if len(details.MetaData) != 1 {
		t.Fatalf("Expected there to be 1 item of metadata but got %d", len(details.MetaData))
	}
	if details.MetaData["hello"] != "there" {
		t.Fatalf("Expected `hello` to be `there` but got %q", details.MetaData["there"])
	}

	t.Logf("[DEBUG] Retrieving the Block List..")
	getBlockListInput := GetBlockListInput{
		BlockListType: All,
	}
	blockList, err := blobClient.GetBlockList(ctx, accountName, containerName, fileName, getBlockListInput)
	if err != nil {
		t.Fatalf("Error retrieving Block List: %s", err)
	}

	// since this is a copy from an existing file, all blocks should be present
	if len(blockList.CommittedBlocks.Blocks) == 0 {
		t.Fatalf("Expected there to be committed blocks but there weren't!")
	}
	if len(blockList.UncommittedBlocks.Blocks) != 0 {
		t.Fatalf("Expected all blocks to be committed but got %d uncommitted blocks", len(blockList.UncommittedBlocks.Blocks))
	}

	t.Logf("[DEBUG] Changing the Access Tiers..")
	tiers := []AccessTier{
		Hot,
		Cool,
		Archive,
	}
	for _, tier := range tiers {
		t.Logf("[DEBUG] Updating the Access Tier to %q..", string(tier))
		if _, err := blobClient.SetTier(ctx, accountName, containerName, fileName, tier); err != nil {
			t.Fatalf("Error setting the Access Tier: %s", err)
		}

		t.Logf("[DEBUG] Re-retrieving Blob Properties..")
		details, err = blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
		if err != nil {
			t.Fatalf("Error re-retrieving properties: %s", err)
		}

		if details.AccessTier != tier {
			t.Fatalf("Expected the AccessTier to be %q but got %q", tier, details.AccessTier)
		}
	}

	t.Logf("[DEBUG] Deleting Blob")
	if _, err := blobClient.Delete(ctx, accountName, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting Blob: %s", err)
	}
}
