package web

import (
	"net/http"
	"os"
	"sort"

	"github.com/cybernetlab/links/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func Index(c *gin.Context) {
	uid := c.GetUint("user_id")
	links := models.LinksByUserID(uid)
	sort.Slice(links, func(i, j int) bool {
		return links[i].CreatedAt.After(links[j].CreatedAt)
	})
	session := sessions.Default(c)
	errors := session.Flashes("error")
	session.Save()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Links",
		"links":   links,
		"csrf":    csrf.GetToken(c),
		"errors":  errors,
		"baseURL": baseURL(c),
	})
}

func baseURL(c *gin.Context) string {
	base := os.Getenv("SERVER_BASE_URL")
	if base != "" {
		return base
	}
	scheme := "http://"
	if c.Request.TLS != nil {
		scheme = "https://"
	}
	return scheme + c.Request.Host
}
