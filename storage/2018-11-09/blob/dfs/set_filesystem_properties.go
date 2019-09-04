package dfs

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// SetProperties set properties for the filesystem.  This operation supports conditional HTTP requests.  For more
// information, see [Specifying Conditional Headers for Blob Service
// Operations](https://docs.microsoft.com/en-us/rest/api/storageservices/specifying-conditional-headers-for-blob-service-operations).
// Parameters:
// filesystem - the filesystem identifier.  The value must start and end with a letter or number and must
// contain only letters, numbers, and the dash (-) character.  Consecutive dashes are not permitted.  All
// letters must be lowercase.  The value must have between 3 and 63 characters.
// xMsProperties - optional. User-defined properties to be stored with the filesystem, in the format of a
// comma-separated list of name and value pairs "n1=v1, n2=v2, ...", where each value is a base64 encoded
// string. Note that the string may only contain ASCII characters in the ISO-8859-1 character set.  If the
// filesystem exists, any properties not included in the list will be removed.  All properties are removed if
// the header is omitted.  To merge new and existing properties, first get all existing properties and the
// current E-Tag, then make a conditional request with the E-Tag and include values for all properties.
// ifModifiedSince - optional. A date and time value. Specify this header to perform the operation only if the
// resource has been modified since the specified date and time.
// ifUnmodifiedSince - optional. A date and time value. Specify this header to perform the operation only if
// the resource has not been modified since the specified date and time.
// xMsClientRequestID - a UUID recorded in the analytics logs for troubleshooting and correlation.
// timeout - an optional operation timeout value in seconds. The period begins when the request is received by
// the service. If the timeout value elapses before the operation completes, the operation fails.
// xMsDate - specifies the Coordinated Universal Time (UTC) for the request.  This is required when using
// shared key authorization.
func (client Client) SetProperties(ctx context.Context, accountName string, filesystem string, xMsProperties string, ifModifiedSince string, ifUnmodifiedSince string, xMsClientRequestID string, timeout *int32, xMsDate string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FilesystemClient.SetProperties")
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
		return result, validation.NewError("storagedatalake.FilesystemClient", "SetProperties", err.Error())
	}

	req, err := client.SetPropertiesPreparer(ctx, accountName, filesystem, xMsProperties, ifModifiedSince, ifUnmodifiedSince, xMsClientRequestID, timeout, xMsDate)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "SetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetPropertiesSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "SetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagedatalake.FilesystemClient", "SetProperties", resp, "Failure responding to request")
	}

	return
}

// SetPropertiesPreparer prepares the SetProperties request.
func (client Client) SetPropertiesPreparer(ctx context.Context, accountName string, filesystem string, xMsProperties string, ifModifiedSince string, ifUnmodifiedSince string, xMsClientRequestID string, timeout *int32, xMsDate string) (*http.Request, error) {
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
		autorest.AsPatch(),
		autorest.WithCustomBaseURL("https://{accountName}.{dnsSuffix}", urlParameters),
		autorest.WithPathParameters("/{filesystem}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(xMsProperties) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("x-ms-properties", autorest.String(xMsProperties)))
	}
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

// SetPropertiesSender sends the SetProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetPropertiesSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// SetPropertiesResponder handles the response to the SetProperties request. The method always
// closes the http.Response Body.
func (client Client) SetPropertiesResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}
