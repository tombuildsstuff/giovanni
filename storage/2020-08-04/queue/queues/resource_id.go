package queues

import (
	"fmt"
)

// GetResourceManagerID returns the Resource Manager ID for the given Queue
// This can be useful when, for example, you're using this as a unique identifier
func (c Client) GetResourceManagerID(subscriptionID, resourceGroup, accountName, queueName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/queueServices/default/queues/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, queueName)
}
