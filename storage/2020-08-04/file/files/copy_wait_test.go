package files

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestFilesCopyAndWaitFromURL(t *testing.T) {
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
	sharesClient := shares.NewWithEnvironment(client.AutoRestEnvironment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	input := shares.CreateInput{
		QuotaInGB: 10,
	}
	_, err = sharesClient.Create(ctx, accountName, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, accountName, shareName, false)

	filesClient := NewWithEnvironment(client.AutoRestEnvironment)
	filesClient.Client = client.PrepareWithAuthorizer(filesClient.Client, storageAuth)

	copiedFileName := "ubuntu.iso"
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
	}

	t.Logf("[DEBUG] Copy And Waiting..")
	if _, err := filesClient.CopyAndWait(ctx, accountName, shareName, "", copiedFileName, copyInput, DefaultCopyPollDuration); err != nil {
		t.Fatalf("Error copy & waiting: %s", err)
	}

	t.Logf("[DEBUG] Asserting that the file's ready..")

	props, err := filesClient.GetProperties(ctx, accountName, shareName, "", copiedFileName)
	if err != nil {
		t.Fatalf("Error retrieving file: %s", err)
	}

	if !strings.EqualFold(props.CopyStatus, "success") {
		t.Fatalf("Expected the Copy Status to be `Success` but got %q", props.CopyStatus)
	}
}

func TestFilesCopyAndWaitFromBlob(t *testing.T) {
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
	sharesClient := shares.NewWithEnvironment(client.AutoRestEnvironment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	input := shares.CreateInput{
		QuotaInGB: 10,
	}
	_, err = sharesClient.Create(ctx, accountName, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, accountName, shareName, false)

	filesClient := NewWithEnvironment(client.AutoRestEnvironment)
	filesClient.Client = client.PrepareWithAuthorizer(filesClient.Client, storageAuth)

	originalFileName := "ubuntu.iso"
	copiedFileName := "ubuntu-copied.iso"
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
	}
	t.Logf("[DEBUG] Copy And Waiting the original file..")
	if _, err := filesClient.CopyAndWait(ctx, accountName, shareName, "", originalFileName, copyInput, DefaultCopyPollDuration); err != nil {
		t.Fatalf("Error copy & waiting: %s", err)
	}

	t.Logf("[DEBUG] Now copying that blob..")
	duplicateInput := CopyInput{
		CopySource: fmt.Sprintf("%s/%s/%s", endpoints.GetOrBuildFileEndpoint(client.endpoint, filesClient.BaseURI, accountName), shareName, originalFileName),
	}
	if _, err := filesClient.CopyAndWait(ctx, accountName, shareName, "", copiedFileName, duplicateInput, DefaultCopyPollDuration); err != nil {
		t.Fatalf("Error copying duplicate: %s", err)
	}

	t.Logf("[DEBUG] Asserting that the file's ready..")
	props, err := filesClient.GetProperties(ctx, accountName, shareName, "", copiedFileName)
	if err != nil {
		t.Fatalf("Error retrieving file: %s", err)
	}

	if !strings.EqualFold(props.CopyStatus, "success") {
		t.Fatalf("Expected the Copy Status to be `Success` but got %q", props.CopyStatus)
	}
}
