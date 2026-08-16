package web

import "embed"

// Assets contains the browser bundle and its local visual asset.
//
//go:embed src
var Assets embed.FS
