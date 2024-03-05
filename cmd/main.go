package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.LoadHTMLGlob("views/html/*.html")

	// router.GET("/", mainpage)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
