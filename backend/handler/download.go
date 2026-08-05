package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mogvl/bika/backend/download"
)

type DownloadHandler struct {
	manager *download.Manager
}

func NewDownloadHandler(manager *download.Manager) *DownloadHandler {
	return &DownloadHandler{manager: manager}
}

// List 获取下载列表
func (h *DownloadHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks := h.manager.GetTasks()
	Success(w, H{"tasks": tasks})
}

// Add 添加下载任务
func (h *DownloadHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "仅支持 POST")
		return
	}

	var req struct {
		BookID   string `json:"bookId"`
		Title    string `json:"title"`
		CoverURL string `json:"coverUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	if req.BookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	task := h.manager.AddTask(req.BookID, req.Title, req.CoverURL)
	Success(w, H{"task": task})
}

// Status 获取下载状态
func (h *DownloadHandler) Status(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	task := h.manager.GetTask(bookID)
	if task == nil {
		Error(w, http.StatusNotFound, "任务不存在")
		return
	}

	Success(w, H{"task": task})
}

// Cancel 暂停下载
func (h *DownloadHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	h.manager.CancelTask(bookID)
	Success(w, H{"message": "已暂停"})
}

// Resume 恢复下载
func (h *DownloadHandler) Resume(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	h.manager.ResumeTask(bookID)
	Success(w, H{"message": "已恢复"})
}

// Remove 删除下载任务
func (h *DownloadHandler) Remove(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	deleteFile := r.URL.Query().Get("deleteFile") == "true"

	h.manager.RemoveTask(bookID, deleteFile)
	Success(w, H{"message": "已删除"})
}
