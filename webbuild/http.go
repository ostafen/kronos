package statichttp

import "embed"

//go:embed web/*
var Static embed.FS
