package azure

import (
	"errors"

	azure "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/ne2blink/antenna/pkg/storage"
)

func initTables(tables ...*azure.Table) error {
	for _, table := range tables {
		if err := table.Create(defaultTimeout, azure.NoMetadata, nil); err != nil {
			if err, ok := err.(azure.AzureStorageServiceError); ok && err.Code == "TableAlreadyExists" {
				return nil
			}
			return errors.New(table.Name + ": " + err.Error())
		}
	}
	return nil
}

func queryEntities(table *azure.Table, filter string) ([]*azure.Entity, error) {
	result, err := table.QueryEntities(
		defaultTimeout,
		azure.MinimalMetadata,
		&azure.QueryOptions{
			Filter: filter,
		},
	)
	if err != nil {
		return nil, err
	}

	entities := append([]*azure.Entity(nil), result.Entities...)
	for result.NextLink != nil {
		if result, err = result.NextResults(nil); err != nil {
			return nil, err
		}
		entities = append(entities, result.Entities...)
	}

	return entities, nil
}

func appFromEntity(entity *azure.Entity) storage.App {
	name, _ := entity.Properties["name"].(string)
	secret, _ := entity.Properties["secret"].(string)
	return storage.App{
		ID:     entity.PartitionKey,
		Name:   name,
		Secret: secret,
	}
}

func splitStringSlice(s []string, l int) ([]string, []string) {
	if len(s) > l {
		return s[:l], s[l:]
	}
	return s, nil
}
