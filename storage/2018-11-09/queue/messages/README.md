## Queue Storage Messages SDK for API version 2018-11-09

This package allows you to interact with the Messages Queue Storage API

### Supported Authorizers

* Azure Active Directory (for the Resource Endpoint `https://storage.azure.com`)
* SharedKeyLite (Blob, File & Queue)

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/queue/messages"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    queueName := "myqueue"
    
    storageAuth := autorest.NewSharedKeyLiteAuthorizer(accountName, storageAccountKey)
    messagesClient := messages.New()
    messagesClient.Client.Authorizer = storageAuth
    
    ctx := context.TODO()
    input := messages.PutInput{
    	Message: "<over><message>hello</message></over>",
    }
    if _, err := messagesClient.Put(ctx, accountName, queueName, input); err != nil {
        return fmt.Errorf("Error creating Message: %s", err)
    }
    
    return nil 
}
```