package dfs

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// GetProperties all system and user-defined filesystem properties are specified in the response headers.
// Parameters:
// filesystem - the filesystem identifier.  The value must start and end with a letter or number and must
// contain only letters, numbers, and the dash (-) character.  Consecutive dashes are not permitted.  All
// letters must be lowercase.  The value must have between 3 and 63 characters.
// xMsClientRequestID - a UUID recorded in the analytics logs for troubleshooting and correlation.
// timeout - an optional operation timeout value in seconds. The period begins when the request is received by
// the service. If the timeout value elapses before the operation completes, the operation fails.
// xMsDate - specifies the Coordinated Universal Time (UTC) for the request.  This is required when using
// shared key authorization.
func (client Client) GetFilesystemProperties(ctx context.Context, accountName string, filesystem string, xMsClientRequestID string, timeout *int32, xMsDate string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FilesystemClient.GetProperties")
		defer func() {
			sc := -1
			if result.Response != nil {
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
		return result, validation.NewError("storagedatalake.FilesystemClient", "GetProperties", err.Error())
	}

	req, err := client.GetFilesystemPropertiesPreparer(ctx, accountName, filesystem, xMsClientRequestID, timeout, xMsDate)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "GetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetFilesystemPropertiesSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "GetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.GetFilesystemPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "GetProperties", resp, "Failure responding to request")
	}

	return
}

// GetPropertiesPreparer prepares the GetProperties request.
func (client Client) GetFilesystemPropertiesPreparer(ctx context.Context, accountName string, filesystem string, xMsClientRequestID string, timeout *int32, xMsDate string) (*http.Request, error) {
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
		autorest.AsHead(),
		autorest.WithCustomBaseURL("https://{accountName}.{dnsSuffix}", urlParameters),
		autorest.WithPathParameters("/{filesystem}", pathParameters),
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

// GetPropertiesSender sends the GetProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetFilesystemPropertiesSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// GetPropertiesResponder handles the response to the GetProperties request. The method always
// closes the http.Response Body.
func (client Client) GetFilesystemPropertiesResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}
