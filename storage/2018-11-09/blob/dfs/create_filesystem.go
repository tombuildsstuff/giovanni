package dfs

import (
	"context"
	"github.com/Azure/go-autorest/tracing"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
)

type CreateFilesystemResponse struct {
	autorest.Response
	Error *ErrorResponse `xml:"Error"`
}

// Create create a filesystem rooted at the specified location. If the filesystem already exists, the operation fails.
// This operation does not support conditional HTTP requests.
// Parameters:
// filesystem - the filesystem identifier.  The value must start and end with a letter or number and must
// contain only letters, numbers, and the dash (-) character.  Consecutive dashes are not permitted.  All
// letters must be lowercase.  The value must have between 3 and 63 characters.
// xMsProperties - user-defined properties to be stored with the filesystem, in the format of a comma-separated
// list of name and value pairs "n1=v1, n2=v2, ...", where each value is a base64 encoded string. Note that the
// string may only contain ASCII characters in the ISO-8859-1 character set.
// xMsClientRequestID - a UUID recorded in the analytics logs for troubleshooting and correlation.
// timeout - an optional operation timeout value in seconds. The period begins when the request is received by
// the service. If the timeout value elapses before the operation completes, the operation fails.
// xMsDate - specifies the Coordinated Universal Time (UTC) for the request.  This is required when using
// shared key authorization.
func (client Client) CreateFilesystem(ctx context.Context, accountName string, filesystem string, xMsProperties string, xMsClientRequestID string, timeout *int32, xMsDate string) (result CreateFilesystemResponse, err error) {
	if accountName == "" {
		return result, validation.NewError("containers.Client", "Create", "`accountName` cannot be an empty string.")
	}
	if filesystem == "" {
		return result, validation.NewError("containers.Client", "Create", "`filesystem` cannot be an empty string.")
	}
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FilesystemClient.Create")
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
		return result, validation.NewError("storagedatalake.FilesystemClient", "Create", err.Error())
	}

	req, err := client.CreateFilesystemPreparer(ctx, accountName, filesystem, xMsProperties, xMsClientRequestID, timeout, xMsDate)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateFilesystemSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateFilesystemResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "Create", resp, "Failure responding to request")
	}

	return
}

// CreatePreparer prepares the Create request.
func (client Client) CreateFilesystemPreparer(ctx context.Context, accountName string, filesystem string, xMsProperties string, xMsClientRequestID string, timeout *int32, xMsDate string) (*http.Request, error) {
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
		autorest.AsPut(),
		autorest.WithCustomBaseURL("https://{accountName}.{dnsSuffix}", urlParameters),
		autorest.WithPathParameters("/{filesystem}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(xMsProperties) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("x-ms-properties", autorest.String(xMsProperties)))
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

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client Client) CreateFilesystemSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client Client) CreateFilesystemResponder(resp *http.Response) (result CreateFilesystemResponse, err error) {
	successfulStatusCodes := []int{
		http.StatusOK, http.StatusCreated,
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
