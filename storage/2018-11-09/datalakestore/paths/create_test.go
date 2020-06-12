package paths

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestCreateDirectory(t *testing.T) {
	client, err := testhelpers.Build()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	fileSystemName := fmt.Sprintf("acctestfs-%s", testhelpers.RandomString())
	path := "test"

	if _, err = client.BuildTestResources(ctx, resourceGroup, accountName, storage.BlobStorage); err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	fileSystemsClient := filesystems.NewWithEnvironment(client.Environment)
	fileSystemsClient.Client = client.PrepareWithStorageResourceManagerAuth(fileSystemsClient.Client)

	t.Logf("[DEBUG] Creating an empty File System..")
	fileSystemInput := filesystems.CreateInput{
		Properties: map[string]string{},
	}
	if _, err = fileSystemsClient.Create(ctx, accountName, fileSystemName, fileSystemInput); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Creating path..")
	pathsClient := NewWithEnvironment(client.Environment)
	pathsClient.Client = client.PrepareWithStorageResourceManagerAuth(pathsClient.Client)

	input := CreateInput{
		Resource: PathResourceDirectory,
	}

	if _, err = pathsClient.Create(ctx, accountName, fileSystemName, path, input); err != nil {
		t.Fatal(fmt.Errorf("Error creating path: %s", err))
	}

	t.Logf("[DEBUG] Deleting File System..")
	if _, err := fileSystemsClient.Delete(ctx, accountName, fileSystemName); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
