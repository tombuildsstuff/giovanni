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

type ImmutabilityPolicyBlobInput struct {
	UntilDate  *string
	PolicyMode *string
}

func (client Client) SetImmutabilityPolicyBlob(ctx context.Context, accountName, containerName, blobName string, input ImmutabilityPolicyBlobInput) (result autorest.Response, err error) {
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
	if input.UntilDate == nil || *input.UntilDate == "" {
		return result, validation.NewError("blobs.Client", "ImmutabilityPolicyBlob", "`input.UntilDate` cannot be an empty string.")
	}

	req, err := client.SetImmutabilityPolicyBlobPreparer(ctx, accountName, containerName, blobName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobs.Client", "ImmutabilityPolicyBlob", nil, "Failure preparing request")
		return
	}

	resp, err := client.ImmutabilityPolicyBlobSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "blobs.Client", "ImmutabilityPolicyBlob", resp, "Failure sending request")
		return
	}

	result, err = client.ImmutabilityPolicyBlobResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobs.Client", "ImmutabilityPolicyBlob", resp, "Failure responding to request")
		return
	}

	return
}

// ImmutabilityPolicyBlobPreparer prepares the ImmutabilityPolicyBlob request.
func (client Client) SetImmutabilityPolicyBlobPreparer(ctx context.Context, accountName, containerName, blobName string, input ImmutabilityPolicyBlobInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerName": autorest.Encode("path", containerName),
		"blobName":      autorest.Encode("path", blobName),
	}

	queryParameters := map[string]interface{}{
		"comp": autorest.Encode("query", "immutabilityPolicies"),
	}

	headers := map[string]interface{}{
		"x-ms-version":                 APIVersion,
		"x-ms-immutability-until-date": input.UntilDate,
	}

	if *input.PolicyMode != "" {
		headers["x-ms-immutability-policy-mode"] = *input.PolicyMode
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetBlobEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{containerName}/{blobName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
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
func (client Client) ImmutabilityPolicyBlobSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ImmutabilityPolicyBlobResponder handles the response to the ImmutabilityPolicyBlob request. The method always
// closes the http.Response Body.
func (client Client) ImmutabilityPolicyBlobResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}
	return
}
