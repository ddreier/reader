package assets

import "embed"

//go:embed dist
//go:embed stylesheets
//go:embed fonts
//go:embed img
var FS embed.FS
