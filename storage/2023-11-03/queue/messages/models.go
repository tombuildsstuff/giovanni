package messages

import (
	"net/http"
)

type QueueMessage struct {
	MessageText string `xml:"MessageText"`
}

type QueueMessagesListResponse struct {
	HttpResponse *http.Response

	QueueMessages *[]QueueMessageResponse `xml:"QueueMessage"`
}

type QueueMessageResponse struct {
	MessageId       string `xml:"MessageId"`
	InsertionTime   string `xml:"InsertionTime"`
	ExpirationTime  string `xml:"ExpirationTime"`
	PopReceipt      string `xml:"PopReceipt"`
	TimeNextVisible string `xml:"TimeNextVisible"`
}
