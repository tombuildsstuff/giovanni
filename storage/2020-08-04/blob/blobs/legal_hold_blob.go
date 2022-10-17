package blobs

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type LegalHoldBlobInput struct {
	HasLegalHold bool
}

func (client Client) SetLegalHoldBlob(ctx context.Context, accountName, containerName, blobName string, input LegalHoldBlobInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("blobs.Client", "LegalHoldBlob", "`accountName` cannot be an empty string.")
	}
	if containerName == "" {
		return result, validation.NewError("blobs.Client", "LegalHoldBlob", "`containerName` cannot be an empty string.")
	}
	if strings.ToLower(containerName) != containerName {
		return result, validation.NewError("blobs.Client", "LegalHoldBlob", "`containerName` must be a lower-cased string.")
	}
	if blobName == "" {
		return result, validation.NewError("blobs.Client", "LegalHoldBlob", "`blobName` cannot be an empty string.")
	}

	req, err := client.SetLegalHoldBlobPreparer(ctx, accountName, containerName, blobName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobs.Client", "LegalHoldBlob", nil, "Failure preparing request")
		return
	}

	resp, err := client.LegalHoldBlobSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "blobs.Client", "LegalHoldBlob", resp, "Failure sending request")
		return
	}

	result, err = client.LegalHoldBlobResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobs.Client", "LegalHoldBlob", resp, "Failure responding to request")
		return
	}

	return
}

// LegalHoldBlobPreparer prepares the LegalHoldBlob request.
func (client Client) SetLegalHoldBlobPreparer(ctx context.Context, accountName, containerName, blobName string, input LegalHoldBlobInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerName": autorest.Encode("path", containerName),
		"blobName":      autorest.Encode("path", blobName),
	}

	queryParameters := map[string]interface{}{
		"comp": autorest.Encode("query", "legalhold"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	if input.HasLegalHold {
		headers["x-ms-legal-hold"] = true
	} else {
		headers["x-ms-legal-hold"] = false
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetBlobEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{containerName}/{blobName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// LegalHoldBlobSender sends the LegalHoldBlob request. The method will close the
// http.Response Body if it receives an error.
func (client Client) LegalHoldBlobSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// LegalHoldBlobResponder handles the response to the LegalHoldBlob request. The method always
// closes the http.Response Body.
func (client Client) LegalHoldBlobResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}
	return
}
