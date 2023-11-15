package files

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const DefaultCopyPollDuration = 15 * time.Second

// CopyAndWait is a convenience method which doesn't exist in the API, which copies the file and then waits for the copy to complete
func (c Client) CopyAndWait(ctx context.Context, shareName, path, fileName string, input CopyInput, pollDuration time.Duration) (resp CopyResponse, err error) {
	copy, e := c.Copy(ctx, shareName, path, fileName, input)
	if err != nil {
		resp.HttpResponse = copy.HttpResponse
		err = fmt.Errorf("error copying: %s", e)
		return
	}

	resp.CopyID = copy.CopyID

	// since the API doesn't return a LRO, this is a hack which also polls every 10s, but should be sufficient
	for true {
		props, e := c.GetProperties(ctx, shareName, path, fileName)
		if e != nil {
			resp.HttpResponse = copy.HttpResponse
			err = fmt.Errorf("error waiting for copy: %s", e)
			return
		}

		switch strings.ToLower(props.CopyStatus) {
		case "pending":
			time.Sleep(pollDuration)
			continue

		case "success":
			return

		default:
			err = fmt.Errorf("Unexpected CopyState %q", e)
			return
		}
	}

	return
}
