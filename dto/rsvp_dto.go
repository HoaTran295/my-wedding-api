package dto

// CreateRsvpRequest binding dữ liệu từ Client Payload
type CreateRsvpRequest struct {
	Name       string `json:"name" binding:"required"`
	Message    string `json:"message"`
	Attendance string `json:"attendance" binding:"required"`
	GuestCount string `json:"guestCount"`
	InvitedBy  string `json:"invitedBy"`
}

// PublicWishResponse DTO cho danh sách lời chúc hiển thị công khai
type PublicWishResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// ApiResponse định dạng Response chung (Đã thêm field Data)
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// UpdateRsvpRequest cập nhật trạng thái 1 phản hồi từ trang quản lý.
// Dùng con trỏ để chỉ cập nhật field nào được gửi lên.
type UpdateRsvpRequest struct {
	IsVisible *bool `json:"is_visible"` // hiển thị lời chúc lên thiệp
	IsSeen    *bool `json:"is_seen"`    // đã xem / chưa xem
}
