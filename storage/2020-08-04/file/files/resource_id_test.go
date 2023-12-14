package files

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

func TestParseFileIDStandard(t *testing.T) {
	input := "https://example1.file.core.windows.net/share1/some/path/file.txt"
	expected := FileId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		ShareName:     "share1",
		DirectoryPath: "some/path",
		FileName:      "file.txt",
	}
	actual, err := ParseFileID(input, "core.windows.net")
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
	if actual.ShareName != expected.ShareName {
		t.Fatalf("expected ShareName to be %q but got %q", expected.ShareName, actual.ShareName)
	}
	if actual.DirectoryPath != expected.DirectoryPath {
		t.Fatalf("expected DirectoryPath to be %q but got %q", expected.DirectoryPath, actual.DirectoryPath)
	}
	if actual.FileName != expected.FileName {
		t.Fatalf("expected FileName to be %q but got %q", expected.FileName, actual.FileName)
	}
}

func TestParseFileIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.file.storage.azure.net/share1/path/file.txt"
	expected := FileId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		ShareName:     "share1",
		DirectoryPath: "path",
		FileName:      "file.txt",
	}
	actual, err := ParseFileID(input, "storage.azure.net")
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
	if actual.ShareName != expected.ShareName {
		t.Fatalf("expected ShareName to be %q but got %q", expected.ShareName, actual.ShareName)
	}
	if actual.DirectoryPath != expected.DirectoryPath {
		t.Fatalf("expected DirectoryPath to be %q but got %q", expected.DirectoryPath, actual.DirectoryPath)
	}
	if actual.FileName != expected.FileName {
		t.Fatalf("expected FileName to be %q but got %q", expected.FileName, actual.FileName)
	}
}

func TestParseFileIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.file.zone1.edgestorage.azure.net/share1/some/path/file.txt"
	expected := FileId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		ShareName:     "share1",
		DirectoryPath: "some/path",
		FileName:      "file.txt",
	}
	actual, err := ParseFileID(input, "edgestorage.azure.net")
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
	if actual.ShareName != expected.ShareName {
		t.Fatalf("expected ShareName to be %q but got %q", expected.ShareName, actual.ShareName)
	}
	if actual.DirectoryPath != expected.DirectoryPath {
		t.Fatalf("expected DirectoryPath to be %q but got %q", expected.DirectoryPath, actual.DirectoryPath)
	}
	if actual.FileName != expected.FileName {
		t.Fatalf("expected FileName to be %q but got %q", expected.FileName, actual.FileName)
	}
}

func TestFormatContainerIDStandard(t *testing.T) {
	actual := FileId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		ShareName:     "share1",
		DirectoryPath: "some/path",
		FileName:      "file.txt",
	}.ID()
	expected := "https://example1.file.core.windows.net/share1/some/path/file.txt"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatContainerIDInDNSZone(t *testing.T) {
	actual := FileId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		ShareName:     "share1",
		DirectoryPath: "some/path",
		FileName:      "file.txt",
	}.ID()
	expected := "https://example1.zone2.file.storage.azure.net/share1/some/path/file.txt"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatContainerIDInEdgeZone(t *testing.T) {
	actual := FileId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		ShareName:     "share1",
		DirectoryPath: "path",
		FileName:      "file.txt",
	}.ID()
	expected := "https://example1.file.zone2.edgestorage.azure.net/share1/path/file.txt"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
