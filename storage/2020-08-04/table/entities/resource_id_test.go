package entities

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

func TestParseEntityIDStandard(t *testing.T) {
	input := "https://example1.table.core.windows.net/table1(PartitionKey='partition1',RowKey='row1')"
	expected := EntityId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		TableName:    "table1",
		PartitionKey: "partition1",
		RowKey:       "row1",
	}
	actual, err := ParseEntityID(input, "core.windows.net")
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
	if actual.PartitionKey != expected.PartitionKey {
		t.Fatalf("expected PartitionKey to be %q but got %q", expected.PartitionKey, actual.PartitionKey)
	}
	if actual.RowKey != expected.RowKey {
		t.Fatalf("expected RowKey to be %q but got %q", expected.RowKey, actual.RowKey)
	}
}

func TestParseEntityIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.table.storage.azure.net/table1(PartitionKey='partition1',RowKey='row1')"
	expected := EntityId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		TableName:    "table1",
		PartitionKey: "partition1",
		RowKey:       "row1",
	}
	actual, err := ParseEntityID(input, "storage.azure.net")
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
	if actual.PartitionKey != expected.PartitionKey {
		t.Fatalf("expected PartitionKey to be %q but got %q", expected.PartitionKey, actual.PartitionKey)
	}
	if actual.RowKey != expected.RowKey {
		t.Fatalf("expected RowKey to be %q but got %q", expected.RowKey, actual.RowKey)
	}
}

func TestParseEntityIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.table.zone1.edgestorage.azure.net/table1(PartitionKey='partition1',RowKey='row1')"
	expected := EntityId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		TableName:    "table1",
		PartitionKey: "partition1",
		RowKey:       "row1",
	}
	actual, err := ParseEntityID(input, "edgestorage.azure.net")
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
	if actual.PartitionKey != expected.PartitionKey {
		t.Fatalf("expected PartitionKey to be %q but got %q", expected.PartitionKey, actual.PartitionKey)
	}
	if actual.RowKey != expected.RowKey {
		t.Fatalf("expected RowKey to be %q but got %q", expected.RowKey, actual.RowKey)
	}
}

func TestFormatEntityIDStandard(t *testing.T) {
	actual := EntityId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		TableName:    "table1",
		PartitionKey: "partition1",
		RowKey:       "row1",
	}.ID()
	expected := "https://example1.table.core.windows.net/table1(PartitionKey='partition1',RowKey='row1')"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatEntityIDInDNSZone(t *testing.T) {
	actual := EntityId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		TableName:    "table1",
		PartitionKey: "partition1",
		RowKey:       "row1",
	}.ID()
	expected := "https://example1.zone2.table.storage.azure.net/table1(PartitionKey='partition1',RowKey='row1')"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatEntityIDInEdgeZone(t *testing.T) {
	actual := EntityId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.TableSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		TableName:    "table1",
		PartitionKey: "partition1",
		RowKey:       "row1",
	}.ID()
	expected := "https://example1.table.zone2.edgestorage.azure.net/table1(PartitionKey='partition1',RowKey='row1')"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
