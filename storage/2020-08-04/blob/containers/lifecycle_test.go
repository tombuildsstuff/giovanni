package containers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var _ StorageContainer = Client{}

func TestContainerLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.KindBlobStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth, err := autorest.NewSharedKeyAuthorizer(accountName, testData.StorageAccountKey, autorest.SharedKeyLite)
	if err != nil {
		t.Fatalf("building SharedKeyAuthorizer: %+v", err)
	}
	containersClient := NewWithEnvironment(client.AutoRestEnvironment)
	containersClient.Client = client.PrepareWithAuthorizer(containersClient.Client, storageAuth)

	// first let's test an empty container
	input := CreateInput{}
	_, err = containersClient.Create(ctx, accountName, containerName, input)
	if err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	container, err := containersClient.GetProperties(ctx, accountName, containerName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error retrieving: %s", err))
	}

	if container.AccessLevel != Private {
		t.Fatalf("Expected Access Level to be Private but got %q", container.AccessLevel)
	}
	if len(container.MetaData) != 0 {
		t.Fatalf("Expected MetaData to be empty but got: %s", container.MetaData)
	}
	if container.LeaseStatus != Unlocked {
		t.Fatalf("Expected Container Lease to be Unlocked but was: %s", container.LeaseStatus)
	}

	// then update the metadata
	metaData := map[string]string{
		"dont": "kill-my-vibe",
	}
	_, err = containersClient.SetMetaData(ctx, accountName, containerName, metaData)
	if err != nil {
		t.Fatal(fmt.Errorf("Error updating metadata: %s", err))
	}

	// give azure time to replicate
	time.Sleep(2 * time.Second)

	// then assert that
	container, err = containersClient.GetProperties(ctx, accountName, containerName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error re-retrieving: %s", err))
	}
	if len(container.MetaData) != 1 {
		t.Fatalf("Expected 1 item in the metadata but got: %s", container.MetaData)
	}
	if container.MetaData["dont"] != "kill-my-vibe" {
		t.Fatalf("Expected `kill-my-vibe` but got %q", container.MetaData["dont"])
	}
	if container.AccessLevel != Private {
		t.Fatalf("Expected Access Level to be Private but got %q", container.AccessLevel)
	}
	if container.LeaseStatus != Unlocked {
		t.Fatalf("Expected Container Lease to be Unlocked but was: %s", container.LeaseStatus)
	}

	// then update the ACL
	_, err = containersClient.SetAccessControl(ctx, accountName, containerName, Blob)
	if err != nil {
		t.Fatal(fmt.Errorf("Error updating ACL's: %s", err))
	}

	// give azure some time to replicate
	time.Sleep(2 * time.Second)

	// then assert that
	container, err = containersClient.GetProperties(ctx, accountName, containerName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error re-retrieving: %s", err))
	}
	if container.AccessLevel != Blob {
		t.Fatalf("Expected Access Level to be Blob but got %q", container.AccessLevel)
	}
	if len(container.MetaData) != 1 {
		t.Fatalf("Expected 1 item in the metadata but got: %s", container.MetaData)
	}
	if container.LeaseStatus != Unlocked {
		t.Fatalf("Expected Container Lease to be Unlocked but was: %s", container.LeaseStatus)
	}

	// acquire a lease for 30s
	acquireLeaseInput := AcquireLeaseInput{
		LeaseDuration: 30,
	}
	acquireLeaseResp, err := containersClient.AcquireLease(ctx, accountName, containerName, acquireLeaseInput)
	if err != nil {
		t.Fatalf("Error acquiring lease: %s", err)
	}
	t.Logf("[DEBUG] Lease ID: %s", acquireLeaseResp.LeaseID)

	// we should then be able to update the ID
	t.Logf("[DEBUG] Changing lease..")
	updateLeaseInput := ChangeLeaseInput{
		ExistingLeaseID: acquireLeaseResp.LeaseID,
		ProposedLeaseID: "aaaabbbb-aaaa-bbbb-cccc-aaaabbbbcccc",
	}
	updateLeaseResp, err := containersClient.ChangeLease(ctx, accountName, containerName, updateLeaseInput)
	if err != nil {
		t.Fatalf("Error changing lease: %s", err)
	}

	// then renew it
	_, err = containersClient.RenewLease(ctx, accountName, containerName, updateLeaseResp.LeaseID)
	if err != nil {
		t.Fatalf("Error renewing lease: %s", err)
	}

	// and then give it a timeout
	breakPeriod := 20
	breakLeaseInput := BreakLeaseInput{
		LeaseID:     updateLeaseResp.LeaseID,
		BreakPeriod: &breakPeriod,
	}
	breakLeaseResp, err := containersClient.BreakLease(ctx, accountName, containerName, breakLeaseInput)
	if err != nil {
		t.Fatalf("Error breaking lease: %s", err)
	}
	if breakLeaseResp.LeaseTime == 0 {
		t.Fatalf("Lease broke immediately when should have waited: %d", breakLeaseResp.LeaseTime)
	}

	// and finally ditch it
	_, err = containersClient.ReleaseLease(ctx, accountName, containerName, updateLeaseResp.LeaseID)
	if err != nil {
		t.Fatalf("Error releasing lease: %s", err)
	}

	t.Logf("[DEBUG] Listing blobs in the container..")
	listInput := ListBlobsInput{}
	listResult, err := containersClient.ListBlobs(ctx, accountName, containerName, listInput)
	if err != nil {
		t.Fatalf("Error listing blobs: %s", err)
	}

	if len(listResult.Blobs.Blobs) != 0 {
		t.Fatalf("Expected there to be no blobs in the container but got %d", len(listResult.Blobs.Blobs))
	}

	t.Logf("[DEBUG] Deleting..")
	_, err = containersClient.Delete(ctx, accountName, containerName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error deleting: %s", err))
	}
}
