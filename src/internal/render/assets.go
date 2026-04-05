package render

import "embed"

//go:embed templates/report.html.tmpl templates/report.css
var embeddedAssets embed.FS
