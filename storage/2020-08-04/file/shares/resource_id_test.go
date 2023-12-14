package shares

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

func TestGetResourceManagerResourceID(t *testing.T) {
	actual := Client{}.GetResourceManagerResourceID("11112222-3333-4444-5555-666677778888", "group1", "account1", "share1")
	expected := "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/share1"
	if actual != expected {
		t.Fatalf("Expected the Resource Manager Resource ID to be %q but got %q", expected, actual)
	}
}

func TestParseShareIDStandard(t *testing.T) {
	input := "https://example1.file.core.windows.net/share1"
	expected := ShareId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		ShareName: "share1",
	}
	actual, err := ParseShareID(input, "core.windows.net")
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
}

func TestParseShareIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.file.storage.azure.net/share1"
	expected := ShareId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		ShareName: "share1",
	}
	actual, err := ParseShareID(input, "storage.azure.net")
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
}

func TestParseShareIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.file.zone1.edgestorage.azure.net/share1"
	expected := ShareId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		ShareName: "share1",
	}
	actual, err := ParseShareID(input, "edgestorage.azure.net")
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
}

func TestFormatContainerIDStandard(t *testing.T) {
	actual := ShareId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		ShareName: "share1",
	}.ID()
	expected := "https://example1.file.core.windows.net/share1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatContainerIDInDNSZone(t *testing.T) {
	actual := ShareId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		ShareName: "share1",
	}.ID()
	expected := "https://example1.zone2.file.storage.azure.net/share1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatContainerIDInEdgeZone(t *testing.T) {
	actual := ShareId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.FileSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		ShareName: "share1",
	}.ID()
	expected := "https://example1.file.zone2.edgestorage.azure.net/share1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
