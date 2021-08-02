package blobs

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
			Expected:    "https://account1.blob.core.chinacloudapi.cn/container1/blob1.vhd",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "https://account1.blob.core.cloudapi.de/container1/blob1.vhd",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "https://account1.blob.core.windows.net/container1/blob1.vhd",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "https://account1.blob.core.usgovcloudapi.net/container1/blob1.vhd",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual := c.GetResourceID("account1", "container1", "blob1.vhd")
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
			Input:       "https://account1.blob.core.chinacloudapi.cn/container1/blob1.vhd",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.blob.core.cloudapi.de/container1/blob1.vhd",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.blob.core.windows.net/container1/blob1.vhd",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.blob.core.usgovcloudapi.net/container1/blob1.vhd",
		},
	}
	t.Logf("[DEBUG] Top Level Files")
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		actual, err := ParseResourceID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != "account1" {
			t.Fatalf("Expected Account Name to be `account1` but got %q", actual.AccountName)
		}
		if actual.ContainerName != "container1" {
			t.Fatalf("Expected Container Name to be `container1` but got %q", actual.ContainerName)
		}
		if actual.BlobName != "blob1.vhd" {
			t.Fatalf("Expected Blob Name to be `blob1.vhd` but got %q", actual.BlobName)
		}
	}

	testData = []struct {
		Environment azure.Environment
		Input       string
	}{
		{
			Environment: azure.ChinaCloud,
			Input:       "https://account1.blob.core.chinacloudapi.cn/container1/example/blob1.vhd",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.blob.core.cloudapi.de/container1/example/blob1.vhd",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.blob.core.windows.net/container1/example/blob1.vhd",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.blob.core.usgovcloudapi.net/container1/example/blob1.vhd",
		},
	}
	t.Logf("[DEBUG] Nested Files")
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		actual, err := ParseResourceID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != "account1" {
			t.Fatalf("Expected Account Name to be `account1` but got %q", actual.AccountName)
		}
		if actual.ContainerName != "container1" {
			t.Fatalf("Expected Container Name to be `container1` but got %q", actual.ContainerName)
		}
		if actual.BlobName != "example/blob1.vhd" {
			t.Fatalf("Expected Blob Name to be `example/blob1.vhd` but got %q", actual.BlobName)
		}
	}
}
