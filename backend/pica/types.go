package pica

// API 配置常量
const (
	BaseURL      = "https://picaapi.picacomic.com/"
	APIKey       = "C69BAF41DA5ABD1FFEDC6D2FEA56B"
	APIVersion   = "2.2.1.3.3.4"
	BuildVersion = "45"
	AppChannel   = "3"
	Platform     = "android"
	Accept       = "application/vnd.picacomic.com.v1+json"
	Agent        = "okhttp/3.8.1"
	ImageQuality = "original"
	UUID         = "defaultUuid"
	AppVersion   = "v1.5.4"
)

// ======================== API 响应结构 ========================

// APIResponse 通用 API 响应包装
type APIResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// ======================== 分类 ========================

// Category 漫画分类
type Category struct {
	Title  string `json:"title"`
	Thumb  string `json:"thumb"`
	Active bool   `json:"isActive"`
	Link   string `json:"link"`
}

// CategoriesResponse 分类列表响应
type CategoriesResponse struct {
	Categories []Category `json:"categories"`
}

// ======================== 漫画 ========================

// ComicPage 漫画分页
type ComicPage struct {
	Comics      []Comic `json:"comics"`
	Page        int     `json:"page"`
	Pages       int     `json:"pages"`
	Total       int     `json:"total"`
	Limit       int     `json:"limit"`
}

// Comic 漫画基本信息
type Comic struct {
	ID           string   `json:"_id"`
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	Categories   []string `json:"categories"`
	Tags         []string `json:"tags"`
	Description  string   `json:"description"`
	Thumb        Thumb    `json:"thumb"`
	EPSCount     int      `json:"epsCount"`
	PagesCount   int      `json:"pagesCount"`
	Finished     bool     `json:"finished"`
	LikesCount   int      `json:"likesCount"`
	TotalViews   int      `json:"totalViews"`
	TotalLikes   int      `json:"totalLikes"`
	CommentsCount int     `json:"commentsCount"`
	IsFavourite  bool     `json:"isFavourite"`
	IsLiked      bool     `json:"isLiked"`
	UpdatedAt    string   `json:"updated_at"`
	CreatedAt    string   `json:"created_at"`
}

// Thumb 缩略图
type Thumb struct {
	OriginalName string `json:"originalName"`
	Path         string `json:"path"`
	FileServer   string `json:"fileServer"`
}

// ======================== 章节 ========================

// ComicDetail 漫画详情
type ComicDetail struct {
	Comic
	// 额外的详情字段
}

// ComicDetailResponse 漫画详情响应
type ComicDetailResponse struct {
	Comic Comic `json:"comic"`
}

// EP 章节(话)
type EP struct {
	ID         string `json:"_id"`
	Title      string `json:"title"`
	Order      int    `json:"order"`
	UpdatedAt  string `json:"updated_at"`
	PagesCount int    `json:"pagesCount"`
}

// EpsResponse 章节列表响应
type EpsResponse struct {
	Eps  []EP `json:"eps"`
	Page int  `json:"page"`
	Pages int `json:"pages"`
	Total int `json:"total"`
}

// ======================== 页面(图片) ========================

// Page 漫画页面
type Page struct {
	ID           int    `json:"id"`
	Media        Media  `json:"media"`
}

// Media 媒体信息
type Media struct {
	OriginalName string `json:"originalName"`
	Path         string `json:"path"`
	FileServer   string `json:"fileServer"`
}

// PagesResponse 页面列表响应
type PagesResponse struct {
	Pages []Page `json:"pages"`
	Page  int    `json:"page"`
	Total int    `json:"total"`
	EP    EP     `json:"ep"`
}

// ======================== 排行榜 ========================

// LeaderboardResponse 排行榜响应
type LeaderboardResponse struct {
	Comics []Comic `json:"comics"`
}

// ======================== 搜索 ========================

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword    string   `json:"keyword"`
	Categories []string `json:"categories"`
	Sort       string   `json:"sort"`
	Page       int      `json:"page"`
}

// KeywordsResponse 热词响应
type KeywordsResponse struct {
	Keywords []string `json:"keywords"`
}

// ======================== 收藏 ========================

// FavouritesResponse 收藏列表响应
type FavouritesResponse struct {
	Comics []Comic `json:"comics"`
	Page   int     `json:"page"`
	Pages  int     `json:"pages"`
	Total  int     `json:"total"`
}
