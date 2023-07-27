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

type UpdateInput struct {
	// A message must be in a format that can be included in an XML request with UTF-8 encoding.
	// The encoded message can be up to 64 KB in size.
	Message string

	// Specifies the valid pop receipt value required to modify this message.
	PopReceipt string

	// Specifies the new visibility timeout value, in seconds, relative to server time.
	// The new value must be larger than or equal to 0, and cannot be larger than 7 days.
	// The visibility timeout of a message cannot be set to a value later than the expiry time.
	// A message can be updated until it has been deleted or has expired.
	VisibilityTimeout int
}

// Update updates an existing message based on it's Pop Receipt
func (client Client) Update(ctx context.Context, accountName, queueName string, messageID string, input UpdateInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("messages.Client", "Update", "`accountName` cannot be an empty string.")
	}
	if queueName == "" {
		return result, validation.NewError("messages.Client", "Update", "`queueName` cannot be an empty string.")
	}
	if strings.ToLower(queueName) != queueName {
		return result, validation.NewError("messages.Client", "Update", "`queueName` must be a lower-cased string.")
	}
	if input.PopReceipt == "" {
		return result, validation.NewError("messages.Client", "Update", "`input.PopReceipt` cannot be an empty string.")
	}

	req, err := client.UpdatePreparer(ctx, accountName, queueName, messageID, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "messages.Client", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Update", resp, "Failure responding to request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client Client) UpdatePreparer(ctx context.Context, accountName, queueName string, messageID string, input UpdateInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"queueName": autorest.Encode("path", queueName),
		"messageID": autorest.Encode("path", messageID),
	}

	queryParameters := map[string]interface{}{
		"popreceipt":        autorest.Encode("query", input.PopReceipt),
		"visibilitytimeout": autorest.Encode("query", input.VisibilityTimeout),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	body := QueueMessage{
		MessageText: input.Message,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetOrBuildQueueEndpoint(client.endpoint, client.BaseURI, accountName)),
		autorest.WithPathParameters("/{queueName}/messages/{messageID}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithXML(body),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client Client) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client Client) UpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
