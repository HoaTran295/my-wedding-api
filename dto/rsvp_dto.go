package dto

// CreateRsvpRequest binding dữ liệu từ Client Payload
type CreateRsvpRequest struct {
	Name       string `json:"name" binding:"required"`
	Message    string `json:"message"`
	Attendance string `json:"attendance" binding:"required"`
	GuestCount string `json:"guestCount"`
	InvitedBy  string `json:"invitedBy"`
}

// ApiResponse định dạng Response chung
type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
