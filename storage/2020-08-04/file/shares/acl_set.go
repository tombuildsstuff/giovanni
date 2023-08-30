package shares

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type setAclResponse struct {
	HttpResponse *client.Response
}

type SetAclInput struct {
	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`

	XMLName xml.Name `xml:"SignedIdentifiers"`
}

// SetACL sets the specified Access Control List on the specified Storage Share
func (c Client) SetACL(ctx context.Context, shareName string, input SetAclInput) (resp setAclResponse, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}
	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: setAclOptions{},
		Path:          shareName,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	err = req.Marshal(&input)
	if err != nil {
		return resp, fmt.Errorf("marshalling request: %v", err)
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}
	return
}

type setAclOptions struct {
	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
}

func (s setAclOptions) ToHeaders() *client.Headers {
	return nil
}

func (s setAclOptions) ToOData() *odata.Query {
	return nil
}

func (s setAclOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("comp", "acl")
	return out
}
