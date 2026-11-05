package handlers

import (
	"net/http"
	"strconv"
	"wedding-api/config"
	"wedding-api/dto"
	"wedding-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// Giá trị attendance khi khách xác nhận tham dự (khớp với option trong form thiệp)
const attendanceYes = "Tôi sẽ tham dự"

// ListRsvp trả về danh sách phản hồi chưa bị soft delete.
// Sắp xếp: khách xác nhận tham dự lên trước, sau đó mới nhất lên đầu.
func ListRsvp(c *gin.Context) {
	var rows []models.RsvpResponse

	err := config.DB.
		Where("is_deleted = ?", false).
		// Gộp 1 biểu thức: nếu gọi Order() 2 lần, GORM ghi đè mất phần CASE
		Order(clause.OrderBy{Expression: clause.Expr{
			SQL:                "CASE WHEN attendance = ? THEN 0 ELSE 1 END, created_at DESC",
			Vars:               []interface{}{attendanceYes},
			WithoutParentheses: true,
		}}).
		Find(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Message: "Không thể lấy danh sách khách mời!",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Lấy danh sách khách mời thành công!",
		Data:    rows,
	})
}

// SoftDeleteRsvp đánh dấu is_deleted = 1 cho 1 phản hồi (không xoá hẳn khỏi database)
func SoftDeleteRsvp(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Message: "ID không hợp lệ",
		})
		return
	}

	result := config.DB.Model(&models.RsvpResponse{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("is_deleted", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Message: "Xoá thất bại, vui lòng thử lại sau!",
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Message: "Không tìm thấy khách mời (hoặc đã xoá trước đó)",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Đã xoá khỏi danh sách",
	})
}
