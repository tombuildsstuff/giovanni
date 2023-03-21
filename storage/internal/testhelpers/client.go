package testhelpers

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	authWrapper "github.com/hashicorp/go-azure-sdk/sdk/auth/autorest"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Environment          environments.Environment
	ResourceGroupsClient resources.GroupsClient
	StorageClient        storage.AccountsClient
	SubscriptionId       string

	resourceManagerAuthorizer autorest.Authorizer
	storageAuthorizer         autorest.Authorizer

	AutoRestEnvironment azure.Environment
}

type TestResources struct {
	ResourceGroup      string
	StorageAccountName string
	StorageAccountKey  string
}

func (c Client) BuildTestResources(ctx context.Context, resourceGroup, name string, kind storage.Kind) (*TestResources, error) {
	return c.buildTestResources(ctx, resourceGroup, name, kind, false, "")
}
func (c Client) BuildTestResourcesWithHns(ctx context.Context, resourceGroup, name string, kind storage.Kind) (*TestResources, error) {
	return c.buildTestResources(ctx, resourceGroup, name, kind, true, "")
}
func (c Client) BuildTestResourcesWithSku(ctx context.Context, resourceGroup, name string, kind storage.Kind, sku storage.SkuName) (*TestResources, error) {
	return c.buildTestResources(ctx, resourceGroup, name, kind, false, sku)
}
func (c Client) buildTestResources(ctx context.Context, resourceGroup, name string, kind storage.Kind, enableHns bool, sku storage.SkuName) (*TestResources, error) {
	location := pointer.To(os.Getenv("ARM_TEST_LOCATION"))
	_, err := c.ResourceGroupsClient.CreateOrUpdate(ctx, resourceGroup, resources.Group{
		Location: location,
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating Resource Group %q: %s", resourceGroup, err)
	}

	props := storage.AccountPropertiesCreateParameters{}
	if kind == storage.KindBlobStorage {
		props.AccessTier = storage.AccessTierHot
	}
	if enableHns {
		props.IsHnsEnabled = &enableHns
	}
	if sku == "" {
		sku = storage.SkuNameStandardLRS
	}

	future, err := c.StorageClient.Create(ctx, resourceGroup, name, storage.AccountCreateParameters{
		Location: location,
		Sku: &storage.Sku{
			Name: sku,
		},
		Kind:                              kind,
		AccountPropertiesCreateParameters: &props,
	})

	if err != nil {
		return nil, fmt.Errorf("Error creating Account %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, c.StorageClient.Client)
	if err != nil {
		return nil, fmt.Errorf("Error waiting for the creation of Account %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	keys, err := c.StorageClient.ListKeys(ctx, resourceGroup, name, "")
	if err != nil {
		return nil, fmt.Errorf("Error listing keys for Storage Account %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	// sure we could poll to get around the inconsistency, but where's the fun in that
	time.Sleep(5 * time.Second)

	accountKeys := *keys.Keys
	return &TestResources{
		ResourceGroup:      resourceGroup,
		StorageAccountName: name,
		StorageAccountKey:  *(accountKeys[0]).Value,
	}, nil
}

func (c Client) DestroyTestResources(ctx context.Context, resourceGroup, name string) error {
	_, err := c.StorageClient.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Account %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	future, err := c.ResourceGroupsClient.Delete(ctx, resourceGroup)
	if err != nil {
		return fmt.Errorf("Error deleting Resource Group %q: %s", resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, c.ResourceGroupsClient.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of Resource Group %q: %s", resourceGroup, err)
	}

	return nil
}

func Build(t *testing.T) (*Client, error) {
	if os.Getenv("ACCTEST") == "" {
		t.Skip("Skipping as `ACCTEST` hasn't been set")
	}

	environmentName := os.Getenv("ARM_ENVIRONMENT")
	env, err := environments.FromName(environmentName)
	if err != nil {
		return nil, fmt.Errorf("determing environment %q: %+v", environmentName, err)
	}
	if env == nil {
		return nil, fmt.Errorf("Environment was nil: %s", err)
	}

	autorestEnv, err := authentication.DetermineEnvironment(environmentName)
	if err != nil {
		return nil, fmt.Errorf("determing autorest environment %q: %+v", environmentName, err)
	}
	if env == nil {
		return nil, fmt.Errorf("Autorest Environment was nil: %s", err)
	}

	authConfig := auth.Credentials{
		Environment:  *env,
		ClientID:     os.Getenv("ARM_CLIENT_ID"),
		TenantID:     os.Getenv("ARM_TENANT_ID"),
		ClientSecret: os.Getenv("ARM_CLIENT_SECRET"),

		EnableAuthenticatingUsingClientCertificate: true,
		EnableAuthenticatingUsingClientSecret:      true,
		EnableAuthenticatingUsingAzureCLI:          false,
		EnableAuthenticatingUsingManagedIdentity:   false,
		EnableAuthenticationUsingOIDC:              false,
		EnableAuthenticationUsingGitHubOIDC:        false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	resourceManagerAuth, err := auth.NewAuthorizerFromCredentials(ctx, authConfig, authConfig.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Resource Manager API: %+v", err)
	}

	storageAuthorizer, err := auth.NewAuthorizerFromCredentials(ctx, authConfig, authConfig.Environment.Storage)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Storage API: %+v", err)
	}

	client := Client{
		Environment:    *env,
		SubscriptionId: os.Getenv("ARM_SUBSCRIPTION_ID"),

		// internal
		resourceManagerAuthorizer: authWrapper.AutorestAuthorizer(resourceManagerAuth),
		storageAuthorizer:         authWrapper.AutorestAuthorizer(storageAuthorizer),

		// Legacy / to be removed
		AutoRestEnvironment: *autorestEnv,
	}

	resourceManagerEndpoint, ok := authConfig.Environment.ResourceManager.Endpoint()
	if !ok {
		return nil, fmt.Errorf("Resource Manager Endpoint is not configured for this environment")
	}

	resourceGroupsClient := resources.NewGroupsClientWithBaseURI(*resourceManagerEndpoint, client.SubscriptionId)
	resourceGroupsClient.Client = client.PrepareWithAuthorizer(resourceGroupsClient.Client, client.resourceManagerAuthorizer)
	client.ResourceGroupsClient = resourceGroupsClient

	storageClient := storage.NewAccountsClientWithBaseURI(*resourceManagerEndpoint, client.SubscriptionId)
	storageClient.Client = client.PrepareWithAuthorizer(storageClient.Client, client.resourceManagerAuthorizer)
	client.StorageClient = storageClient

	return &client, nil
}

func (c Client) Configure(client *client.Client, authorizer auth.Authorizer) {
	client.Authorizer = authorizer
	// TODO: add logging
}

func (c Client) PrepareWithStorageResourceManagerAuth(input autorest.Client) autorest.Client {
	return c.PrepareWithAuthorizer(input, c.storageAuthorizer)
}

func (c Client) PrepareWithAuthorizer(input autorest.Client, authorizer autorest.Authorizer) autorest.Client {
	input.Authorizer = authorizer
	input.Sender = buildSender()
	input.SkipResourceProviderRegistration = true
	return input
}
