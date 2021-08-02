package messages

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-03-28/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestLifeCycle(t *testing.T) {
	client, err := testhelpers.Build(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	queueName := fmt.Sprintf("queue-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.KindStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	queuesClient := queues.NewWithEnvironment(client.Environment)
	queuesClient.Client = client.PrepareWithStorageResourceManagerAuth(queuesClient.Client)

	storageAuth := auth.NewSharedKeyLiteAuthorizer(accountName, testData.StorageAccountKey)
	messagesClient := NewWithEnvironment(client.Environment)
	messagesClient.Client = client.PrepareWithAuthorizer(messagesClient.Client, storageAuth)

	_, err = queuesClient.Create(ctx, accountName, queueName, map[string]string{})
	if err != nil {
		t.Fatalf("Error creating queue: %s", err)
	}
	defer queuesClient.Delete(ctx, accountName, queueName)

	input := PutInput{
		Message: "ohhai",
	}
	putResp, err := messagesClient.Put(ctx, accountName, queueName, input)
	if err != nil {
		t.Fatalf("Error putting message in queue: %s", err)
	}

	messageId := (*putResp.QueueMessages)[0].MessageId
	popReceipt := (*putResp.QueueMessages)[0].PopReceipt

	_, err = messagesClient.Update(ctx, accountName, queueName, messageId, UpdateInput{
		PopReceipt:        popReceipt,
		Message:           "Updated message",
		VisibilityTimeout: 65,
	})
	if err != nil {
		t.Fatalf("Error updating: %s", err)
	}

	for i := 0; i < 5; i++ {
		input := PutInput{
			Message: fmt.Sprintf("Message %d", i),
		}
		_, err := messagesClient.Put(ctx, accountName, queueName, input)
		if err != nil {
			t.Fatalf("Error putting message %d in queue: %s", i, err)
		}
	}

	peakedMessages, err := messagesClient.Peek(ctx, accountName, queueName, 3)
	if err != nil {
		t.Fatalf("Error peaking messages: %s", err)
	}

	for _, v := range *peakedMessages.QueueMessages {
		t.Logf("Message: %q", v.MessageId)
	}

	retrievedMessages, err := messagesClient.Get(ctx, accountName, queueName, 6, GetInput{})
	if err != nil {
		t.Fatalf("Error retrieving messages: %s", err)
	}

	for _, v := range *retrievedMessages.QueueMessages {
		t.Logf("Message: %q", v.MessageId)

		_, err = messagesClient.Delete(ctx, accountName, queueName, v.MessageId, v.PopReceipt)
		if err != nil {
			t.Fatalf("Error deleting message from queue: %s", err)
		}
	}
}
