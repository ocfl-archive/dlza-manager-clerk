package config

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	"github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/stashconfig"
	"github.com/ocfl-archive/dlza-manager-clerk/models"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"io/fs"
	"os"
)

type Config struct {
	GraphQLConfig           models.GraphQLConfig `toml:"graphqlconfig"`
	LocalAddr               string               `toml:"localaddr"`
	Domain                  string               `toml:"domain"`
	ExternalAddr            string               `toml:"externaladdr"`
	Bearer                  string               `toml:"bearer"`
	ResolverAddr            string               `toml:"resolveraddr"`
	ResolverTimeout         config.Duration      `toml:"resolvertimeout"`
	ResolverNotFoundTimeout config.Duration      `toml:"resolvernotfoundtimeout"`
	ActionTimeout           config.Duration      `toml:"actiontimeout"`
	ServerTLS               *loader.Config       `toml:"server"`
	ClientTLS               *loader.Config       `toml:"client"`
	GRPCClient              map[string]string    `toml:"grpcclient"`
	Log                     stashconfig.Config   `toml:"log"`
	Jwt                     string               `toml:"jwt"`
}

func LoadConfig(fSys fs.FS, fp string, conf *Config) error {
	if _, err := fs.Stat(fSys, fp); err != nil {
		path, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "cannot get current working directory")
		}
		fSys = os.DirFS(path)
		fp = "clerk.toml"
	}
	data, err := fs.ReadFile(fSys, fp)
	if err != nil {
		return errors.Wrapf(err, "cannot read file [%v] %s", fSys, fp)
	}
	_, err = toml.Decode(string(data), conf)
	if err != nil {
		return errors.Wrapf(err, "error loading config file %v", fp)
	}
	if conf.Jwt == "" {
		conf.Jwt = os.Getenv("JWT_KEY")
	}
	return nil
}
