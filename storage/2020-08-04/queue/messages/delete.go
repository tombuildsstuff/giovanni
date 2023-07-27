package messages

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
)

// Delete deletes a specific message
func (client Client) Delete(ctx context.Context, queueName, messageID, popReceipt string) (result autorest.Response, err error) {
	if queueName == "" {
		return result, validation.NewError("messages.Client", "Delete", "`queueName` cannot be an empty string.")
	}
	if strings.ToLower(queueName) != queueName {
		return result, validation.NewError("messages.Client", "Delete", "`queueName` must be a lower-cased string.")
	}
	if messageID == "" {
		return result, validation.NewError("messages.Client", "Delete", "`messageID` cannot be an empty string.")
	}
	if popReceipt == "" {
		return result, validation.NewError("messages.Client", "Delete", "`popReceipt` cannot be an empty string.")
	}

	req, err := client.DeletePreparer(ctx, queueName, messageID, popReceipt)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "messages.Client", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messages.Client", "Delete", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client Client) DeletePreparer(ctx context.Context, queueName, messageID, popReceipt string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"queueName": autorest.Encode("path", queueName),
		"messageID": autorest.Encode("path", messageID),
	}

	queryParameters := map[string]interface{}{
		"popreceipt": autorest.Encode("query", popReceipt),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(client.endpoint),
		autorest.WithPathParameters("/{queueName}/messages/{messageID}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client Client) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client Client) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
