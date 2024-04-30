package main

import (
	"backend/routes"
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
	gin.SetMode(gin.ReleaseMode)
	router.Use(cors.Default())
	router.Use(gin.Logger())
	routes.Routes(router)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	log.Fatal(router.Run(":" + port))
}
