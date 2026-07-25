package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Mogvl/bika/backend/pica"
)

// AuthHandler 认证相关处理器
type AuthHandler struct {
	client *pica.Client
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(client *pica.Client) *AuthHandler {
	return &AuthHandler{client: client}
}

// Login 登录
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "仅支持 POST 请求")
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

	if req.Email == "" || req.Password == "" {
		Error(w, http.StatusBadRequest, "邮箱和密码不能为空")
		return
	}

	resp, err := h.client.Login(req.Email, req.Password)
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 构建返回数据
	data := H{
		"token": h.client.GetToken(),
	}
	if resp.Data != nil {
		// 复制用户信息
		for k, v := range resp.Data {
			if k != "token" {
				data[k] = v
			}
		}
	}

	Success(w, data)
}

// Profile 获取用户信息
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "仅支持 GET 请求")
		return
	}

	resp, err := h.client.GetProfile()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, resp.Data)
}

// PunchIn 每日签到
func (h *AuthHandler) PunchIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "仅支持 POST 请求")
		return
	}

	resp, err := h.client.PunchIn()
	if err != nil {
		// 即使签到失败（已签到）也返回结果
		Success(w, resp.Data)
		return
	}

	Success(w, resp.Data)
}

// AuthMiddleware token 认证中间件
func AuthMiddleware(client *pica.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			// 从 Cookie 中获取
			cookie, err := r.Cookie("token")
			if err == nil && cookie.Value != "" {
				token = cookie.Value
			}
		}
		if token == "" {
			Error(w, http.StatusUnauthorized, "未登录")
			return
		}

		// 添加 Bearer 前缀（如果没有的话）
		if !strings.HasPrefix(token, "Bearer ") {
			token = "Bearer " + token
		}
		client.SetToken(token)
		next(w, r)
	}
}
