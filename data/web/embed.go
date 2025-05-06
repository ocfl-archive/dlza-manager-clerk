package web

import "embed"

//go:embed static/*
var WebFS embed.FS
