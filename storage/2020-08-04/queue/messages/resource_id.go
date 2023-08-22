package messages

import "fmt"

// GetResourceManagerResourceID returns the Resource Manager specific
// ResourceID for a specific Queue

// TODO update for message
func (c Client) GetResourceManagerID(subscriptionID, resourceGroup, accountName, queueName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/queueServices/default/queues/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, queueName)
}
