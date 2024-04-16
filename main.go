package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(cors.Default())
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	log.Fatal(router.Run(":" + port))
}