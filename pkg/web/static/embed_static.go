package static

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"runtime"
)

//go:embed static favicon.ico
var EmbeddedFS embed.FS

var FS fs.FS = EmbeddedFS

// init sets DiskFS if the templates are available on disk.
func init() {
	// e.g. /path/to/web/templates/templates.go
	_, staticGoPath, _, _ := runtime.Caller(0)
	staticDir := path.Dir(staticGoPath)

	stat, err := os.Stat(staticDir)
	if err != nil || !stat.IsDir() {
		return
	}

	diskFS := os.DirFS(staticDir)
	FS = diskFS
}
