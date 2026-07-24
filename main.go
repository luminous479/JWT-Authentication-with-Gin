package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/luminous479/JWT-Authentication-with-Gin/routes"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("api_1/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
			"detail":  "Access granted for API 1",
		})
	})

	router.GET("api_2/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
			"detail":  "Access granted for API 2",
		})
	})

	router.Run(":" + port)
}
