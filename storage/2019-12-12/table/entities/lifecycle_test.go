package entities

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

var _ StorageTableEntity = Client{}

func TestEntitiesLifecycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	tableName := fmt.Sprintf("table%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.Storage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth := auth.NewSharedKeyLiteTableAuthorizer(accountName, testData.StorageAccountKey)
	tablesClient := tables.NewWithEnvironment(client.Environment)
	tablesClient.Client = client.PrepareWithAuthorizer(tablesClient.Client, storageAuth)

	t.Logf("[DEBUG] Creating Table..")
	if _, err := tablesClient.Create(ctx, accountName, tableName); err != nil {
		t.Fatalf("Error creating Table %q: %s", tableName, err)
	}
	defer tablesClient.Delete(ctx, accountName, tableName)

	entitiesClient := NewWithEnvironment(client.Environment)
	entitiesClient.Client = client.PrepareWithAuthorizer(entitiesClient.Client, storageAuth)

	partitionKey := "hello"
	rowKey := "there"

	t.Logf("[DEBUG] Inserting..")
	insertInput := InsertEntityInput{
		MetaDataLevel: NoMetaData,
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
		Entity: map[string]interface{}{
			"hello": "world",
		},
	}
	if _, err := entitiesClient.Insert(ctx, accountName, tableName, insertInput); err != nil {
		t.Logf("Error retrieving: %s", err)
	}

	t.Logf("[DEBUG] Insert or Merging..")
	insertOrMergeInput := InsertOrMergeEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
		Entity: map[string]interface{}{
			"hello": "ther88e",
		},
	}
	if _, err := entitiesClient.InsertOrMerge(ctx, accountName, tableName, insertOrMergeInput); err != nil {
		t.Logf("Error insert/merging: %s", err)
	}

	t.Logf("[DEBUG] Insert or Replacing..")
	insertOrReplaceInput := InsertOrReplaceEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
		Entity: map[string]interface{}{
			"hello": "pandas",
		},
	}
	if _, err := entitiesClient.InsertOrReplace(ctx, accountName, tableName, insertOrReplaceInput); err != nil {
		t.Logf("Error inserting/replacing: %s", err)
	}

	t.Logf("[DEBUG] Querying..")
	queryInput := QueryEntitiesInput{
		MetaDataLevel: NoMetaData,
	}
	results, err := entitiesClient.Query(ctx, accountName, tableName, queryInput)
	if err != nil {
		t.Logf("Error querying: %s", err)
	}

	if len(results.Entities) != 1 {
		t.Fatalf("Expected 1 item but got %d", len(results.Entities))
	}

	for _, v := range results.Entities {
		thisPartitionKey := v["PartitionKey"].(string)
		thisRowKey := v["RowKey"].(string)
		if partitionKey != thisPartitionKey {
			t.Fatalf("Expected Partition Key to be %q but got %q", partitionKey, thisPartitionKey)
		}
		if rowKey != thisRowKey {
			t.Fatalf("Expected Partition Key to be %q but got %q", rowKey, thisRowKey)
		}
	}

	t.Logf("[DEBUG] Retrieving..")
	getInput := GetEntityInput{
		MetaDataLevel: MinimalMetaData,
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
	}
	getResults, err := entitiesClient.Get(ctx, accountName, tableName, getInput)
	if err != nil {
		t.Logf("Error querying: %s", err)
	}

	partitionKey2 := getResults.Entity["PartitionKey"].(string)
	rowKey2 := getResults.Entity["RowKey"].(string)
	if partitionKey2 != partitionKey {
		t.Fatalf("Expected Partition Key to be %q but got %q", partitionKey, partitionKey2)
	}
	if rowKey2 != rowKey {
		t.Fatalf("Expected Row Key to be %q but got %q", rowKey, rowKey2)
	}

	t.Logf("[DEBUG] Deleting..")
	deleteInput := DeleteEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
	}
	if _, err := entitiesClient.Delete(ctx, accountName, tableName, deleteInput); err != nil {
		t.Logf("Error deleting: %s", err)
	}
}
