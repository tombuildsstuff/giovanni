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
)

type Client struct {
	ResourceGroupsClient resources.GroupsClient
	StorageClient        storage.AccountsClient

	auth                *autorest.Authorizer
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

	authClient, env, err := buildAuthClient()
	if err != nil {
		return nil, fmt.Errorf("Error building Auth Client: %s", err)
	}

	if env == nil {
		return nil, fmt.Errorf("Environment was nil: %s", err)
	}

	apiClient, err := buildAPIClient(authClient, *env)
	if err != nil {
		return nil, fmt.Errorf("Error building API Client: %s", err)
	}

	return apiClient, nil
}

func buildAPIClient(config *authentication.Config, env azure.Environment) (*Client, error) {
	oauthConfig, err := config.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", config.TenantID)
	}

	sender := buildSender()
	armAuth, err := config.GetAuthorizationToken(sender, oauthConfig, env.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}

	storageAuth, err := config.GetAuthorizationToken(sender, oauthConfig, "https://storage.azure.com/")
	if err != nil {
		return nil, err
	}

	client := Client{
		AutoRestEnvironment: env,
		auth:                &storageAuth,
	}

	resourceGroupsClient := resources.NewGroupsClientWithBaseURI(env.ResourceManagerEndpoint, config.SubscriptionID)
	resourceGroupsClient.Client = client.PrepareWithStorageResourceManagerAuth(resourceGroupsClient.Client)
	resourceGroupsClient.Authorizer = armAuth
	client.ResourceGroupsClient = resourceGroupsClient

	storageClient := storage.NewAccountsClientWithBaseURI(env.ResourceManagerEndpoint, config.SubscriptionID)
	storageClient.Client = client.PrepareWithStorageResourceManagerAuth(storageClient.Client)
	storageClient.Authorizer = armAuth
	client.StorageClient = storageClient

	return &client, nil
}

func (c Client) PrepareWithStorageResourceManagerAuth(input autorest.Client) autorest.Client {
	return c.PrepareWithAuthorizer(input, *c.auth)
}

func (c Client) PrepareWithAuthorizer(input autorest.Client, authorizer autorest.Authorizer) autorest.Client {
	input.Authorizer = authorizer
	input.Sender = buildSender()
	input.SkipResourceProviderRegistration = true
	return input
}

func buildAuthClient() (*authentication.Config, *azure.Environment, error) {
	builder := &authentication.Builder{
		SubscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
		ClientID:       os.Getenv("ARM_CLIENT_ID"),
		ClientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
		TenantID:       os.Getenv("ARM_TENANT_ID"),
		Environment:    os.Getenv("ARM_ENVIRONMENT"),

		// Feature Toggles
		SupportsClientSecretAuth: true,
	}

	c, err := builder.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("Error building AzureRM Client: %s", err)
	}

	env, err := authentication.DetermineEnvironment(c.Environment)
	if err != nil {
		return nil, nil, err
	}

	return c, env, nil
}
