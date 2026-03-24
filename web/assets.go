package web

import "embed"

// Assets embeds the static frontend so it ships with the Go function bundle.
//
//go:embed index.html styles.css js/*
var Assets embed.FS
