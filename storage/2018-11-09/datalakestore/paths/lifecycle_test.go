package paths

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestLifecycle(t *testing.T) {

	const defaultACLString = "user::rwx,group::r-x,other::---"

	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	fileSystemName := fmt.Sprintf("acctestfs-%s", testhelpers.RandomString())
	path := "test"

	if _, err = client.BuildTestResourcesWithHns(ctx, resourceGroup, accountName, storage.BlobStorage); err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)
	fileSystemsClient := filesystems.NewWithEnvironment(client.Environment)
	fileSystemsClient.Client = client.PrepareWithStorageResourceManagerAuth(fileSystemsClient.Client)
	pathsClient := NewWithEnvironment(client.Environment)
	pathsClient.Client = client.PrepareWithStorageResourceManagerAuth(fileSystemsClient.Client)

	t.Logf("[DEBUG] Creating an empty File System..")
	fileSystemInput := filesystems.CreateInput{}
	if _, err = fileSystemsClient.Create(ctx, accountName, fileSystemName, fileSystemInput); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Creating folder 'test' ..")
	input := CreateInput{
		Resource: PathResourceDirectory,
	}
	if _, err = pathsClient.Create(ctx, accountName, fileSystemName, path, input); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Getting properties for folder 'test' ..")
	props, err := pathsClient.GetProperties(ctx, accountName, fileSystemName, path, GetPropertiesActionGetAccessControl)
	if err != nil {
		t.Fatal(fmt.Errorf("Error getting properties: %s", err))
	}
	t.Logf("[DEBUG] Props.Owner: %q", props.Owner)
	t.Logf("[DEBUG] Props.Group: %q", props.Group)
	t.Logf("[DEBUG] Props.ACL: %q", props.ACL)
	t.Logf("[DEBUG] Props.ETag: %q", props.ETag)
	t.Logf("[DEBUG] Props.LastModified: %q", props.LastModified)
	if props.ACL != defaultACLString {
		t.Fatal(fmt.Errorf("Expected Default ACL %q, got %q", defaultACLString, props.ACL))
	}

	newACL := "user::rwx,group::r-x,other::r-x,default:user::rwx,default:group::r-x,default:other::---"
	accessControlInput := SetAccessControlInput{
		ACL: &newACL,
	}
	t.Logf("[DEBUG] Setting Access Control for folder 'test' ..")
	if _, err = pathsClient.SetAccessControl(ctx, accountName, fileSystemName, path, accessControlInput); err != nil {
		t.Fatal(fmt.Errorf("Error setting Access Control %s", err))
	}

	t.Logf("[DEBUG] Getting properties for folder 'test' (2) ..")
	props, err = pathsClient.GetProperties(ctx, accountName, fileSystemName, path, GetPropertiesActionGetAccessControl)
	if err != nil {
		t.Fatal(fmt.Errorf("Error getting properties (2): %s", err))
	}
	if props.ACL != newACL {
		t.Fatal(fmt.Errorf("Expected new ACL %q, got %q", newACL, props.ACL))
	}

	t.Logf("[DEBUG] Deleting path 'test' ..")
	if _, err = pathsClient.Delete(ctx, accountName, fileSystemName, path); err != nil {
		t.Fatal(fmt.Errorf("Error deleting path: %s", err))
	}

	t.Logf("[DEBUG] Deleting File System..")
	if _, err := fileSystemsClient.Delete(ctx, accountName, fileSystemName); err != nil {
		t.Fatalf("Error deleting filesystem: %s", err)
	}
}
