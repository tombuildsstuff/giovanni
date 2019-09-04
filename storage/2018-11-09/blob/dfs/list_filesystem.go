package dfs

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// List list filesystems and their properties in given account.
// Parameters:
// prefix - filters results to filesystems within the specified prefix.
// continuation - the number of filesystems returned with each invocation is limited. If the number of
// filesystems to be returned exceeds this limit, a continuation token is returned in the response header
// x-ms-continuation. When a continuation token is  returned in the response, it must be specified in a
// subsequent invocation of the list operation to continue listing the filesystems.
// maxResults - an optional value that specifies the maximum number of items to return. If omitted or greater
// than 5,000, the response will include up to 5,000 items.
// xMsClientRequestID - a UUID recorded in the analytics logs for troubleshooting and correlation.
// timeout - an optional operation timeout value in seconds. The period begins when the request is received by
// the service. If the timeout value elapses before the operation completes, the operation fails.
// xMsDate - specifies the Coordinated Universal Time (UTC) for the request.  This is required when using
// shared key authorization.
func (client Client) ListFilesystem(ctx context.Context, accountName string, prefix string, continuation string, maxResults *int32, xMsClientRequestID string, timeout *int32, xMsDate string) (result FilesystemList, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FilesystemClient.List")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: maxResults,
			Constraints: []validation.Constraint{{Target: "maxResults", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "maxResults", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil}}}}},
		{TargetValue: xMsClientRequestID,
			Constraints: []validation.Constraint{{Target: "xMsClientRequestID", Name: validation.Empty, Rule: false,
				Chain: []validation.Constraint{{Target: "xMsClientRequestID", Name: validation.Pattern, Rule: `^[{(]?[0-9a-f]{8}[-]?([0-9a-f]{4}[-]?){3}[0-9a-f]{12}[)}]?$`, Chain: nil}}}}},
		{TargetValue: timeout,
			Constraints: []validation.Constraint{{Target: "timeout", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "timeout", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("storagedatalake.FilesystemClient", "List", err.Error())
	}

	req, err := client.ListFilesystemPreparer(ctx, accountName, prefix, continuation, maxResults, xMsClientRequestID, timeout, xMsDate)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListFilesystemSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListFilesystemResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client Client) ListFilesystemPreparer(ctx context.Context, accountName string, prefix string, continuation string, maxResults *int32, xMsClientRequestID string, timeout *int32, xMsDate string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"accountName": accountName,
		"dnsSuffix":   DefaultDNSSuffix,
	}

	queryParameters := map[string]interface{}{
		"resource": autorest.Encode("query", "account"),
	}
	if len(prefix) > 0 {
		queryParameters["prefix"] = autorest.Encode("query", prefix)
	}
	if len(continuation) > 0 {
		queryParameters["continuation"] = autorest.Encode("query", continuation)
	}
	if maxResults != nil {
		queryParameters["maxResults"] = autorest.Encode("query", *maxResults)
	}
	if timeout != nil {
		queryParameters["timeout"] = autorest.Encode("query", *timeout)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("https://{accountName}.{dnsSuffix}", urlParameters),
		autorest.WithPath("/"),
		autorest.WithQueryParameters(queryParameters))
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

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client Client) ListFilesystemSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client Client) ListFilesystemResponder(resp *http.Response) (result FilesystemList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
