package main

import "github.com/gin-gonic/gin"

// TODO: THIS IS WILL BE REPLACED WITH CONTROLLERS AND ROUTERS AS IT SHOULD BE.
func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
