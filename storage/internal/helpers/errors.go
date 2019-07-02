package helpers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// TODO: trial switching over to this, then upstream it

// WithErrorUnlessStatusCode returns a RespondDecorator that emits an
// azure.RequestError by reading the response body unless the response HTTP status code
// is among the set passed.
//
// If there is a chance service may return responses other than the Azure error
// format and the response cannot be parsed into an error, a decoding error will
// be returned containing the response body. In any case, the Responder will
// return an error if the status code is not satisfied.
//
// If this Responder returns an error, the response body will be replaced with
// an in-memory reader, which needs no further closing.
func WithErrorUnlessStatusCode(codes ...int) autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			err := r.Respond(resp)
			if err == nil && !autorest.ResponseHasStatusCode(resp, codes...) {
				var e azure.RequestError
				defer resp.Body.Close()

				contentType := autorest.EncodedAsJSON
				if resp != nil {
					contentTypeStr := resp.Header.Get("Content-Type")
					if strings.EqualFold(contentTypeStr, "application/xml") {
						contentType = autorest.EncodedAsXML
					}
				}

				// Copy and replace the Body in case it does not contain an error object.
				// This will leave the Body available to the caller.
				b, decodeErr := autorest.CopyAndDecode(contentType, resp.Body, &e)
				resp.Body = ioutil.NopCloser(&b)
				if decodeErr != nil {
					return fmt.Errorf("autorest/azure: error response cannot be parsed: %q error: %v", b.String(), decodeErr)
				}
				if e.ServiceError == nil {
					switch contentType {
					case autorest.EncodedAsJSON:
						// Check if error is unwrapped ServiceError
						if err := json.Unmarshal(b.Bytes(), &e.ServiceError); err != nil {
							return err
						}
						break
					case autorest.EncodedAsXML:
						// Check if error is unwrapped ServiceError
						if err := xml.Unmarshal(b.Bytes(), &e.ServiceError); err != nil {
							return err
						}
						break
					}
				}
				if e.ServiceError.Message == "" {

					rawBody := map[string]interface{}{}

					switch contentType {
					case autorest.EncodedAsJSON:
						// if we're here it means the returned error wasn't OData v4 compliant.
						// try to unmarshal the body as raw JSON in hopes of getting something.
						if err := json.Unmarshal(b.Bytes(), &rawBody); err != nil {
							return err
						}
						break

					case autorest.EncodedAsXML:
						// if we're here it means the returned error wasn't OData v4 compliant.
						// try to unmarshal the body as raw XML in hopes of getting something.
						if err := xml.Unmarshal(b.Bytes(), &rawBody); err != nil {
							return err
						}
						break
					}

					e.ServiceError = &azure.ServiceError{
						Code:    "Unknown",
						Message: "Unknown service error",
					}
					if len(rawBody) > 0 {
						e.ServiceError.Details = []map[string]interface{}{rawBody}
					}
				}
				e.Response = resp
				e.RequestID = azure.ExtractRequestID(resp)
				if e.StatusCode == nil {
					e.StatusCode = resp.StatusCode
				}
				err = &e
			}
			return err
		})
	}
}
