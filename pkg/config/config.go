package config

import (
	"os"
	"strings"

	"go.uber.org/config"
)

// New creates a new config
func New(paths ...string) (config.Provider, error) {
	if len(paths) == 0 {
		if path, ok := os.LookupEnv("ANTENNA_CONFIG"); ok {
			paths = append(paths, path)
		} else {
			paths = append(paths, "config.yml")
		}
	}

	opts := make([]config.YAMLOption, 0, len(paths)+1)
	for _, path := range paths {
		opts = append(opts, config.File(path))
	}
	replacer := strings.NewReplacer(".", "_")
	opts = append(opts, config.Expand(func(key string) (string, bool) {
		key = "ANTENNA_" + replacer.Replace(strings.ToUpper(key))
		return os.LookupEnv(key)
	}))

	return config.NewYAML(opts...)
}
