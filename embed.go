package ip2cloud

import (
	"embed"
	"io/fs"
)

//go:embed data/*.txt
var dataFS embed.FS

func EmbeddedData() (fs.FS, error) {
	return fs.Sub(dataFS, "data")
}
