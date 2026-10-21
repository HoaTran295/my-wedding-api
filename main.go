package main

import (
	"wedding-api/config"
	"wedding-api/handlers"

	"github.com/gin-gonic/gin"
)

// Middleware CORS cho phép Frontend HTML call API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// Kết nối Database
	config.ConnectDatabase()

	r := gin.Default()

	// Sử dụng middleware CORS
	r.Use(CORSMiddleware())

	// Route API
	r.POST("/api/v1/wedding/rsvp", handlers.CreateRsvp)

	// Chạy server tại port 8080
	r.Run(":2711")
}
