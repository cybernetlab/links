package web

import (
	"net/http"

	"github.com/cybernetlab/links/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Ref      string `form:"_ref"`
}

func ShowLogin(c *gin.Context) {
	session := sessions.Default(c)
	ref := session.Get("ref")
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Links login",
		"ref":   ref,
		"csrf":  csrf.GetToken(c),
	})
}

func Login(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	user, err := models.Login(form.Username, form.Password)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Delete("ref")
	session.Save()
	url := "/"
	if form.Ref != "" {
		url = form.Ref
	}
	c.Redirect(http.StatusFound, url)
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	session.Delete("ref")
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
