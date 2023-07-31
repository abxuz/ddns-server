package resource

import (
	"embed"
	"io/fs"
)

//go:embed html
var htmlFs embed.FS

func HtmlFs() (fs.FS, error) {
	return fs.Sub(htmlFs, "html")
}
