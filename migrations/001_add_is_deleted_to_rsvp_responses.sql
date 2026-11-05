-- Thêm cột soft delete cho bảng rsvp_responses
-- Chạy 1 lần trên database "wedding" trước khi deploy API mới.
ALTER TABLE rsvp_responses
  ADD COLUMN is_deleted TINYINT(1) NOT NULL DEFAULT 0 AFTER is_visible,
  ADD INDEX idx_rsvp_responses_is_deleted (is_deleted);
