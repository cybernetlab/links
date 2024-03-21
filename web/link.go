package web

import (
	"net/http"
	"net/url"
	"time"

	"github.com/cybernetlab/links/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type createLinkForm struct {
	Link string `form:"link" binding:"required"`
}

type hashParams struct {
	Hash string `uri:"hash" binding:"required"`
}

func CreateLink(c *gin.Context) {
	var form createLinkForm
	if err := c.ShouldBind(&form); err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	target, err := url.Parse(form.Link)
	if err != nil {
		abortWithError(c, "Invalid link")
		return
	}
	if target.Scheme != "http" && target.Scheme != "https" {
		abortWithError(c, "Link should start with 'http://' or 'https://'")
		return
	}
	if target.Host == c.Request.Host {
		abortWithError(c, "Can't create link to myself")
		return
	}
	hours, _ := time.ParseDuration("1h")
	uid := c.GetUint("user_id")
	models.CreateLink(uid, form.Link, hours)
	c.Redirect(http.StatusFound, "/")
}

func DeleteLink(c *gin.Context) {
	var params hashParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	uid := c.GetUint("user_id")
	models.DeleteLink(uid, params.Hash)
	c.Redirect(http.StatusFound, "/")
}

func GoToLink(c *gin.Context) {
	var params hashParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	link := models.LinkByHash(params.Hash)
	if link == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	link.Count += 1
	models.UpdateLink(link)
	c.Redirect(http.StatusTemporaryRedirect, link.Link)
}

func abortWithError(c *gin.Context, err string) {
	session := sessions.Default(c)
	session.AddFlash(err, "error")
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
