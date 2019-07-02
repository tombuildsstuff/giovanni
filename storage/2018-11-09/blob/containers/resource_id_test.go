package containers

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
			Expected:    "https://account1.blob.core.chinacloudapi.cn/container1",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "https://account1.blob.core.cloudapi.de/container1",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "https://account1.blob.core.windows.net/container1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "https://account1.blob.core.usgovcloudapi.net/container1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual := c.GetResourceID("account1", "container1")
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
			Input:       "https://account1.blob.core.chinacloudapi.cn/container1",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.blob.core.cloudapi.de/container1",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.blob.core.windows.net/container1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.blob.core.usgovcloudapi.net/container1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual, err := c.ParseResourceID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != "account1" {
			t.Fatalf("Expected the account name to be `account1` but got %q", actual.AccountName)
		}

		if actual.ContainerName != "container1" {
			t.Fatalf("Expected the container name to be `container1` but got %q", actual.ContainerName)
		}
	}
}
