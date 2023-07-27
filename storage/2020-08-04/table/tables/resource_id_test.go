package tables

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
			Expected:    "https://account1.table.core.chinacloudapi.cn/Tables('table1')",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "https://account1.table.core.cloudapi.de/Tables('table1')",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "https://account1.table.core.windows.net/Tables('table1')",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "https://account1.table.core.usgovcloudapi.net/Tables('table1')",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment("account1", v.Environment)
		actual := c.GetResourceID("table1")
		if actual != v.Expected {
			t.Fatalf("Expected the Resource ID to be %q but got %q", v.Expected, actual)
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
			Input:       "https://account1.table.core.chinacloudapi.cn/Tables('table1')",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.table.core.cloudapi.de/Tables('table1')",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.table.core.windows.net/Tables('table1')",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.table.core.usgovcloudapi.net/Tables('table1')",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		actual, err := ParseResourceID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != "account1" {
			t.Fatalf("Expected Account Name to be `account1` but got %q", actual.AccountName)
		}
		if actual.TableName != "table1" {
			t.Fatalf("Expected Table Name to be `table1` but got %q", actual.TableName)
		}
	}
}
