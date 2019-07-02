package messages

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
			Expected:    "https://account1.queue.core.chinacloudapi.cn/queue1/messages/message1",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "https://account1.queue.core.cloudapi.de/queue1/messages/message1",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "https://account1.queue.core.windows.net/queue1/messages/message1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "https://account1.queue.core.usgovcloudapi.net/queue1/messages/message1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual := c.GetResourceID("account1", "queue1", "message1")
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
			Input:       "https://account1.queue.core.chinacloudapi.cn/queue1/messages/message1",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.queue.core.cloudapi.de/queue1/messages/message1",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.queue.core.windows.net/queue1/messages/message1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.queue.core.usgovcloudapi.net/queue1/messages/message1",
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
			t.Fatalf("Expected Account Name to be `account1` but got %q", actual.AccountName)
		}
		if actual.QueueName != "queue1" {
			t.Fatalf("Expected Queue Name to be `queue1` but got %q", actual.QueueName)
		}
		if actual.MessageID != "message1" {
			t.Fatalf("Expected Message ID to be `message1` but got %q", actual.MessageID)
		}
	}
}
