package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

// New creates a new config
func New(path string) (*viper.Viper, error) {
	if c, ok := os.LookupEnv("ANTENNA_CONFIG"); ok {
		path = c
	}

	v := viper.New()
	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}
	v.SetEnvPrefix("ANTENNA")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	return v, nil
}
