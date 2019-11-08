package accounts

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"net/http"
)

type SetServicePropertiesResult struct {
	autorest.Response
	Error *ErrorResponse `xml:"Error"`
}

type ErrorResponse struct {
	Code    *string `xml:"Code"`
	Message *string `xml:"Message"`
}

type StorageServiceProperties struct {
	// Cors - Specifies CORS rules for the Blob service. You can include up to five CorsRule elements in the request. If no CorsRule elements are included in the request body, all CORS rules will be deleted, and CORS will be disabled for the Blob service.
	Cors *CorsRules `json:"cors,omitempty"`
	// DefaultServiceVersion - DefaultServiceVersion indicates the default version to use for requests to the Blob service if an incoming request’s version is not specified. Possible values include version 2008-10-27 and all more recent versions.
	DefaultServiceVersion *string `json:"defaultServiceVersion,omitempty"`
	// DeleteRetentionPolicy - The blob service properties for soft delete.
	DeleteRetentionPolicy *DeleteRetentionPolicy `json:"deleteRetentionPolicy,omitempty"`
	// AutomaticSnapshotPolicyEnabled - Automatic Snapshot is enabled if set to true.
	AutomaticSnapshotPolicyEnabled *bool `json:"automaticSnapshotPolicyEnabled,omitempty"`
	// StaticWebsite - Optional
	StaticWebsite *StaticWebsite
}

type StaticWebsite struct {
	// Enabled - Required. Indicates whether static website support is enabled for the given account.
	Enabled  *bool
	// IndexDocument - Optional. The webpage that Azure Storage serves for requests to the root of a website or any subfolder. For example, index.html. The value is case-sensitive.
	IndexDocument *string
	// ErrorDocument404Path - Optional. The absolute path to a webpage that Azure Storage serves for requests that do not correspond to an existing file. For example, error/404.html. Only a single custom 404 page is supported in each static website. The value is case-sensitive.
	ErrorDocument404Path *string
}

// CorsRules sets the CORS rules. You can include up to five CorsRule elements in the request.
type CorsRules struct {
	// CorsRules - The List of CORS rules. You can include up to five CorsRule elements in the request.
	CorsRules *[]CorsRule `json:"corsRules,omitempty"`
}

// DeleteRetentionPolicy the blob service properties for soft delete.
type DeleteRetentionPolicy struct {
	// Enabled - Indicates whether DeleteRetentionPolicy is enabled for the Blob service.
	Enabled *bool `json:"enabled,omitempty"`
	// Days - Indicates the number of days that the deleted blob should be retained. The minimum specified value can be 1 and the maximum value can be 365.
	Days *int32 `json:"days,omitempty"`
}

// CorsRule specifies a CORS rule for the Blob service.
type CorsRule struct {
	// AllowedOrigins - Required if CorsRule element is present. A list of origin domains that will be allowed via CORS, or "*" to allow all domains
	AllowedOrigins *[]string `json:"allowedOrigins,omitempty"`
	// AllowedMethods - Required if CorsRule element is present. A list of HTTP methods that are allowed to be executed by the origin.
	AllowedMethods *[]string `json:"allowedMethods,omitempty"`
	// MaxAgeInSeconds - Required if CorsRule element is present. The number of seconds that the client/browser should cache a preflight response.
	MaxAgeInSeconds *int32 `json:"maxAgeInSeconds,omitempty"`
	// ExposedHeaders - Required if CorsRule element is present. A list of response headers to expose to CORS clients.
	ExposedHeaders *[]string `json:"exposedHeaders,omitempty"`
	// AllowedHeaders - Required if CorsRule element is present. A list of headers allowed to be part of the cross-origin request.
	AllowedHeaders *[]string `json:"allowedHeaders,omitempty"`
}

// SetPropertiesSender sends the SetProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetServicePropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetPropertiesPreparer prepares the SetProperties request.
func (client Client) SetPropertiesPreparer(ctx context.Context, accountName string, input StorageServiceProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"restype": "service",
		"comp": "properties",
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetBlobEndpoint(client.BaseURI, accountName)),
		autorest.WithHeaders(headers),
		autorest.WithXML(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetPropertiesResponder handles the response to the SetProperties request. The method always
// closes the http.Response Body.
func (client Client) SetPropertiesResponder(resp *http.Response) (result SetServicePropertiesResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}


func (client Client) SetServiceProperties(ctx context.Context, accountName string, input StorageServiceProperties) (result SetServicePropertiesResult, err error) {
	if accountName == "" {
		return result, validation.NewError("accounts.Client", "SetServiceProperties", "`accountName` cannot be an empty string.")
	}

	req, err := client.SetPropertiesPreparer(ctx, accountName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetServicePropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetProperties", resp, "Failure responding to request")
		return
	}

	return
}

