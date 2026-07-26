package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Mogvl/bika/backend/pica"
)

// decodeJSON 通用 JSON 解码
func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

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
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	if req.Email == "" || req.Password == "" {
		Error(w, http.StatusBadRequest, "用户名和密码不能为空")
		return
	}

	resp, err := h.client.Login(req.Email, req.Password)
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	data := H{
		"token": h.client.GetToken(),
	}
	if resp.Data != nil {
		for k, v := range resp.Data {
			if k != "token" {
				data[k] = v
			}
		}
	}

	Success(w, data)
}

// Register 注册
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "仅支持 POST 请求")
		return
	}

	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Name      string `json:"name"`
		Birthday  string `json:"birthday"`
		Gender    string `json:"gender"`
		Question1 int    `json:"question1"`
		Question2 int    `json:"question2"`
		Question3 int    `json:"question3"`
		Answer1   string `json:"answer1"`
		Answer2   string `json:"answer2"`
		Answer3   string `json:"answer3"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		Error(w, http.StatusBadRequest, "邮箱、密码和昵称不能为空")
		return
	}

	questions := []int{req.Question1, req.Question2, req.Question3}
	answers := []string{req.Answer1, req.Answer2, req.Answer3}

	resp, err := h.client.Register(req.Email, req.Password, req.Name, req.Birthday, req.Gender, questions, answers)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, resp.Data)
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
		Success(w, resp.Data)
		return
	}

	Success(w, resp.Data)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	resp, err := h.client.ChangePassword(h.client.GetToken(), req.OldPassword, req.NewPassword)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// ForgotPassword 忘记密码
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	resp, err := h.client.ForgotPassword(req.Email)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
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
			cookie, err := r.Cookie("token")
			if err == nil && cookie.Value != "" {
				token = cookie.Value
			}
		}
		if token == "" {
			Error(w, http.StatusUnauthorized, "未登录")
			return
		}

		// PicACG API 的 token 可能已包含 Bearer 前缀
		if !strings.HasPrefix(token, "Bearer ") {
			token = "Bearer " + token
		}
		client.SetToken(token)
		next(w, r)
	}
}
