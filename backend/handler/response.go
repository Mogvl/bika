package handler

import (
	"encoding/json"
	"net/http"
)

// JSON 发送 JSON 响应
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error 发送错误响应
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]any{
		"code":    status,
		"message": message,
		"data":    nil,
	})
}

// H 快捷创建 JSON map
type H map[string]any

// Success 发送成功响应
func Success(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}
