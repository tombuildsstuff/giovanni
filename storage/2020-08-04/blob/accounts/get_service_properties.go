package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetServicePropertiesResult struct {
	autorest.Response

	ContentType              string
	StorageServiceProperties *StorageServiceProperties
}

// GetServicePropertiesPreparer prepares the GetServiceProperties request.
func (client Client) GetServicePropertiesPreparer(ctx context.Context) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"restype": "service",
		"comp":    "properties",
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.endpoint),
		autorest.WithHeaders(headers),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (client Client) GetServicePropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

func (client Client) GetServicePropertiesResponder(resp *http.Response) (result GetServicePropertiesResult, err error) {
	if resp != nil && resp.Header != nil {
		result.ContentType = resp.Header.Get("Content-Type")
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingXML(&result.StorageServiceProperties),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

func (client Client) GetServiceProperties(ctx context.Context) (result GetServicePropertiesResult, err error) {
	req, err := client.GetServicePropertiesPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.Client", "GetServiceProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetServicePropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "accounts.Client", "GetServiceProperties", resp, "Failure sending request")
		return
	}

	result, err = client.GetServicePropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.Client", "GetServiceProperties", resp, "Failure responding to request")
		return
	}

	return
}
