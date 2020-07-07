package directories

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

var StorageFile = Client{}

func TestDirectoriesLifeCycle(t *testing.T) {
	client, err := testhelpers.Build()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.Storage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth := auth.NewSharedKeyLiteAuthorizer(accountName, testData.StorageAccountKey)
	sharesClient := shares.NewWithEnvironment(client.Environment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	directoriesClient := NewWithEnvironment(client.Environment)
	directoriesClient.Client = client.PrepareWithAuthorizer(directoriesClient.Client, storageAuth)

	input := shares.CreateInput{
		QuotaInGB: 1,
	}
	_, err = sharesClient.Create(ctx, accountName, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, accountName, shareName, true)

	metaData := map[string]string{
		"hello": "world",
	}

	log.Printf("[DEBUG] Creating Top Level..")
	if _, err := directoriesClient.Create(ctx, accountName, shareName, "hello", metaData); err != nil {
		t.Fatalf("Error creating Top Level Directory: %s", err)
	}

	log.Printf("[DEBUG] Creating Inner..")
	if _, err := directoriesClient.Create(ctx, accountName, shareName, "hello/there", metaData); err != nil {
		t.Fatalf("Error creating Inner Directory: %s", err)
	}

	log.Printf("[DEBUG] Retrieving share")
	innerDir, err := directoriesClient.Get(ctx, accountName, shareName, "hello/there")
	if err != nil {
		t.Fatalf("Error retrieving Inner Directory: %s", err)
	}

	if innerDir.DirectoryMetaDataEncrypted != true {
		t.Fatalf("Expected MetaData to be encrypted but got: %t", innerDir.DirectoryMetaDataEncrypted)
	}

	if len(innerDir.MetaData) != 1 {
		t.Fatalf("Expected MetaData to contain 1 item but got %d", len(innerDir.MetaData))
	}
	if innerDir.MetaData["hello"] != "world" {
		t.Fatalf("Expected MetaData `hello` to be `world`: %s", innerDir.MetaData["hello"])
	}

	log.Printf("[DEBUG] Setting MetaData")
	updatedMetaData := map[string]string{
		"panda": "pops",
	}
	if _, err := directoriesClient.SetMetaData(ctx, accountName, shareName, "hello/there", updatedMetaData); err != nil {
		t.Fatalf("Error updating MetaData: %s", err)
	}

	log.Printf("[DEBUG] Retrieving MetaData")
	retrievedMetaData, err := directoriesClient.GetMetaData(ctx, accountName, shareName, "hello/there")
	if err != nil {
		t.Fatalf("Error retrieving the updated metadata: %s", err)
	}
	if len(retrievedMetaData.MetaData) != 1 {
		t.Fatalf("Expected the updated metadata to have 1 item but got %d", len(retrievedMetaData.MetaData))
	}
	if retrievedMetaData.MetaData["panda"] != "pops" {
		t.Fatalf("Expected the metadata `panda` to be `pops` but got %q", retrievedMetaData.MetaData["panda"])
	}

	t.Logf("[DEBUG] Deleting Inner..")
	if _, err := directoriesClient.Delete(ctx, accountName, shareName, "hello/there"); err != nil {
		t.Fatalf("Error deleting Inner Directory: %s", err)
	}

	t.Logf("[DEBUG] Deleting Top Level..")
	if _, err := directoriesClient.Delete(ctx, accountName, shareName, "hello"); err != nil {
		t.Fatalf("Error deleting Top Level Directory: %s", err)
	}
}
