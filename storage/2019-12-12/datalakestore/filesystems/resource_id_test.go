package filesystems

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
			Expected:    "https://account1.dfs.core.chinacloudapi.cn/directory1",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "https://account1.dfs.core.cloudapi.de/directory1",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "https://account1.dfs.core.windows.net/directory1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "https://account1.dfs.core.usgovcloudapi.net/directory1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual := c.GetResourceID("account1", "directory1")
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
			Expected:    "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "/subscriptions/11112222-3333-4444-5555-666677778888/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual := c.GetResourceManagerResourceID("11112222-3333-4444-5555-666677778888", "group1", "account1", "container1")
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
			Input:       "https://account1.dfs.core.chinacloudapi.cn/directory1",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.dfs.core.cloudapi.de/directory1",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.dfs.core.windows.net/directory1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.dfs.core.usgovcloudapi.net/directory1",
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

		if actual.DirectoryName != "directory1" {
			t.Fatalf("Expected the directory name to be `directory1` but got %q", actual.DirectoryName)
		}
	}
}
