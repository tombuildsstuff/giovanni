package messages

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// GetResourceID returns the Resource ID for the given Message within a Queue
// This can be useful when, for example, you're using this as a unique identifier
func (client Client) GetResourceID(accountName, queueName, messageID string) string {
	domain := endpoints.GetQueueEndpoint(client.BaseURI, accountName)
	return fmt.Sprintf("%s/%s/messages/%s", domain, queueName, messageID)
}

type ResourceID struct {
	AccountName string
	QueueName   string
	MessageID   string
}

// ParseResourceID parses the specified Resource ID and returns an object
// which can be used to interact with the Message within a Queue
func (client Client) ParseResourceID(id string) (*ResourceID, error) {
	// example: https://account1.queue.core.chinacloudapi.cn/queue1/messages/message1

	if id == "" {
		return nil, fmt.Errorf("`id` was empty")
	}

	uri, err := url.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("Error parsing ID as a URL: %s", err)
	}

	accountName, err := endpoints.GetAccountNameFromEndpoint(uri.Host)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Account Name: %s", err)
	}

	path := strings.TrimPrefix(uri.Path, "/")
	segments := strings.Split(path, "/")
	if len(segments) != 3 {
		return nil, fmt.Errorf("Expected the path to contain 3 segments but got %d", len(segments))
	}

	queueName := segments[0]
	messageID := segments[2]
	return &ResourceID{
		AccountName: *accountName,
		MessageID:   messageID,
		QueueName:   queueName,
	}, nil
}
