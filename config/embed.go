package config

import (
	"embed"
)

//go:embed server.toml
var ConfigFS embed.FS
