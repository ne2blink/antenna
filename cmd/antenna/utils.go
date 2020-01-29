package main

import (
	"encoding/json"
	"os"

	"github.com/ne2blink/antenna/pkg/config"
	"github.com/ne2blink/antenna/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func flagString(cmd *cobra.Command, name string) string {
	return cmd.Flag(name).Value.String()
}

func flagBool(cmd *cobra.Command, name string) bool {
	value, _ := cmd.Flags().GetBool(name)
	return value
}

func getConfig(cmd *cobra.Command) (*viper.Viper, error) {
	return config.New(flagString(cmd, "config"))
}

func getStore(cmd *cobra.Command) (storage.Store, error) {
	config, err := getConfig(cmd)
	if err != nil {
		return nil, err
	}
	return storage.New(
		config.GetString("storage.type"),
		config.GetStringMap("storage.options"),
	)
}

func printJSON(obj interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	return enc.Encode(obj)
}
