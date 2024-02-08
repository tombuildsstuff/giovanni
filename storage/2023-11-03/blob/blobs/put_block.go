package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type PutBlockInput struct {
	BlockID    string
	Content    []byte
	ContentMD5 *string
	LeaseID    *string
}

type PutBlockResponse struct {
	HttpResponse *http.Response

	ContentMD5 string
}

// PutBlock creates a new block to be committed as part of a blob.
func (c Client) PutBlock(ctx context.Context, containerName, blobName string, input PutBlockInput) (result PutBlockResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}

	if strings.ToLower(containerName) != containerName {
		err = fmt.Errorf("`containerName` must be a lower-cased string")
		return
	}

	if blobName == "" {
		err = fmt.Errorf("`blobName` cannot be an empty string")
		return
	}

	if input.BlockID == "" {
		err = fmt.Errorf("`input.BlockID` cannot be an empty string")
		return
	}

	if len(input.Content) == 0 {
		err = fmt.Errorf("`input.Content` cannot be empty")
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putBlockOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	err = req.Marshal(&input.Content)
	if err != nil {
		err = fmt.Errorf("marshalling request: %+v", err)
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.HttpResponse = resp.Response
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type putBlockOptions struct {
	input PutBlockInput
}

func (p putBlockOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Content-Length", strconv.Itoa(len(p.input.Content)))

	if p.input.ContentMD5 != nil {
		headers.Append("x-ms-blob-content-md5", *p.input.ContentMD5)
	}
	if p.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *p.input.LeaseID)
	}

	return headers
}

func (p putBlockOptions) ToOData() *odata.Query {
	return nil
}

func (p putBlockOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "block")
	out.Append("blockid", p.input.BlockID)
	return out
}
