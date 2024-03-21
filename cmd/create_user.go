package cmd

import (
	"fmt"
	"log"

	"github.com/cybernetlab/links/models"
)

func CreateUser(username string, password string) {
	if username == "" || password == "" {
		usage()
		return
	}
	models.Setup()
	password = models.EncodePassword(password)
	result := models.DB.Create(&models.User{Username: username, Password: password})
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func usage() {
	fmt.Println("Usage: links create_user username password")
}
