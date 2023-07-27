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

type PutInput struct {
	// A message must be in a format that can be included in an XML request with UTF-8 encoding.
	// The encoded message can be up to 64 KB in size.
	Message string

	// The maximum time-to-live can be any positive number,
	// as well as -1 indicating that the message does not expire.
	// If this parameter is omitted, the default time-to-live is 7 days.
	MessageTtl *int

	// Specifies the new visibility timeout value, in seconds, relative to server time.
	// The new value must be larger than or equal to 0, and cannot be larger than 7 days.
	// The visibility timeout of a message cannot be set to a value later than the expiry time.
	// visibilitytimeout should be set to a value smaller than the time-to-live value.
	// If not specified, the default value is 0.
	VisibilityTimeout *int
}

// Put adds a new message to the back of the message queue
func (client Client) Put(ctx context.Context, accountName, queueName string, input PutInput) (result QueueMessagesListResult, err error) {
	if accountName == "" {
		return result, validation.NewError("messages.Client", "Put", "`accountName` cannot be an empty string.")
	}
	if queueName == "" {
		return result, validation.NewError("messages.Client", "Put", "`queueName` cannot be an empty string.")
	}
	if strings.ToLower(queueName) != queueName {
		return result, validation.NewError("messages.Client", "Put", "`queueName` must be a lower-cased string.")
	}

	req, err := client.PutPreparer(ctx, accountName, queueName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Put", nil, "Failure preparing request")
		return
	}

	resp, err := client.PutSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "messages.Client", "Put", resp, "Failure sending request")
		return
	}

	result, err = client.PutResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Put", resp, "Failure responding to request")
		return
	}

	return
}

// PutPreparer prepares the Put request.
func (client Client) PutPreparer(ctx context.Context, accountName, queueName string, input PutInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"queueName": autorest.Encode("path", queueName),
	}

	queryParameters := map[string]interface{}{}

	if input.MessageTtl != nil {
		queryParameters["messagettl"] = autorest.Encode("path", *input.MessageTtl)
	}

	if input.VisibilityTimeout != nil {
		queryParameters["visibilitytimeout"] = autorest.Encode("path", *input.VisibilityTimeout)
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	body := QueueMessage{
		MessageText: input.Message,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(endpoints.GetOrBuildQueueEndpoint(client.endpoint, client.BaseURI, accountName)),
		autorest.WithPathParameters("/{queueName}/messages", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithXML(body),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PutSender sends the Put request. The method will close the
// http.Response Body if it receives an error.
func (client Client) PutSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// PutResponder handles the response to the Put request. The method always
// closes the http.Response Body.
func (client Client) PutResponder(resp *http.Response) (result QueueMessagesListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.ByUnmarshallingXML(&result),
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
