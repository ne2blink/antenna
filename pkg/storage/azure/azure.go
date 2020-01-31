package azure

import (
	"errors"

	azure "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/ne2blink/antenna/pkg/storage"
)

const (
	tableNameApp          = "AntennaApp"
	tableNameSubscription = "AntennaSubscription"

	defaultTimeout = 30
)

type store struct {
	app          *azure.Table
	subscription *azure.Table
}

// New creates a store based on Azure Storage.
func New(conn string, shouldInitTables bool) (storage.Store, error) {
	cli, err := azure.NewClientFromConnectionString(conn)
	if err != nil {
		return nil, err
	}

	tableCli := cli.GetTableService()
	tableApp := tableCli.GetTableReference(tableNameApp)
	tableSubscription := tableCli.GetTableReference(tableNameSubscription)

	if shouldInitTables {
		if err := initTables(
			tableApp,
			tableSubscription,
		); err != nil {
			return nil, err
		}
	}

	return &store{
		app:          tableApp,
		subscription: tableSubscription,
	}, nil
}

func init() {
	storage.Register("azure", func(options map[string]interface{}) (storage.Store, error) {
		conn, _ := options["conn"].(string)
		if conn == "" {
			return nil, errors.New("azure: missing connection string")
		}
		notInitTables, _ := options["no_init"].(bool)
		return New(conn, !notInitTables)
	})
}
