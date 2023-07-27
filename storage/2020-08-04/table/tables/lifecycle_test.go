package tables

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var _ StorageTable = Client{}

func TestTablesLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	tableName := fmt.Sprintf("table%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.KindStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLiteForTable)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
	tablesClient := NewWithEnvironment(accountName, client.AutoRestEnvironment)
	tablesClient.Client = client.PrepareWithAuthorizer(tablesClient.Client, storageAuth)

	t.Logf("[DEBUG] Creating Table..")
	if _, err := tablesClient.Create(ctx, tableName); err != nil {
		t.Fatalf("Error creating Table %q: %s", tableName, err)
	}

	// first look it up directly and confirm it's there
	t.Logf("[DEBUG] Checking if Table exists..")
	if _, err := tablesClient.Exists(ctx, tableName); err != nil {
		t.Fatalf("Error checking if Table %q exists: %s", tableName, err)
	}

	// then confirm it exists in the Query too
	t.Logf("[DEBUG] Querying for Tables..")
	result, err := tablesClient.Query(ctx, NoMetaData)
	if err != nil {
		t.Fatalf("Error retrieving Tables: %s", err)
	}
	found := false
	for _, v := range result.Tables {
		log.Printf("[DEBUG] Table: %q", v.TableName)

		if v.TableName == tableName {
			found = true
		}
	}
	if !found {
		t.Fatalf("%q was not found in the Query response!", tableName)
	}

	t.Logf("[DEBUG] Setting ACL's for Table %q..", tableName)
	acls := []SignedIdentifier{
		{
			Id: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			AccessPolicy: AccessPolicy{
				Permission: "raud",
				Start:      "2020-11-26T08:49:37.0000000Z",
				Expiry:     "2020-11-27T08:49:37.0000000Z",
			},
		},
	}
	if _, err := tablesClient.SetACL(ctx, tableName, acls); err != nil {
		t.Fatalf("Error setting ACLs: %s", err)
	}

	t.Logf("[DEBUG] Retrieving ACL's for Table %q..", tableName)
	retrievedACLs, err := tablesClient.GetACL(ctx, tableName)
	if err != nil {
		t.Fatalf("Error retrieving ACLs: %s", err)
	}

	if len(retrievedACLs.SignedIdentifiers) != len(acls) {
		t.Fatalf("Expected %d but got %q ACLs", len(retrievedACLs.SignedIdentifiers), len(acls))
	}

	for i, retrievedAcl := range retrievedACLs.SignedIdentifiers {
		expectedAcl := acls[i]

		if retrievedAcl.Id != expectedAcl.Id {
			t.Fatalf("Expected ID to be %q but got %q", retrievedAcl.Id, expectedAcl.Id)
		}

		if retrievedAcl.AccessPolicy.Start != expectedAcl.AccessPolicy.Start {
			t.Fatalf("Expected Start to be %q but got %q", retrievedAcl.AccessPolicy.Start, expectedAcl.AccessPolicy.Start)
		}

		if retrievedAcl.AccessPolicy.Expiry != expectedAcl.AccessPolicy.Expiry {
			t.Fatalf("Expected Expiry to be %q but got %q", retrievedAcl.AccessPolicy.Expiry, expectedAcl.AccessPolicy.Expiry)
		}

		if retrievedAcl.AccessPolicy.Permission != expectedAcl.AccessPolicy.Permission {
			t.Fatalf("Expected Permission to be %q but got %q", retrievedAcl.AccessPolicy.Permission, expectedAcl.AccessPolicy.Permission)
		}
	}

	t.Logf("[DEBUG] Deleting Table %q..", tableName)
	if _, err := tablesClient.Delete(ctx, tableName); err != nil {
		t.Fatalf("Error deleting %q: %s", tableName, err)
	}
}
