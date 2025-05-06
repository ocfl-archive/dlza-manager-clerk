package config

import "embed"

//go:embed clerk.toml
var ConfigFS embed.FS
