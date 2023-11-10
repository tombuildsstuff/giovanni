package blobs

import (
	"context"
	"fmt"
	"time"
)

// CopyAndWait copies a blob to a destination within the storage account and waits for it to finish copying.
func (c Client) CopyAndWait(ctx context.Context, containerName, blobName string, input CopyInput, pollingInterval time.Duration) error {
	if _, err := c.Copy(ctx, containerName, blobName, input); err != nil {
		return fmt.Errorf("error copying: %s", err)
	}

	for true {
		getInput := GetPropertiesInput{
			LeaseID: input.LeaseID,
		}
		getResult, err := c.GetProperties(ctx, containerName, blobName, getInput)
		if err != nil {
			return fmt.Errorf("")
		}

		switch getResult.CopyStatus {
		case Aborted:
			return fmt.Errorf("Copy was aborted: %s", getResult.CopyStatusDescription)

		case Failed:
			return fmt.Errorf("Copy failed: %s", getResult.CopyStatusDescription)

		case Success:
			return nil

		case Pending:
			time.Sleep(pollingInterval)
			continue
		}
	}

	return fmt.Errorf("unexpected error waiting for the copy to complete")
}
