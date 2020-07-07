package filesystems

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestCreateHasNoTagsByDefault(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	fileSystemName := fmt.Sprintf("acctestfs-%s", testhelpers.RandomString())

	if _, err = client.BuildTestResources(ctx, resourceGroup, accountName, storage.BlobStorage); err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	fileSystemsClient := NewWithEnvironment(client.Environment)
	fileSystemsClient.Client = client.PrepareWithStorageResourceManagerAuth(fileSystemsClient.Client)

	t.Logf("[DEBUG] Creating an empty File System..")
	input := CreateInput{
		Properties: map[string]string{},
	}
	if _, err = fileSystemsClient.Create(ctx, accountName, fileSystemName, input); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Retrieving the Properties..")
	props, err := fileSystemsClient.GetProperties(ctx, accountName, fileSystemName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error getting properties: %s", err))
	}

	if len(props.Properties) != 0 {
		t.Fatalf("Expected 0 properties by default but got %d", len(props.Properties))
	}

	t.Logf("[DEBUG] Deleting File System..")
	if _, err := fileSystemsClient.Delete(ctx, accountName, fileSystemName); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
