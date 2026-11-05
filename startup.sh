#!/usr/bin/env bash
# Khởi động wedding-api: build 1 lần ra file chạy rồi chạy file đó.
# KHÔNG dùng "go run": lệnh đó biên dịch lại toàn bộ code + thư viện mỗi lần chạy
# (tốn rất nhiều CPU/RAM) và giữ thêm 1 tiến trình "go" chạy kèm.
#
# Cách dùng:  ./startup.sh          (build + khởi động lại)
#             ./startup.sh stop     (dừng)
set -euo pipefail
cd "$(dirname "$0")"

APP=./wedding-api-server
PID_FILE=./wedding-api.pid
LOG_FILE=./wedding-api.log

stop_app() {
  if [[ -f "$PID_FILE" ]] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
    echo "Dừng tiến trình cũ (PID $(cat "$PID_FILE"))..."
    kill "$(cat "$PID_FILE")"
    for _ in $(seq 1 20); do kill -0 "$(cat "$PID_FILE")" 2>/dev/null || break; sleep 0.5; done
  fi
  rm -f "$PID_FILE"
  # Dọn các bản cũ chạy bằng "go run main.go" nếu còn sót
  pkill -f "go run main.go" 2>/dev/null || true
}

if [[ "${1:-}" == "stop" ]]; then
  stop_app
  echo "Đã dừng."
  exit 0
fi

echo "Build..."
# Giới hạn số luồng biên dịch để VPS nhỏ không bị treo trong lúc build
GOMAXPROCS=1 go build -p 1 -trimpath -ldflags="-s -w" -o "$APP" .

stop_app

echo "Khởi động..."
# GIN_MODE=release: tắt log debug của Gin
# GOMEMLIMIT: giới hạn mềm bộ nhớ của Go runtime (GC chạy sớm hơn khi gần chạm ngưỡng)
GIN_MODE=release GOMEMLIMIT=100MiB nohup "$APP" >>"$LOG_FILE" 2>&1 &
echo $! > "$PID_FILE"
sleep 1
if kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
  echo "Đang chạy (PID $(cat "$PID_FILE")). Log: $LOG_FILE"
else
  echo "Khởi động thất bại, xem log:"; tail -n 30 "$LOG_FILE"; exit 1
fi
