package models

import "time"

type RsvpResponse struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Message     string    `gorm:"type:text" json:"message"`
	Attendance  string    `gorm:"type:varchar(100);not null" json:"attendance"`
	GuestCount  string    `gorm:"type:varchar(255)" json:"guest_count"`
	InvitedBy   string    `gorm:"type:varchar(100)" json:"invited_by"`
	IsVisible   bool      `gorm:"type:tinyint(1);default:0;not null" json:"is_visible"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName chỉ định tên bảng trong database MySQL
func (RsvpResponse) TableName() string {
	return "rsvp_responses"
}
