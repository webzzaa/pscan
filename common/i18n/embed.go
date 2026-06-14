package i18n

import "embed"

//go:embed locales/*.yaml
var localeFS embed.FS
