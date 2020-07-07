package blobs

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2017-07-29/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestPageBlobLifecycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "append-blob.txt"

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.StorageV2)
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

	t.Logf("[DEBUG] Putting Page Blob..")
	fileSize := int64(10240000)
	if _, err := blobClient.PutPageBlob(ctx, accountName, containerName, fileName, PutPageBlobInput{
		BlobContentLengthBytes: fileSize,
	}); err != nil {
		t.Fatalf("Error putting page blob: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties..")
	props, err := blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error retrieving properties: %s", err)
	}
	if props.ContentLength != fileSize {
		t.Fatalf("Expected Content-Length to be %d but it was %d", fileSize, props.ContentLength)
	}

	for iteration := 1; iteration <= 3; iteration++ {
		t.Logf("[DEBUG] Putting Page %d of 3..", iteration)
		byteArray := func() []byte {
			o := make([]byte, 0)

			for i := 0; i < 512; i++ {
				o = append(o, byte(i))
			}

			return o
		}()
		startByte := int64(512 * iteration)
		endByte := int64(startByte + 511)
		putPageInput := PutPageUpdateInput{
			StartByte: startByte,
			EndByte:   endByte,
			Content:   byteArray,
		}
		if _, err := blobClient.PutPageUpdate(ctx, accountName, containerName, fileName, putPageInput); err != nil {
			t.Fatalf("Error putting page: %s", err)
		}
	}

	t.Logf("[DEBUG] Deleting..")
	if _, err := blobClient.Delete(ctx, accountName, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
