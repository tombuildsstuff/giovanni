package paths

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
)

func TestGetResourceID(t *testing.T) {
	testData := []struct {
		AccountName    string
		FileSystemName string
		Path           string
		Expected       string
	}{
		{
			AccountName:    "account1",
			FileSystemName: "fs1",
			Path:           "test",
			Expected:       "https://account1.dfs.core.windows.net/fs1/test",
		},
		{
			AccountName:    "account1",
			FileSystemName: "fs1",
			Path:           "test/test2",
			Expected:       "https://account1.dfs.core.windows.net/fs1/test/test2",
		},
		{
			AccountName:    "account1",
			FileSystemName: "fs1",
			Path:           "",
			Expected:       "https://account1.dfs.core.windows.net/fs1/",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Path %q", v.Path)
		c := NewWithEnvironment(azure.PublicCloud)
		actual := c.GetResourceID(v.AccountName, v.FileSystemName, v.Path)
		if actual != v.Expected {
			t.Fatalf("Expected the Resource ID to be %q but got %q", v.Expected, actual)
		}
	}
}

func TestParseResourceID(t *testing.T) {
	testData := []struct {
		Input          string
		AccountName    string
		FileSystemName string
		Path           string
	}{
		{
			Input:          "https://account1.dfs.core.windows.net/fs1/test",
			AccountName:    "account1",
			FileSystemName: "fs1",
			Path:           "test",
		},
		{
			Input:          "https://account1.dfs.core.windows.net/fs1/test/test2",
			AccountName:    "account1",
			FileSystemName: "fs1",
			Path:           "test/test2",
		},
		{
			Input:          "https://account1.dfs.core.windows.net/fs1/",
			AccountName:    "account1",
			FileSystemName: "fs1",
			Path:           "",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Path %q", v.Path)
		actual, err := ParseResourceID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != v.AccountName {
			t.Fatalf("Expected the account name to be `%q` but got %q", v.AccountName, actual.AccountName)
		}

		if actual.FileSystemName != v.FileSystemName {
			t.Fatalf("Expected the file system name to be `%q` but got %q", v.FileSystemName, actual.FileSystemName)
		}

		if actual.Path != v.Path {
			t.Fatalf("Expected the path to be `%q` but got %q", v.Path, actual.Path)
		}
	}
}
