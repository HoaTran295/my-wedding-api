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

// GetPublicWishes lấy danh sách các lời chúc đã được duyệt (is_visible = true)
func GetPublicWishes(c *gin.Context) {
	var wishes []dto.PublicWishResponse

	// Query lọc: is_visible = 1, message không rỗng, sắp xếp mới nhất lên đầu
	err := config.DB.Model(&models.RsvpResponse{}).
		Select("name, message").
		Where("is_visible = ? AND message IS NOT NULL AND message != ''", true).
		Order("created_at DESC").
		Find(&wishes).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Message: "Khổng thể lấy danh sách lời chúc!",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Lấy danh sách lời chúc thành công!",
		Data:    wishes,
	})
}
