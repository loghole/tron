package swagger

import "embed"

// Content holds our static swagger web server content.
//go:embed favicon-16x16.png
//go:embed favicon-32x32.png
//go:embed index.html
//go:embed oauth2-redirect.html
//go:embed swagger-ui-bundle.js
//go:embed swagger-ui-es-bundle-core.js
//go:embed swagger-ui-es-bundle.js
//go:embed swagger-ui-standalone-preset.js
//go:embed swagger-ui.css
//go:embed swagger-ui.js
var Content embed.FS // nolint:gochecknoglobals // is swagger web server embed content.
