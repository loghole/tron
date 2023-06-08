package swagger

import "embed"

// Content holds our static swagger web server content.
//
//go:embed favicon-16x16.png
//go:embed favicon-32x32.png
//go:embed index.css
//go:embed index.html
//go:embed oauth2-redirect.html
//go:embed swagger-initializer.js
//go:embed swagger-ui-bundle.js
//go:embed swagger-ui-es-bundle-core.js
//go:embed swagger-ui-es-bundle.js
//go:embed swagger-ui-standalone-preset.js
//go:embed swagger-ui.css
//go:embed swagger-ui.js
//go:embed swagger_data.go
var Content embed.FS
