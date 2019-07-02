## Blob Storage Container SDK for API version 2017-07-29

This package allows you to interact with the Containers Blob Storage API

### Supported Authorizers

* SharedKeyLite (Blob, File & Queue)

Note: when using the `ListBlobs` operation, only `SharedKeyLite` authentication is supported.

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2017-07-29/blob/containers"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    containerName := "mycontainer"
    
    storageAuth := autorest.NewSharedKeyLiteAuthorizer(accountName, storageAccountKey)
    containersClient := containers.New()
    containersClient.Client.Authorizer = storageAuth
    
    ctx := context.TODO()
    createInput := containers.CreateInput{
        AccessLevel: containers.Private,
    }
    if _, err := containersClient.Create(ctx, accountName, containerName, createInput); err != nil {
        return fmt.Errorf("Error creating Container: %s", err)
    }
    
    return nil 
}
```