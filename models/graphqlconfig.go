package models

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
)

type GraphQLConfig struct {
	Addr      string         `toml:"addr"`           // server will start at this "host:port"
	ExtAddr   string         `toml:"extaddr"`        // server will assume running at this base url
	TLSCert   string         `toml:"certificate"`    // TLS Certificate
	TLSKey    string         `toml:"certificatekey"` // TLS Certificate Private Key
	RootCA    []string       `toml:"rootca"`         // additional root CA to trust
	WebStatic string         `toml:"webstatic"`      // folder with static web files
	Logging   *LoggingConfig `toml:"logging"`
	Keycloak  Keycloak       `toml:"keycloak"`
}

func LoadGraphQLConfig(data []byte) (*GraphQLConfig, error) {
	cfg := &GraphQLConfig{
		Addr:    "localhost:4444",
		ExtAddr: "https://localhost:4444",
		TLSCert: "",
		TLSKey:  "",
		RootCA:  []string{},
	}
	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, errors.Wrap(err, "cannot decode toml config")
	}
	return cfg, nil
}
