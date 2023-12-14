package filesystems

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

func TestGetResourceManagerResourceID(t *testing.T) {
	actual := Client{}.GetResourceManagerResourceID("11112222-3333-4444-5555-666677778888", "group1", "account1", "container1")
	expected := "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1"
	if actual != expected {
		t.Fatalf("Expected the Resource Manager Resource ID to be %q but got %q", expected, actual)
	}
}

func TestParseFileSystemIDStandard(t *testing.T) {
	input := "https://example1.dfs.core.windows.net/fileSystem1"
	expected := FileSystemId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		FileSystemName: "fileSystem1",
	}
	actual, err := ParseFileSystemID(input, "core.windows.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountId.AccountName != expected.AccountId.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountId.AccountName, actual.AccountId.AccountName)
	}
	if actual.AccountId.SubDomainType != expected.AccountId.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.AccountId.SubDomainType, actual.AccountId.SubDomainType)
	}
	if actual.AccountId.DomainSuffix != expected.AccountId.DomainSuffix {
		t.Fatalf("expected DomainSuffix to be %q but got %q", expected.AccountId.DomainSuffix, actual.AccountId.DomainSuffix)
	}
	if actual.FileSystemName != expected.FileSystemName {
		t.Fatalf("expected FileSystemName to be %q but got %q", expected.FileSystemName, actual.FileSystemName)
	}
}

func TestParseFileSystemIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.dfs.storage.azure.net/fileSystem1"
	expected := FileSystemId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		FileSystemName: "fileSystem1",
	}
	actual, err := ParseFileSystemID(input, "storage.azure.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountId.AccountName != expected.AccountId.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountId.AccountName, actual.AccountId.AccountName)
	}
	if actual.AccountId.SubDomainType != expected.AccountId.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.AccountId.SubDomainType, actual.AccountId.SubDomainType)
	}
	if actual.AccountId.DomainSuffix != expected.AccountId.DomainSuffix {
		t.Fatalf("expected DomainSuffix to be %q but got %q", expected.AccountId.DomainSuffix, actual.AccountId.DomainSuffix)
	}
	if pointer.From(actual.AccountId.ZoneName) != pointer.From(expected.AccountId.ZoneName) {
		t.Fatalf("expected ZoneName to be %q but got %q", pointer.From(expected.AccountId.ZoneName), pointer.From(actual.AccountId.ZoneName))
	}
	if actual.FileSystemName != expected.FileSystemName {
		t.Fatalf("expected FileSystemName to be %q but got %q", expected.FileSystemName, actual.FileSystemName)
	}
}

func TestParseFileSystemIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.dfs.zone1.edgestorage.azure.net/fileSystem1"
	expected := FileSystemId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		FileSystemName: "fileSystem1",
	}
	actual, err := ParseFileSystemID(input, "edgestorage.azure.net")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual.AccountId.AccountName != expected.AccountId.AccountName {
		t.Fatalf("expected AccountName to be %q but got %q", expected.AccountId.AccountName, actual.AccountId.AccountName)
	}
	if actual.AccountId.SubDomainType != expected.AccountId.SubDomainType {
		t.Fatalf("expected SubDomainType to be %q but got %q", expected.AccountId.SubDomainType, actual.AccountId.SubDomainType)
	}
	if actual.AccountId.DomainSuffix != expected.AccountId.DomainSuffix {
		t.Fatalf("expected DomainSuffix to be %q but got %q", expected.AccountId.DomainSuffix, actual.AccountId.DomainSuffix)
	}
	if pointer.From(actual.AccountId.ZoneName) != pointer.From(expected.AccountId.ZoneName) {
		t.Fatalf("expected ZoneName to be %q but got %q", pointer.From(expected.AccountId.ZoneName), pointer.From(actual.AccountId.ZoneName))
	}
	if !actual.AccountId.IsEdgeZone {
		t.Fatalf("expected the Account to be in an Edge Zone but it wasn't")
	}
	if actual.FileSystemName != expected.FileSystemName {
		t.Fatalf("expected FileSystemName to be %q but got %q", expected.FileSystemName, actual.FileSystemName)
	}
}

func TestFormatFileSystemIDStandard(t *testing.T) {
	actual := FileSystemId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		FileSystemName: "fileSystem1",
	}.ID()
	expected := "https://example1.dfs.core.windows.net/fileSystem1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatFileSystemIDInDNSZone(t *testing.T) {
	actual := FileSystemId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		FileSystemName: "fileSystem1",
	}.ID()
	expected := "https://example1.zone2.dfs.storage.azure.net/fileSystem1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatFileSystemIDInEdgeZone(t *testing.T) {
	actual := FileSystemId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		FileSystemName: "fileSystem1",
	}.ID()
	expected := "https://example1.dfs.zone2.edgestorage.azure.net/fileSystem1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
