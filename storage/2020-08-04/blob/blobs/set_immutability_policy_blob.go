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

type ImmutabilityPolicyMode string

const (
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = "Locked"
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = "Unlocked"
)

type ImmutabilityPolicyBlobInput struct {
	UntilDate  string
	PolicyMode ImmutabilityPolicyMode
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
	if input.UntilDate == "" {
		return result, validation.NewError("blobs.Client", "ImmutabilityPolicyBlob", "`input.UntilDate` cannot be an empty string.")
	}

	req, err := client.SetImmutabilityPolicyBlobPreparer(ctx, accountName, containerName, blobName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobs.Client", "ImmutabilityPolicyBlob", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetImmutabilityPolicyBlobSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "blobs.Client", "ImmutabilityPolicyBlob", resp, "Failure sending request")
		return
	}

	result, err = client.SetImmutabilityPolicyBlobResponder(resp)
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
		"v":    autorest.Encode("query", "2021-06-01"),
	}

	headers := map[string]interface{}{
		"x-ms-version":                        "2021-06-08",
		"x-ms-immutability-policy-until-date": input.UntilDate,
		"x-ms-immutability-policy-mode":       string(ImmutabilityPolicyModeUnlocked),
	}

	if input.PolicyMode != "" {
		headers["x-ms-immutability-policy-mode"] = string(input.PolicyMode)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetBlobEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{containerName}/{blobName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ImmutabilityPolicyBlobSender sends the ImmutabilityPolicyBlob request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetImmutabilityPolicyBlobSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ImmutabilityPolicyBlobResponder handles the response to the ImmutabilityPolicyBlob request. The method always
// closes the http.Response Body.
func (client Client) SetImmutabilityPolicyBlobResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}
	return
}
