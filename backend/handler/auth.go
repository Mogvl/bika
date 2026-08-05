package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mogvl/bika/backend/fried"
	"github.com/Mogvl/bika/backend/pica"
)

// decodeJSON 通用 JSON 解码
func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// AuthHandler 认证相关处理器
type AuthHandler struct {
	client      *pica.Client
	friedClient *fried.FriedClient
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(client *pica.Client, friedClient *fried.FriedClient) *AuthHandler {
	return &AuthHandler{client: client, friedClient: friedClient}
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

	// 同步 token 到锅贴客户端
	if h.friedClient != nil {
		h.friedClient.SetToken(h.client.GetToken())
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

// ResetPassword 重置密码（忘记密码后通过密保问答重置）
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email      string `json:"email"`
		QuestionNo int    `json:"questionNo"`
		Answer     string `json:"answer"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	if req.Email == "" || req.Answer == "" {
		Error(w, http.StatusBadRequest, "邮箱和答案不能为空")
		return
	}
	resp, err := h.client.ResetPassword(req.Email, req.QuestionNo, req.Answer)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// SetAvatar 修改头像
func (h *AuthHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Avatar string `json:"avatar"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	if req.Avatar == "" {
		Error(w, http.StatusBadRequest, "头像数据不能为空")
		return
	}
	resp, err := h.client.SetAvatar(req.Avatar, "jpeg")
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// SetTitle 修改称号
func (h *AuthHandler) SetTitle(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		userID = r.URL.Query().Get("userId")
	}
	var req struct {
		Title string `json:"title"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}
	if req.Title == "" {
		Error(w, http.StatusBadRequest, "称号不能为空")
		return
	}
	resp, err := h.client.SetTitle(userID, req.Title)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// MyComments 我的评论列表
func (h *AuthHandler) MyComments(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	resp, err := h.client.GetUserComment(page)
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

		// 直接使用 token（保留原始格式，PicACG API 不需要 Bearer 前缀）
		client.SetToken(token)
		next(w, r)
	}
}
