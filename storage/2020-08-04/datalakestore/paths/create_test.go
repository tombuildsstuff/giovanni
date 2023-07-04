package paths

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestCreateDirectory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	fileSystemName := fmt.Sprintf("acctestfs-%s", testhelpers.RandomString())
	path := "test"

	if _, err = client.BuildTestResourcesWithHns(ctx, resourceGroup, accountName, storage.KindBlobStorage); err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	fileSystemsClient := filesystems.NewWithEnvironment(client.AutoRestEnvironment)
	fileSystemsClient.Client = client.PrepareWithStorageResourceManagerAuth(fileSystemsClient.Client)

	t.Logf("[DEBUG] Creating an empty File System..")
	fileSystemInput := filesystems.CreateInput{
		Properties: map[string]string{},
	}
	if _, err = fileSystemsClient.Create(ctx, accountName, fileSystemName, fileSystemInput); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Creating path..")
	pathsClient := NewWithEnvironment(client.AutoRestEnvironment)
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
