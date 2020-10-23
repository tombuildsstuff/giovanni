package shares

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
)

func TestGetResourceID(t *testing.T) {
	testData := []struct {
		Environment azure.Environment
		Expected    string
	}{
		{
			Environment: azure.ChinaCloud,
			Expected:    "https://account1.file.core.chinacloudapi.cn/share1",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "https://account1.file.core.cloudapi.de/share1",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "https://account1.file.core.windows.net/share1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "https://account1.file.core.usgovcloudapi.net/share1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual := c.GetResourceID("account1", "share1")
		if actual != v.Expected {
			t.Fatalf("Expected the Resource ID to be %q but got %q", v.Expected, actual)
		}
	}
}

func TestGetResourceManagerResourceID(t *testing.T) {
	testData := []struct {
		Environment azure.Environment
		Expected    string
	}{
		{
			Environment: azure.ChinaCloud,
			Expected:    "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/share1",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/share1",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/share1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/share1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual := c.GetResourceManagerResourceID("11112222-3333-4444-5555-666677778888", "group1", "account1", "share1")
		if actual != v.Expected {
			t.Fatalf("Expected the Resource Manager Resource ID to be %q but got %q", v.Expected, actual)
		}
	}
}

func TestParseResourceID(t *testing.T) {
	testData := []struct {
		Environment azure.Environment
		Input       string
	}{
		{
			Environment: azure.ChinaCloud,
			Input:       "https://account1.file.core.chinacloudapi.cn/share1",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.file.core.cloudapi.de/share1",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.file.core.windows.net/share1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.file.core.usgovcloudapi.net/share1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		actual, err := ParseResourceID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != "account1" {
			t.Fatalf("Expected the account name to be `account1` but got %q", actual.AccountName)
		}

		if actual.ShareName != "share1" {
			t.Fatalf("Expected the share name to be `share1` but got %q", actual.ShareName)
		}
	}
}
