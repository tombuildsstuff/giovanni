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

func (client Client) DeleteImmutabilityPolicyBlob(ctx context.Context, accountName, containerName, blobName string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("blobs.Client", "ImmutabilityPolicyBlob", "`accountName` cannot be an empty string.")
	}
	if containerName == "" {
		return result, validation.NewError("blobs.Client", "ImmutabilityPolicyBlob", "`containerName` cannot be an empty string.")
	}
	if strings.ToLower(containerName) != containerName {
		return result, validation.NewError("blobs.Client", "ImmutabilityPolicyBlob", "`containerName` must be a lower-cased string.")
	}
	if blobName == "" {
		return result, validation.NewError("blobs.Client", "ImmutabilityPolicyBlob", "`blobName` cannot be an empty string.")
	}

	req, err := client.DeleteImmutabilityPolicyBlobPreparer(ctx, accountName, containerName, blobName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobs.Client", "ImmutabilityPolicyBlob", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteImmutabilityPolicyBlobSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "blobs.Client", "ImmutabilityPolicyBlob", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteImmutabilityPolicyBlobResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobs.Client", "ImmutabilityPolicyBlob", resp, "Failure responding to request")
		return
	}

	return
}

// ImmutabilityPolicyBlobPreparer prepares the ImmutabilityPolicyBlob request.
func (client Client) DeleteImmutabilityPolicyBlobPreparer(ctx context.Context, accountName, containerName, blobName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerName": autorest.Encode("path", containerName),
		"blobName":      autorest.Encode("path", blobName),
	}

	queryParameters := map[string]interface{}{
		"comp": autorest.Encode("query", "immutabilityPolicies"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(endpoints.GetBlobEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{containerName}/{blobName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ImmutabilityPolicyBlobSender sends the ImmutabilityPolicyBlob request. The method will close the
// http.Response Body if it receives an error.
func (client Client) DeleteImmutabilityPolicyBlobSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ImmutabilityPolicyBlobResponder handles the response to the ImmutabilityPolicyBlob request. The method always
// closes the http.Response Body.
func (client Client) DeleteImmutabilityPolicyBlobResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}
	return
}
