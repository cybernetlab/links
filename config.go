package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	//--go:embed dist/assets
	// StaticAssetsFS embed.FS

	//--go:embed templates
	StaticTemplateFS embed.FS
)

func TemplatesFS() fs.FS {
	if gin.Mode() == gin.DebugMode {
		fmt.Println("Using local templates")
		return os.DirFS("templates")
	} else {
		return StaticTemplateFS
	}
}
