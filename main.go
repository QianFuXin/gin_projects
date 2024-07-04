package main

import (
	. "gin_projects/config"
	. "gin_projects/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
	r.GET("/redis-set/:key/:value", func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("value")
		err := RDB.Set(Ctx, key, value, 0).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "value set"})
	})
	r.GET("/redis-get/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, err := RDB.Get(Ctx, key).Result()
		if err == redis.Nil {
			c.JSON(404, gin.H{"error": "key not found"})
			return
		} else if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"value": value})
	})
	err = r.Run()
	if err != nil {
		return
	} //listen and serve on 0.0.0.0:8080
}
