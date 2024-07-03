package main

import (
	. "gin_projects/config"
	. "gin_projects/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	InitDB()

	r := gin.Default()
	// Enable CORS
	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	// Setup User routes
	SetupUserRoutes(r)
	err = r.Run()
	if err != nil {
		return
	} //listen and serve on 0.0.0.0:8080
}
