package paths

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

func TestParsePathIDStandard(t *testing.T) {
	input := "https://example1.dfs.core.windows.net/fileSystem1/some/path"
	expected := PathId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		FileSystemName: "fileSystem1",
		Path:           "some/path",
	}
	actual, err := ParsePathID(input, "core.windows.net")
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
	if actual.Path != expected.Path {
		t.Fatalf("expected Path to be %q but got %q", expected.Path, actual.Path)
	}
}

func TestParsePathIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.dfs.storage.azure.net/fileSystem1/some/path"
	expected := PathId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		FileSystemName: "fileSystem1",
		Path:           "some/path",
	}
	actual, err := ParsePathID(input, "storage.azure.net")
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
	if actual.Path != expected.Path {
		t.Fatalf("expected Path to be %q but got %q", expected.Path, actual.Path)
	}
}

func TestParsePathIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.dfs.zone1.edgestorage.azure.net/fileSystem1/some/path"
	expected := PathId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		FileSystemName: "fileSystem1",
		Path:           "some/path",
	}
	actual, err := ParsePathID(input, "edgestorage.azure.net")
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
	if actual.Path != expected.Path {
		t.Fatalf("expected Path to be %q but got %q", expected.Path, actual.Path)
	}
}

func TestFormatPathIDStandard(t *testing.T) {
	actual := PathId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		FileSystemName: "fileSystem1",
		Path:           "some/path",
	}.ID()
	expected := "https://example1.dfs.core.windows.net/fileSystem1/some/path"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatPathIDInDNSZone(t *testing.T) {
	actual := PathId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		FileSystemName: "fileSystem1",
		Path:           "some/path",
	}.ID()
	expected := "https://example1.zone2.dfs.storage.azure.net/fileSystem1/some/path"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatPathIDInEdgeZone(t *testing.T) {
	actual := PathId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.DataLakeStoreSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		FileSystemName: "fileSystem1",
		Path:           "some",
	}.ID()
	expected := "https://example1.dfs.zone2.edgestorage.azure.net/fileSystem1/some"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
