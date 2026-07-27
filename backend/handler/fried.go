package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mogvl/bika/backend/fried"
)

type FriedHandler struct {
	client *fried.FriedClient
}

func NewFriedHandler(client *fried.FriedClient) *FriedHandler {
	return &FriedHandler{client: client}
}

// Posts 获取帖子列表
func (h *FriedHandler) Posts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	posts, total, err := h.client.GetPosts(page - 1)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, H{"posts": posts, "total": total})
}

// Comments 获取评论列表
func (h *FriedHandler) Comments(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if postID == "" {
		Error(w, http.StatusBadRequest, "帖子ID不能为空")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	comments, err := h.client.GetComments(postID, page-1)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, H{"comments": comments})
}

// SendComment 发送评论
func (h *FriedHandler) SendComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}

	postID := r.PathValue("id")
	if postID == "" {
		Error(w, http.StatusBadRequest, "帖子ID不能为空")
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	err := h.client.SendComment(postID, req.Content)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, H{"message": "发送成功"})
}

// LikeComment 点赞评论
func (h *FriedHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("id")
	if commentID == "" {
		Error(w, http.StatusBadRequest, "评论ID不能为空")
		return
	}

	err := h.client.LikeComment(commentID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, H{"message": "点赞成功"})
}
