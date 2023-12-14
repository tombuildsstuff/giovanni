package queues

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
)

func TestGetResourceManagerResourceID(t *testing.T) {
	actual := Client{}.GetResourceManagerResourceID("11112222-3333-4444-5555-666677778888", "group1", "account1", "queue1")
	expected := "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/queueServices/default/queues/queue1"
	if actual != expected {
		t.Fatalf("Expected the Resource Manager Resource ID to be %q but got %q", expected, actual)
	}
}

func TestParseQueueIDStandard(t *testing.T) {
	input := "https://example1.queue.core.windows.net/queue1"
	expected := QueueId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.QueueSubDomainType,
			DomainSuffix:  "core.windows.net",
		},
		QueueName: "queue1",
	}
	actual, err := ParseQueueID(input, "core.windows.net")
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
	if actual.QueueName != expected.QueueName {
		t.Fatalf("expected QueueName to be %q but got %q", expected.QueueName, actual.QueueName)
	}
}

func TestParseQueueIDInADNSZone(t *testing.T) {
	input := "https://example1.zone1.queue.storage.azure.net/queue1"
	expected := QueueId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.QueueSubDomainType,
			DomainSuffix:  "storage.azure.net",
			ZoneName:      pointer.To("zone1"),
		},
		QueueName: "queue1",
	}
	actual, err := ParseQueueID(input, "storage.azure.net")
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
	if actual.QueueName != expected.QueueName {
		t.Fatalf("expected QueueName to be %q but got %q", expected.QueueName, actual.QueueName)
	}
}

func TestParseQueueIDInAnEdgeZone(t *testing.T) {
	input := "https://example1.queue.zone1.edgestorage.azure.net/queue1"
	expected := QueueId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.QueueSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			ZoneName:      pointer.To("zone1"),
			IsEdgeZone:    true,
		},
		QueueName: "queue1",
	}
	actual, err := ParseQueueID(input, "edgestorage.azure.net")
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
	if actual.QueueName != expected.QueueName {
		t.Fatalf("expected QueueName to be %q but got %q", expected.QueueName, actual.QueueName)
	}
}

func TestFormatQueueIDStandard(t *testing.T) {
	actual := QueueId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			SubDomainType: accounts.QueueSubDomainType,
			DomainSuffix:  "core.windows.net",
			IsEdgeZone:    false,
		},
		QueueName: "queue1",
	}.ID()
	expected := "https://example1.queue.core.windows.net/queue1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatQueueIDInDNSZone(t *testing.T) {
	actual := QueueId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.QueueSubDomainType,
			DomainSuffix:  "storage.azure.net",
			IsEdgeZone:    false,
		},
		QueueName: "queue1",
	}.ID()
	expected := "https://example1.zone2.queue.storage.azure.net/queue1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}

func TestFormatQueueIDInEdgeZone(t *testing.T) {
	actual := QueueId{
		AccountId: accounts.AccountId{
			AccountName:   "example1",
			ZoneName:      pointer.To("zone2"),
			SubDomainType: accounts.QueueSubDomainType,
			DomainSuffix:  "edgestorage.azure.net",
			IsEdgeZone:    true,
		},
		QueueName: "queue1",
	}.ID()
	expected := "https://example1.queue.zone2.edgestorage.azure.net/queue1"
	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
