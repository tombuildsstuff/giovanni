package messages

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// Peek retrieves one or more messages from the front of the queue, but doesn't alter the visibility of the messages
func (client Client) Peek(ctx context.Context, accountName, queueName string, numberOfMessages int) (result QueueMessagesListResult, err error) {
	if accountName == "" {
		return result, validation.NewError("messages.Client", "Peek", "`accountName` cannot be an empty string.")
	}
	if queueName == "" {
		return result, validation.NewError("messages.Client", "Peek", "`queueName` cannot be an empty string.")
	}
	if strings.ToLower(queueName) != queueName {
		return result, validation.NewError("messages.Client", "Peek", "`queueName` must be a lower-cased string.")
	}
	if numberOfMessages < 1 || numberOfMessages > 32 {
		return result, validation.NewError("messages.Client", "Peek", "`numberOfMessages` must be between 1 and 32.")
	}

	req, err := client.PeekPreparer(ctx, accountName, queueName, numberOfMessages)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Peek", nil, "Failure preparing request")
		return
	}

	resp, err := client.PeekSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "messages.Client", "Peek", resp, "Failure sending request")
		return
	}

	result, err = client.PeekResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Peek", resp, "Failure responding to request")
		return
	}

	return
}

// PeekPreparer prepares the Peek request.
func (client Client) PeekPreparer(ctx context.Context, accountName, queueName string, numberOfMessages int) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"queueName": autorest.Encode("path", queueName),
	}

	queryParameters := map[string]interface{}{
		"numofmessages": autorest.Encode("query", numberOfMessages),
		"peekonly":      autorest.Encode("query", true),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetQueueEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{queueName}/messages", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PeekSender sends the Peek request. The method will close the
// http.Response Body if it receives an error.
func (client Client) PeekSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// PeekResponder handles the response to the Peek request. The method always
// closes the http.Response Body.
func (client Client) PeekResponder(resp *http.Response) (result QueueMessagesListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.ByUnmarshallingXML(&result),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
