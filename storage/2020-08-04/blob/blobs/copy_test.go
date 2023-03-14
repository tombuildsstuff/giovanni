package blobs

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestCopyFromExistingFile(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "ubuntu.iso"
	copiedFileName := "copied.iso"

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

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLite)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
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

	t.Logf("[DEBUG] Duplicating that file..")
	copiedInput := CopyInput{
		CopySource: fmt.Sprintf("%s/%s/%s", endpoints.GetBlobEndpoint(blobClient.BaseURI, accountName), containerName, fileName),
	}
	if err := blobClient.CopyAndWait(ctx, accountName, containerName, copiedFileName, copiedInput, refreshInterval); err != nil {
		t.Fatalf("Error duplicating file: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties for the Original File..")
	props, err := blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error getting properties for the original file: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties for the Copied File..")
	copiedProps, err := blobClient.GetProperties(ctx, accountName, containerName, copiedFileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error getting properties for the copied file: %s", err)
	}

	if props.ContentLength != copiedProps.ContentLength {
		t.Fatalf("Expected the content length to be %d but it was %d", props.ContentLength, copiedProps.ContentLength)
	}

	t.Logf("[DEBUG] Deleting copied file..")
	if _, err := blobClient.Delete(ctx, accountName, containerName, copiedFileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting file: %s", err)
	}

	t.Logf("[DEBUG] Deleting original file..")
	if _, err := blobClient.Delete(ctx, accountName, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting file: %s", err)
	}
}

func TestCopyFromURL(t *testing.T) {
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

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLite)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
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

	t.Logf("[DEBUG] Retrieving Properties..")
	props, err := blobClient.GetProperties(ctx, accountName, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error getting properties: %s", err)
	}

	if props.ContentLength == 0 {
		t.Fatalf("Expected the file to be there but looks like it isn't: %d", props.ContentLength)
	}

	t.Logf("[DEBUG] Deleting file..")
	if _, err := blobClient.Delete(ctx, accountName, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting file: %s", err)
	}
}
