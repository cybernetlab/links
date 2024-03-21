package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//go:embed templates
var templates embed.FS

//go:embed dist/assets
var assets embed.FS

func Configure(server *gin.Engine) {
	salt := os.Getenv("SERVER_SALT")
	if salt == "" && gin.Mode() == gin.DebugMode {
		salt = "secret"
	}
	store := cookie.NewStore([]byte(salt))
	server.Use(sessions.Sessions("links", store))
	server.Use(csrf.Middleware(csrf.Options{
		Secret: salt,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	server.SetFuncMap(template.FuncMap{
		"formatTime": formatTime,
	})

	configureTemplates(server)
	configureAssets(server)

	server.GET("/login", ShowLogin)
	server.POST("/login", Login)
	server.POST("/logout", Logout)

	needLogin := server.Group("/")
	needLogin.Use(AuthMiddleware)
	needLogin.GET("/", Index)
	needLogin.POST("/private/link", CreateLink)
	needLogin.POST("/private/link/:hash/delete", DeleteLink)

	server.GET("/:hash", GoToLink)
}

func configureAssets(server *gin.Engine) {
	var serverAssets fs.FS
	if gin.Mode() == gin.DebugMode {
		serverAssets = os.DirFS("web/dist/assets")
	} else {
		sub, err := fs.Sub(assets, "dist/assets")
		if err != nil {
			log.Fatal(err)
		}
		serverAssets = sub
	}
	server.StaticFS("/assets", http.FS(serverAssets))
}

func configureTemplates(server *gin.Engine) {
	if gin.Mode() == gin.DebugMode {
		server.LoadHTMLGlob("web/templates/*")
	} else {
		templ := template.Must(template.New("").ParseFS(templates, "*"))
		server.SetHTMLTemplate(templ)
	}
}

func formatTime(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%02d.%02d.%d %02d:%02d:%02d", day, month, year, t.Hour(), t.Minute(), t.Second())
}
