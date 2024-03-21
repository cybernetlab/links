package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("user_id")
	if uid == nil || uid == "" {
		session.Set("ref", c.Request.RequestURI)
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	} else {
		str := fmt.Sprintf("%v", uid)
		id, err := strconv.Atoi(str)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			c.Set("user_id", id)
			c.Next()
		}
	}
}
