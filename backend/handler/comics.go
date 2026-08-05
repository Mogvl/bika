package handler

import (
	"net/http"
	"strconv"

	"github.com/Mogvl/bika/backend/pica"
)

// ComicsHandler 漫画相关处理器
type ComicsHandler struct {
	client *pica.Client
}

// NewComicsHandler 创建漫画处理器
func NewComicsHandler(client *pica.Client) *ComicsHandler {
	return &ComicsHandler{client: client}
}

// Categories 获取分类列表
func (h *ComicsHandler) Categories(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetCategories()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// ListByCategory 按分类获取漫画列表
func (h *ComicsHandler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	category := r.URL.Query().Get("c")
	sort := r.URL.Query().Get("s")
	if sort == "" {
		sort = "ua"
	}

	if category == "" {
		Error(w, http.StatusBadRequest, "分类参数不能为空")
		return
	}

	resp, err := h.client.GetComicsByCategory(page, category, sort)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Search 高级搜索
func (h *ComicsHandler) Search(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	category := r.URL.Query().Get("c")
	sort := r.URL.Query().Get("s")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	var categories []string
	if category != "" {
		categories = []string{category}
	}

	resp, err := h.client.AdvancedSearch(keyword, categories, sort, page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Detail 获取漫画详情
func (h *ComicsHandler) Detail(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	resp, err := h.client.GetComicDetail(bookID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Eps 获取漫画章节列表
func (h *ComicsHandler) Eps(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	resp, err := h.client.GetComicEps(bookID, page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Pages 获取章节页面(图片)列表
func (h *ComicsHandler) Pages(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	epsID := r.PathValue("epsId")
	if bookID == "" || epsID == "" {
		Error(w, http.StatusBadRequest, "参数不完整")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	resp, err := h.client.GetComicPages(bookID, epsID, page)
	if err != nil {
		// PicACG API 偶尔会返回 500，给出友好提示
		Error(w, http.StatusInternalServerError, "获取页面失败，PicACG 服务可能暂时不可用")
		return
	}
	Success(w, resp.Data)
}

// Leaderboard 获取排行榜
func (h *ComicsHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	tt := r.URL.Query().Get("tt")
	if tt == "" {
		tt = "H24"
	}

	resp, err := h.client.GetLeaderboard(tt)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// KnightLeaderboard 获取骑士榜（用户排行）
func (h *ComicsHandler) KnightLeaderboard(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetKnightLeaderboard()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Random 获取随机漫画
func (h *ComicsHandler) Random(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetRandomComics()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Collections 获取推荐合集
func (h *ComicsHandler) Collections(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetCollections()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Keywords 获取搜索热词
func (h *ComicsHandler) Keywords(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetKeywords()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Favourites 获取收藏列表
func (h *ComicsHandler) Favourites(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	sort := r.URL.Query().Get("s")
	if sort == "" {
		sort = "da"
	}

	resp, err := h.client.GetFavourites(page, sort)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// AddFavourite 添加收藏
func (h *ComicsHandler) AddFavourite(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	resp, err := h.client.AddFavourite(bookID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Comments 获取评论列表
func (h *ComicsHandler) Comments(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	resp, err := h.client.GetComments(bookID, page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// SendComment 发送评论
func (h *ComicsHandler) SendComment(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}
	var req struct {
		Content string `json:"content"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Content == "" {
		Error(w, http.StatusBadRequest, "评论内容不能为空")
		return
	}
	resp, err := h.client.SendComment(bookID, req.Content)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// LikeComic 点赞漫画
func (h *ComicsHandler) LikeComic(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}
	resp, err := h.client.LikeComic(bookID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// LikeComment 点赞评论
func (h *ComicsHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("id")
	if commentID == "" {
		Error(w, http.StatusBadRequest, "评论ID不能为空")
		return
	}
	resp, err := h.client.LikeComment(commentID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// SubComments 获取子评论（楼中楼）
func (h *ComicsHandler) SubComments(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("id")
	if commentID == "" {
		Error(w, http.StatusBadRequest, "评论ID不能为空")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	resp, err := h.client.GetSubComments(commentID, page)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// SendSubComment 发送子评论
func (h *ComicsHandler) SendSubComment(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("id")
	if commentID == "" {
		Error(w, http.StatusBadRequest, "评论ID不能为空")
		return
	}
	var req struct {
		Content string `json:"content"`
	}
	if err := decodeJSON(r, &req); err != nil || req.Content == "" {
		Error(w, http.StatusBadRequest, "评论内容不能为空")
		return
	}
	resp, err := h.client.SendSubComment(commentID, req.Content)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// ReportComment 举报评论
func (h *ComicsHandler) ReportComment(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("id")
	if commentID == "" {
		Error(w, http.StatusBadRequest, "评论ID不能为空")
		return
	}
	resp, err := h.client.ReportComment(commentID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

// Recommendation 获取漫画相关推荐
func (h *ComicsHandler) Recommendation(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")
	if bookID == "" {
		Error(w, http.StatusBadRequest, "漫画ID不能为空")
		return
	}
	resp, err := h.client.GetComicRecommendation(bookID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, resp.Data)
}

