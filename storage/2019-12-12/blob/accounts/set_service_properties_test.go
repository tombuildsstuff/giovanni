package accounts

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestContainerLifecycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())

	_, err = client.BuildTestResources(ctx, resourceGroup, accountName, storage.KindStorageV2)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	accountsClient := NewWithEnvironment(client.Environment)
	accountsClient.Client = client.PrepareWithStorageResourceManagerAuth(accountsClient.Client)

	input := StorageServiceProperties{}
	_, err = accountsClient.SetServiceProperties(ctx, accountName, input)
	if err != nil {
		t.Fatal(fmt.Errorf("error setting properties: %s", err))
	}

	var index = "index.html"
	//var enabled = true
	var errorDocument = "404.html"

	input = StorageServiceProperties{
		StaticWebsite: &StaticWebsite{
			Enabled:              true,
			IndexDocument:        index,
			ErrorDocument404Path: errorDocument,
		},
		Logging: &Logging{
			Version: "2.0",
			Delete:  true,
			Read:    true,
			Write:   true,
			RetentionPolicy: DeleteRetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
	}

	_, err = accountsClient.SetServiceProperties(ctx, accountName, input)
	if err != nil {
		t.Fatal(fmt.Errorf("error setting properties: %s", err))
	}

	t.Log("[DEBUG] Waiting 2 seconds..")
	time.Sleep(2 * time.Second)

	result, err := accountsClient.GetServiceProperties(ctx, accountName)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting properties: %s", err))
	}

	website := result.StorageServiceProperties.StaticWebsite
	if website.Enabled != true {
		t.Fatalf("Expected the StaticWebsite %t but got %t", true, website.Enabled)
	}

	logging := result.StorageServiceProperties.Logging
	if logging.Version != "2.0" {
		t.Fatalf("Expected the Logging Version %s but got %s", "2.0", logging.Version)
	}
	if !logging.Read {
		t.Fatalf("Expected the Logging Read %t but got %t", true, logging.Read)
	}
	if !logging.Write {
		t.Fatalf("Expected the Logging Write %t but got %t", true, logging.Write)
	}
	if !logging.Delete {
		t.Fatalf("Expected the Logging Delete %t but got %t", true, logging.Delete)
	}
	if !logging.RetentionPolicy.Enabled {
		t.Fatalf("Expected the Logging RetentionPolicy.Enabled %t but got %t", true, logging.RetentionPolicy.Enabled)
	}
	if logging.RetentionPolicy.Days != 7 {
		t.Fatalf("Expected the Logging RetentionPolicy.Enabled %d but got %d", 7, logging.RetentionPolicy.Days)
	}
}
