package main

import (
	"os"

	"github.com/cybernetlab/links/api"
	"github.com/cybernetlab/links/cmd"
	"github.com/cybernetlab/links/models"
	"github.com/cybernetlab/links/web"
	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "create_user" {
		cmd.CreateUser(os.Args[2], os.Args[3])
		return
	}
	address := os.Getenv("SERVER_LISTEN_ADDRESS")
	if address == "" && gin.Mode() == gin.DebugMode {
		address = "127.0.0.1:4000"
	}

	models.Setup()

	gin.ForceConsoleColor()
	router := gin.Default()
	web.Configure(router)
	apiRouter := router.Group("/api")
	apiRouter.Use(api.AuthMiddleware)
	{
		apiRouter.POST("/link", api.CreateLink)
	}
	router.Run(address)
}
