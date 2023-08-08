package paths

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type pathOptions struct {
}

func (p pathOptions) ToHeaders() *client.Headers {
	return nil
}

func (p pathOptions) ToOData() *odata.Query {
	return nil
}

func (p pathOptions) ToQuery() *client.QueryParams {
	return nil
}
