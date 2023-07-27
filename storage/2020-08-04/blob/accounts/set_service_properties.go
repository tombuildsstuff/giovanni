package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// SetServicePropertiesSender sends the SetServiceProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetServicePropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetServicePropertiesPreparer prepares the SetServiceProperties request.
func (client Client) SetServicePropertiesPreparer(ctx context.Context, input StorageServiceProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"restype": "service",
		"comp":    "properties",
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(client.endpoint),
		autorest.WithHeaders(headers),
		autorest.WithXML(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetServicePropertiesResponder handles the response to the SetServiceProperties request. The method always
// closes the http.Response Body.
func (client Client) SetServicePropertiesResponder(resp *http.Response) (result SetServicePropertiesResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

func (client Client) SetServiceProperties(ctx context.Context, input StorageServiceProperties) (result SetServicePropertiesResult, err error) {
	req, err := client.SetServicePropertiesPreparer(ctx, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetServiceProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetServicePropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetServiceProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetServicePropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetServiceProperties", resp, "Failure responding to request")
		return
	}

	return
}
