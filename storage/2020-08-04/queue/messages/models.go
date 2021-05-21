package messages

import "github.com/Azure/go-autorest/autorest"

type QueueMessage struct {
	MessageText string `xml:"MessageText"`
}

type QueueMessagesListResult struct {
	autorest.Response

	QueueMessages *[]QueueMessageResponse `xml:"QueueMessage"`
}

type QueueMessageResponse struct {
	MessageId       string `xml:"MessageId"`
	InsertionTime   string `xml:"InsertionTime"`
	ExpirationTime  string `xml:"ExpirationTime"`
	PopReceipt      string `xml:"PopReceipt"`
	TimeNextVisible string `xml:"TimeNextVisible"`
}
