package handlers

import (
	"net/http"
	"wedding-api/config"
	"wedding-api/dto"
	"wedding-api/models"

	"github.com/gin-gonic/gin"
)

// CreateRsvp tiếp nhận thông tin xác nhận tham dự
func CreateRsvp(c *gin.Context) {
	var req dto.CreateRsvpRequest

	// Validate dữ liệu truyền vào JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Message: "Vui lòng nhập đầy đủ tên và trạng thái tham dự!",
		})
		return
	}

	// Map DTO sang Model DB
	rsvp := models.RsvpResponse{
		Name:       req.Name,
		Message:    req.Message,
		Attendance: req.Attendance,
		GuestCount: req.GuestCount,
		InvitedBy:  req.InvitedBy,
		IsVisible:  false, // Mặc định là ẩn (0) để duyệt lời chúc sau
	}

	// Lưu vào database
	if err := config.DB.Create(&rsvp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Message: "Lưu dữ liệu thất bại, vui lòng thử lại sau!",
		})
		return
	}

	// Trả về response thành công
	c.JSON(http.StatusCreated, dto.ApiResponse{
		Success: true,
		Message: "Xác nhận tham dự thành công!",
	})
}
