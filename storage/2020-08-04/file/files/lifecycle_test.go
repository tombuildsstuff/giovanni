package files

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var _ StorageFile = Client{}

func TestFilesLifeCycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.KindStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLite)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
	sharesClient := shares.NewWithEnvironment(accountName, client.AutoRestEnvironment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	input := shares.CreateInput{
		QuotaInGB: 1,
	}
	_, err = sharesClient.Create(ctx, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, shareName, false)

	filesClient := NewWithEnvironment(accountName, client.AutoRestEnvironment)
	filesClient.Client = client.PrepareWithAuthorizer(filesClient.Client, storageAuth)

	fileName := "bled5.png"
	contentEncoding := "application/vnd+panda"

	t.Logf("[DEBUG] Creating Top Level File..")
	createInput := CreateInput{
		ContentLength:   1024,
		ContentEncoding: &contentEncoding,
	}
	if _, err := filesClient.Create(ctx, shareName, "", fileName, createInput); err != nil {
		t.Fatalf("Error creating Top-Level File: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties for the Top-Level File..")
	file, err := filesClient.GetProperties(ctx, shareName, "", fileName)
	if err != nil {
		t.Fatalf("Error retrieving Top-Level File: %s", err)
	}

	if *file.ContentLength != 1024 {
		t.Fatalf("Expected the Content-Length to be 1024 but got %d", *file.ContentLength)
	}

	if file.ContentEncoding != contentEncoding {
		t.Fatalf("Expected the Content-Encoding to be %q but got %q", contentEncoding, file.ContentEncoding)
	}

	updatedSize := int64(2048)
	updatedEncoding := "application/vnd+pandas2"
	updatedInput := SetPropertiesInput{
		ContentEncoding: &updatedEncoding,
		ContentLength:   updatedSize,
		MetaData: map[string]string{
			"bingo": "bango",
		},
	}
	t.Logf("[DEBUG] Setting Properties for the Top-Level File..")
	if _, err := filesClient.SetProperties(ctx, shareName, "", fileName, updatedInput); err != nil {
		t.Fatalf("Error setting properties: %s", err)
	}

	t.Logf("[DEBUG] Re-retrieving Properties for the Top-Level File..")
	file, err = filesClient.GetProperties(ctx, shareName, "", fileName)
	if err != nil {
		t.Fatalf("Error retrieving Top-Level File: %s", err)
	}

	if *file.ContentLength != 2048 {
		t.Fatalf("Expected the Content-Length to be 1024 but got %d", *file.ContentLength)
	}

	if file.ContentEncoding != updatedEncoding {
		t.Fatalf("Expected the Content-Encoding to be %q but got %q", updatedEncoding, file.ContentEncoding)
	}

	if len(file.MetaData) != 1 {
		t.Fatalf("Expected 1 item but got %d", len(file.MetaData))
	}
	if file.MetaData["bingo"] != "bango" {
		t.Fatalf("Expected `bingo` to be `bango` but got %q", file.MetaData["bingo"])
	}

	t.Logf("[DEBUG] Setting MetaData..")
	metaData := map[string]string{
		"hello": "there",
	}
	if _, err := filesClient.SetMetaData(ctx, shareName, "", fileName, metaData); err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	t.Logf("[DEBUG] Retrieving MetaData..")
	retrievedMetaData, err := filesClient.GetMetaData(ctx, shareName, "", fileName)
	if err != nil {
		t.Fatalf("Error retrieving MetaData: %s", err)
	}
	if len(retrievedMetaData.MetaData) != 1 {
		t.Fatalf("Expected 1 item but got %d", len(retrievedMetaData.MetaData))
	}
	if retrievedMetaData.MetaData["hello"] != "there" {
		t.Fatalf("Expected `hello` to be `there` but got %q", retrievedMetaData.MetaData["hello"])
	}

	t.Logf("[DEBUG] Re-Setting MetaData..")
	metaData = map[string]string{
		"hello":  "there",
		"second": "thing",
	}
	if _, err := filesClient.SetMetaData(ctx, shareName, "", fileName, metaData); err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	t.Logf("[DEBUG] Re-Retrieving MetaData..")
	retrievedMetaData, err = filesClient.GetMetaData(ctx, shareName, "", fileName)
	if err != nil {
		t.Fatalf("Error retrieving MetaData: %s", err)
	}
	if len(retrievedMetaData.MetaData) != 2 {
		t.Fatalf("Expected 2 items but got %d", len(retrievedMetaData.MetaData))
	}
	if retrievedMetaData.MetaData["hello"] != "there" {
		t.Fatalf("Expected `hello` to be `there` but got %q", retrievedMetaData.MetaData["hello"])
	}
	if retrievedMetaData.MetaData["second"] != "thing" {
		t.Fatalf("Expected `second` to be `thing` but got %q", retrievedMetaData.MetaData["second"])
	}

	t.Logf("[DEBUG] Deleting Top Level File..")
	if _, err := filesClient.Delete(ctx, shareName, "", fileName); err != nil {
		t.Fatalf("Error deleting Top-Level File: %s", err)
	}
}
