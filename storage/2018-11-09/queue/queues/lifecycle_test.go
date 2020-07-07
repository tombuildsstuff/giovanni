package queues

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

var _ StorageQueue = Client{}

func TestQueuesLifecycle(t *testing.T) {
	client, err := testhelpers.Build()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	queueName := fmt.Sprintf("queue-%d", testhelpers.RandomInt())

	_, err = client.BuildTestResources(ctx, resourceGroup, accountName, storage.Storage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	queuesClient := NewWithEnvironment(client.Environment)
	queuesClient.Client = client.PrepareWithStorageResourceManagerAuth(queuesClient.Client)

	// first let's test an empty container
	_, err = queuesClient.Create(ctx, accountName, queueName, map[string]string{})
	if err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	// then let's retrieve it to ensure there's no metadata..
	resp, err := queuesClient.GetMetaData(ctx, accountName, queueName)
	if err != nil {
		t.Fatalf("Error retrieving MetaData: %s", err)
	}
	if len(resp.MetaData) != 0 {
		t.Fatalf("Expected no MetaData but got: %s", err)
	}

	// then let's add some..
	updatedMetaData := map[string]string{
		"band":  "panic",
		"boots": "the-overpass",
	}
	_, err = queuesClient.SetMetaData(ctx, accountName, queueName, updatedMetaData)
	if err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	resp, err = queuesClient.GetMetaData(ctx, accountName, queueName)
	if err != nil {
		t.Fatalf("Error re-retrieving MetaData: %s", err)
	}

	if len(resp.MetaData) != 2 {
		t.Fatalf("Expected metadata to have 2 items but got: %s", resp.MetaData)
	}
	if resp.MetaData["band"] != "panic" {
		t.Fatalf("Expected `band` to be `panic` but got: %s", resp.MetaData["band"])
	}
	if resp.MetaData["boots"] != "the-overpass" {
		t.Fatalf("Expected `boots` to be `the-overpass` but got: %s", resp.MetaData["boots"])
	}

	// and woo let's remove it again
	_, err = queuesClient.SetMetaData(ctx, accountName, queueName, map[string]string{})
	if err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	resp, err = queuesClient.GetMetaData(ctx, accountName, queueName)
	if err != nil {
		t.Fatalf("Error retrieving MetaData: %s", err)
	}
	if len(resp.MetaData) != 0 {
		t.Fatalf("Expected no MetaData but got: %s", err)
	}

	// set some properties
	props := StorageServiceProperties{
		Logging: &LoggingConfig{
			Version: "1.0",
			Delete:  true,
			Read:    true,
			Write:   true,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
		Cors: &Cors{
			CorsRule: []CorsRule{
				CorsRule{
					AllowedMethods:  "GET,PUT",
					AllowedOrigins:  "http://www.example.com",
					ExposedHeaders:  "x-tempo-*",
					AllowedHeaders:  "x-tempo-*",
					MaxAgeInSeconds: 500,
				},
				CorsRule{
					AllowedMethods:  "POST",
					AllowedOrigins:  "http://www.test.com",
					ExposedHeaders:  "*",
					AllowedHeaders:  "x-method-*",
					MaxAgeInSeconds: 200,
				},
			},
		},
		HourMetrics: &MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
		MinuteMetrics: &MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
	}
	_, err = queuesClient.SetServiceProperties(ctx, accountName, props)
	if err != nil {
		t.Fatalf("SetServiceProperties failed: %s", err)
	}

	properties, err := queuesClient.GetServiceProperties(ctx, accountName)
	if err != nil {
		t.Fatalf("GetServiceProperties failed: %s", err)
	}

	if len(properties.Cors.CorsRule) > 1 {
		if properties.Cors.CorsRule[0].AllowedMethods != "GET,PUT" {
			t.Fatalf("CORS Methods weren't set!")
		}
		if properties.Cors.CorsRule[1].AllowedMethods != "POST" {
			t.Fatalf("CORS Methods weren't set!")
		}
	} else {
		t.Fatalf("CORS Methods weren't set!")
	}

	if properties.HourMetrics.Enabled {
		t.Fatalf("HourMetrics were enabled when they shouldn't be!")
	}

	if properties.MinuteMetrics.Enabled {
		t.Fatalf("MinuteMetrics were enabled when they shouldn't be!")
	}

	if !properties.Logging.Write {
		t.Fatalf("Logging Write's was not enabled when they should be!")
	}

	includeAPIS := true
	// set some properties
	props2 := StorageServiceProperties{
		Logging: &LoggingConfig{
			Version: "1.0",
			Delete:  true,
			Read:    true,
			Write:   true,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
		Cors: &Cors{
			CorsRule: []CorsRule{
				CorsRule{
					AllowedMethods:  "PUT",
					AllowedOrigins:  "http://www.example.com",
					ExposedHeaders:  "x-tempo-*",
					AllowedHeaders:  "x-tempo-*",
					MaxAgeInSeconds: 500,
				},
			},
		},
		HourMetrics: &MetricsConfig{
			Version: "1.0",
			Enabled: true,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
			IncludeAPIs: &includeAPIS,
		},
		MinuteMetrics: &MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
	}

	_, err = queuesClient.SetServiceProperties(ctx, accountName, props2)
	if err != nil {
		t.Fatalf("SetServiceProperties failed: %s", err)
	}

	properties, err = queuesClient.GetServiceProperties(ctx, accountName)
	if err != nil {
		t.Fatalf("GetServiceProperties failed: %s", err)
	}

	if len(properties.Cors.CorsRule) == 1 {
		if properties.Cors.CorsRule[0].AllowedMethods != "PUT" {
			t.Fatalf("CORS Methods weren't set!")
		}
	} else {
		t.Fatalf("CORS Methods weren't set!")
	}

	if !properties.HourMetrics.Enabled {
		t.Fatalf("HourMetrics were enabled when they shouldn't be!")
	}

	if properties.MinuteMetrics.Enabled {
		t.Fatalf("MinuteMetrics were enabled when they shouldn't be!")
	}

	if !properties.Logging.Write {
		t.Fatalf("Logging Write's was not enabled when they should be!")
	}

	log.Printf("[DEBUG] Deleting..")
	_, err = queuesClient.Delete(ctx, accountName, queueName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error deleting: %s", err))
	}
}
