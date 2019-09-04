package dfs

import (
	"context"
	"github.com/Azure/go-autorest/tracing"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
)

type DeleteFilesystemResponse struct {
	autorest.Response
	Error *ErrorResponse `xml:"Error"`
}

// Delete marks the filesystem for deletion.  When a filesystem is deleted, a filesystem with the same identifier
// cannot be created for at least 30 seconds. While the filesystem is being deleted, attempts to create a filesystem
// with the same identifier will fail with status code 409 (Conflict), with the service returning additional error
// information indicating that the filesystem is being deleted. All other operations, including operations on any files
// or directories within the filesystem, will fail with status code 404 (Not Found) while the filesystem is being
// deleted. This operation supports conditional HTTP requests.  For more information, see [Specifying Conditional
// Headers for Blob Service
// Operations](https://docs.microsoft.com/en-us/rest/api/storageservices/specifying-conditional-headers-for-blob-service-operations).
// Parameters:
// filesystem - the filesystem identifier.  The value must start and end with a letter or number and must
// contain only letters, numbers, and the dash (-) character.  Consecutive dashes are not permitted.  All
// letters must be lowercase.  The value must have between 3 and 63 characters.
// ifModifiedSince - optional. A date and time value. Specify this header to perform the operation only if the
// resource has been modified since the specified date and time.
// ifUnmodifiedSince - optional. A date and time value. Specify this header to perform the operation only if
// the resource has not been modified since the specified date and time.
// xMsClientRequestID - a UUID recorded in the analytics logs for troubleshooting and correlation.
// timeout - an optional operation timeout value in seconds. The period begins when the request is received by
// the service. If the timeout value elapses before the operation completes, the operation fails.
// xMsDate - specifies the Coordinated Universal Time (UTC) for the request.  This is required when using
// shared key authorization.
func (client Client) DeleteFilesystem(ctx context.Context, accountName string, filesystem string, ifModifiedSince string, ifUnmodifiedSince string, xMsClientRequestID string, timeout *int32, xMsDate string) (result DeleteFilesystemResponse, err error) {
	if accountName == "" {
		return result, validation.NewError("containers.Client", "Create", "`accountName` cannot be an empty string.")
	}
	if filesystem == "" {
		return result, validation.NewError("containers.Client", "Create", "`filesystem` cannot be an empty string.")
	}
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FilesystemClient.Delete")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: filesystem,
			Constraints: []validation.Constraint{{Target: "filesystem", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "filesystem", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: xMsClientRequestID,
			Constraints: []validation.Constraint{{Target: "xMsClientRequestID", Name: validation.Empty, Rule: false,
				Chain: []validation.Constraint{{Target: "xMsClientRequestID", Name: validation.Pattern, Rule: `^[{(]?[0-9a-f]{8}[-]?([0-9a-f]{4}[-]?){3}[0-9a-f]{12}[)}]?$`, Chain: nil}}}}},
		{TargetValue: timeout,
			Constraints: []validation.Constraint{{Target: "timeout", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "timeout", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("storagedatalake.FilesystemClient", "Delete", err.Error())
	}

	req, err := client.DeleteFilesystemPreparer(ctx, accountName, filesystem, ifModifiedSince, ifUnmodifiedSince, xMsClientRequestID, timeout, xMsDate)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteFilesystemSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteFilesystemResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client Client) DeleteFilesystemPreparer(ctx context.Context, accountName string, filesystem string, ifModifiedSince string, ifUnmodifiedSince string, xMsClientRequestID string, timeout *int32, xMsDate string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"accountName": accountName,
		"dnsSuffix":   DefaultDNSSuffix,
	}

	pathParameters := map[string]interface{}{
		"filesystem": autorest.Encode("path", filesystem),
	}

	queryParameters := map[string]interface{}{
		"resource": autorest.Encode("query", "filesystem"),
	}
	if timeout != nil {
		queryParameters["timeout"] = autorest.Encode("query", *timeout)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithCustomBaseURL("https://{accountName}.{dnsSuffix}", urlParameters),
		autorest.WithPathParameters("/{filesystem}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(ifModifiedSince) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("If-Modified-Since", autorest.String(ifModifiedSince)))
	}
	if len(ifUnmodifiedSince) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("If-Unmodified-Since", autorest.String(ifUnmodifiedSince)))
	}
	if len(xMsClientRequestID) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("x-ms-client-request-id", autorest.String(xMsClientRequestID)))
	}
	if len(xMsDate) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("x-ms-date", autorest.String(xMsDate)))
	}
	preparer = autorest.DecoratePreparer(preparer,
		autorest.WithHeader("x-ms-version", autorest.String(APIVersion)))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client Client) DeleteFilesystemSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client Client) DeleteFilesystemResponder(resp *http.Response) (result DeleteFilesystemResponse, err error) {
	successfulStatusCodes := []int{
		http.StatusOK, http.StatusAccepted,
	}
	if autorest.ResponseHasStatusCode(resp, successfulStatusCodes...) {
		// when successful there's no response
		err = autorest.Respond(
			resp,
			client.ByInspecting(),
			azure.WithErrorUnlessStatusCode(successfulStatusCodes...),
			autorest.ByClosing())
		result.Response = autorest.Response{Response: resp}
	} else {
		// however when there's an error the error's in the response
		err = autorest.Respond(
			resp,
			client.ByInspecting(),
			azure.WithErrorUnlessStatusCode(successfulStatusCodes...),
			autorest.ByUnmarshallingXML(&result),
			autorest.ByClosing())
		result.Response = autorest.Response{Response: resp}
	}
	return
}
