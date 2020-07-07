package queues

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestQueuesLifecycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	queueName := fmt.Sprintf("queue-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.Storage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth := auth.NewSharedKeyLiteAuthorizer(accountName, testData.StorageAccountKey)
	queuesClient := NewWithEnvironment(client.Environment)
	queuesClient.Client = client.PrepareWithAuthorizer(queuesClient.Client, storageAuth)

	// first let's test an empty container
	_, err = queuesClient.Create(ctx, accountName, queueName, map[string]string{})
	if err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	// then let's retrieve it to ensure there's no metadata..
	resp, err := queuesClient.GetMetaData(ctx, accountName, queueName)
	if err != nil {
		t.Fatalf("Error retrieving MetaData: %s", err)
	}
	if len(resp.MetaData) != 0 {
		t.Fatalf("Expected no MetaData but got: %s", err)
	}

	// then let's add some..
	updatedMetaData := map[string]string{
		"band":  "panic",
		"boots": "the-overpass",
	}
	_, err = queuesClient.SetMetaData(ctx, accountName, queueName, updatedMetaData)
	if err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	resp, err = queuesClient.GetMetaData(ctx, accountName, queueName)
	if err != nil {
		t.Fatalf("Error re-retrieving MetaData: %s", err)
	}

	if len(resp.MetaData) != 2 {
		t.Fatalf("Expected metadata to have 2 items but got: %s", resp.MetaData)
	}
	if resp.MetaData["band"] != "panic" {
		t.Fatalf("Expected `band` to be `panic` but got: %s", resp.MetaData["band"])
	}
	if resp.MetaData["boots"] != "the-overpass" {
		t.Fatalf("Expected `boots` to be `the-overpass` but got: %s", resp.MetaData["boots"])
	}

	// and woo let's remove it again
	_, err = queuesClient.SetMetaData(ctx, accountName, queueName, map[string]string{})
	if err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	resp, err = queuesClient.GetMetaData(ctx, accountName, queueName)
	if err != nil {
		t.Fatalf("Error retrieving MetaData: %s", err)
	}
	if len(resp.MetaData) != 0 {
		t.Fatalf("Expected no MetaData but got: %s", err)
	}

	log.Printf("[DEBUG] Deleting..")
	_, err = queuesClient.Delete(ctx, accountName, queueName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error deleting: %s", err))
	}
}
