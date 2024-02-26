package tables

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
)

func TestGetResourceManagerResourceID(t *testing.T) {
	actual := Client{}.GetResourceManagerResourceID("11112222-3333-4444-5555-666677778888", "group1", "account1", "table1")
	expected := "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/tableServices/default/tables/table1"
	if actual != expected {
		t.Fatalf("Expected the Resource Manager Resource ID to be %q but got %q", expected, actual)
	}
}

func TestParseTableIDStandard(t *testing.T) {
	input := "https://example1.table.core.windows.net/Tables('table1')"
	expected := TableId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		TableName: "table1",
	}
	actual, err := ParseTableID(input, "core.windows.net")
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
	if actual.TableName != expected.TableName {
		t.Fatalf("expected TableName to be %q but got %q", expected.TableName, actual.TableName)
	}
}

func TestParseTableIDLegacy(t *testing.T) {
	input := "https://example1.table.core.windows.net/table1"
	expected := TableId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		TableName: "table1",
	}
	actual, err := ParseTableID(input, "core.windows.net")
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
	if actual.TableName != expected.TableName {
		t.Fatalf("expected TableName to be %q but got %q", expected.TableName, actual.TableName)
	}
}

func TestParseTableIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.table.storage.azure.net/Tables('table1')"
	expected := TableId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		TableName: "table1",
	}
	actual, err := ParseTableID(input, "storage.azure.net")
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
	if actual.TableName != expected.TableName {
		t.Fatalf("expected TableName to be %q but got %q", expected.TableName, actual.TableName)
	}
}

func TestParseTableIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.table.zone1.edgestorage.azure.net/Tables('table1')"
	expected := TableId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		TableName: "table1",
	}
	actual, err := ParseTableID(input, "edgestorage.azure.net")
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
	if actual.TableName != expected.TableName {
		t.Fatalf("expected TableName to be %q but got %q", expected.TableName, actual.TableName)
	}
}

func TestFormatTableIDStandard(t *testing.T) {
	actual := TableId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		TableName: "table1",
	}.ID()
	expected := "https://example1.table.core.windows.net/Tables('table1')"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatTableIDInDNSZone(t *testing.T) {
	actual := TableId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		TableName: "table1",
	}.ID()
	expected := "https://example1.zone2.table.storage.azure.net/Tables('table1')"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatTableIDInEdgeZone(t *testing.T) {
	actual := TableId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		TableName: "table1",
	}.ID()
	expected := "https://example1.table.zone2.edgestorage.azure.net/Tables('table1')"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
