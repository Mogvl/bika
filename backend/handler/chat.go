package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mogvl/bika/backend/chat"
)

type ChatHandler struct {
	client *chat.ChatClient
}

func NewChatHandler(client *chat.ChatClient) *ChatHandler {
	return &ChatHandler{client: client}
}

// Login 聊天登录
func (h *ChatHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	token, err := h.client.Login(req.Email, req.Password)
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	Success(w, H{"token": token})
}

// Rooms 获取聊天室列表
func (h *ChatHandler) Rooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.client.GetRooms()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, H{"rooms": rooms})
}

// Messages 获取聊天消息
func (h *ChatHandler) Messages(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomId")
	if roomID == "" {
		Error(w, http.StatusBadRequest, "roomId 不能为空")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	messages, err := h.client.GetMessages(roomID, page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, H{"messages": messages})
}

// SendMessage 发送消息
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}

	var req struct {
		RoomID  string `json:"roomId"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	if req.RoomID == "" || req.Message == "" {
		Error(w, http.StatusBadRequest, "roomId 和 message 不能为空")
		return
	}

	err := h.client.SendMessage(req.RoomID, req.Message)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, H{"message": "发送成功"})
}

// Profile 获取聊天用户信息
func (h *ChatHandler) Profile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.client.GetProfile()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, profile)
}
