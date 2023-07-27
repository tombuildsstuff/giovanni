package messages

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
)

type GetInput struct {
	// VisibilityTimeout specifies the new visibility timeout value, in seconds, relative to server time.
	// The new value must be larger than or equal to 0, and cannot be larger than 7 days.
	VisibilityTimeout *int
}

// Get retrieves one or more messages from the front of the queue
func (client Client) Get(ctx context.Context, queueName string, numberOfMessages int, input GetInput) (result QueueMessagesListResult, err error) {
	if queueName == "" {
		return result, validation.NewError("messages.Client", "Get", "`queueName` cannot be an empty string.")
	}
	if strings.ToLower(queueName) != queueName {
		return result, validation.NewError("messages.Client", "Get", "`queueName` must be a lower-cased string.")
	}
	if numberOfMessages < 1 || numberOfMessages > 32 {
		return result, validation.NewError("messages.Client", "Get", "`numberOfMessages` must be between 1 and 32.")
	}
	if input.VisibilityTimeout != nil {
		t := *input.VisibilityTimeout
		maxTime := (time.Hour * 24 * 7).Seconds()
		if t < 1 || t < int(maxTime) {
			return result, validation.NewError("messages.Client", "Get", "`input.VisibilityTimeout` must be larger than or equal to 1 second, and cannot be larger than 7 days.")
		}
	}

	req, err := client.GetPreparer(ctx, queueName, numberOfMessages, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "messages.Client", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client Client) GetPreparer(ctx context.Context, queueName string, numberOfMessages int, input GetInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"queueName": autorest.Encode("path", queueName),
	}

	queryParameters := map[string]interface{}{
		"numofmessages": autorest.Encode("query", numberOfMessages),
	}

	if input.VisibilityTimeout != nil {
		queryParameters["visibilitytimeout"] = autorest.Encode("query", *input.VisibilityTimeout)
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(client.endpoint),
		autorest.WithPathParameters("/{queueName}/messages", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client Client) GetResponder(resp *http.Response) (result QueueMessagesListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.ByUnmarshallingXML(&result),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
