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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

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
	r.GET("/api/v1/wedding/wishes", handlers.GetPublicWishes) // API lấy lời chúc công khai

	// API quản lý danh sách tham dự (trang ds-tham-du.html)
	r.GET("/api/v1/wedding/admin/rsvp", handlers.ListRsvp)              // danh sách chưa soft delete
	r.DELETE("/api/v1/wedding/admin/rsvp/:id", handlers.SoftDeleteRsvp) // soft delete: is_deleted = 1
	r.PATCH("/api/v1/wedding/admin/rsvp/:id", handlers.UpdateRsvp)      // cập nhật is_visible / is_seen

	// Trang thiệp có tên khách (nginx chuyển index.html?invite-name=... sang đây)
	// để Zalo/Messenger hiện "Thân mời <tên> tham dự đám cưới ..." khi gửi link
	r.GET("/api/v1/wedding/invite-page", handlers.GetInvitePage)

	// Chạy server tại port 8080
	r.Run(":2711")
}
