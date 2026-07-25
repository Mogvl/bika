package handler

import (
	"net/http"
	"strconv"

	"github.com/Mogvl/bika/backend/pica"
)

type GameHandler struct {
	client *pica.Client
}

func NewGameHandler(client *pica.Client) *GameHandler {
	return &GameHandler{client: client}
}

// List 游戏列表
func (h *GameHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	resp, err := h.client.GetGames(page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Detail 游戏详情
func (h *GameHandler) Detail(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	if gameID == "" {
		Error(w, http.StatusBadRequest, "游戏ID不能为空")
		return
	}
	resp, err := h.client.GetGameInfo(gameID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Eps 游戏章节
func (h *GameHandler) Eps(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	resp, err := h.client.GetGameEps(gameID, page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Pages 游戏章节页面
func (h *GameHandler) Pages(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	epsID := r.PathValue("epsId")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	resp, err := h.client.GetGamePages(gameID, epsID, page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Comments 游戏评论
func (h *GameHandler) Comments(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	resp, err := h.client.GetGameComments(gameID, page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}
