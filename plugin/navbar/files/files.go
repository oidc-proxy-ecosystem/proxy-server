package files

import "embed"

//go:embed templates
var Templates embed.FS

//go:embed images
var Images embed.FS
