-- Thêm cột đánh dấu "đã xem" cho trang ds-tham-du.html (0 = tin mới, 1 = đã xem)
-- Chạy 1 lần trên database "wedding" trước khi deploy API mới (sau migration 001).
ALTER TABLE rsvp_responses
  ADD COLUMN is_seen TINYINT(1) NOT NULL DEFAULT 0 AFTER is_deleted;
