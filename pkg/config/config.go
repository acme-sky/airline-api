package config

import (
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"strings"
)

// Load config froom environment. Something different than that could be create
// an overthinking of the structure for a container because we should also
// consider volumes to insert config file.
// Every env var is coverted to lowercase and plitted by underscore "_".
//
// Example: `DATABASE_DSN` becomes `database.dsn`
func LoadConfig() (*koanf.Koanf, error) {
	k := koanf.New(".")

	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	return k, nil
}
