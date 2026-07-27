package pica

import (
	"fmt"
	"net/url"
)

// GetCategories 获取漫画分类列表
func (c *Client) GetCategories() (*APIResponse, error) {
	return c.doGet("categories")
}

// GetComicsByCategory 按分类获取漫画列表
// page: 页码, category: 分类名, sort: 排序方式(ua=更新时间, dd=新到旧, da=旧到新, ld=最多喜欢, vv=最多浏览)
func (c *Client) GetComicsByCategory(page int, category, sort string) (*APIResponse, error) {
	// 分类名需要 URL 编码（和原始 picacg-qt 一致）
	encodedCategory := url.QueryEscape(category)
	path := fmt.Sprintf("comics?page=%d&c=%s&s=%s", page, encodedCategory, sort)
	return c.doGet(path)
}

// GetComicsByCategoryPage 直接分页
func (c *Client) GetComicsByCategoryPage(page int, category string) (*APIResponse, error) {
	return c.GetComicsByCategory(page, category, "ua")
}

// AdvancedSearch 高级搜索
// keyword: 关键词, categories: 分类筛选, sort: 排序, page: 页码
func (c *Client) AdvancedSearch(keyword string, categories []string, sort string, page int) (*APIResponse, error) {
	body := map[string]any{
		"keyword":    keyword,
		"categories": categories,
		"sort":       sort,
		"page":       page,
	}
	return c.doPost(fmt.Sprintf("comics/advanced-search?page=%d", page), body)
}

// SearchByCategory 按分类搜索
func (c *Client) SearchByCategory(page int, category, sort string) (*APIResponse, error) {
	encodedCategory := url.QueryEscape(category)
	return c.doGet(fmt.Sprintf("comics?page=%d&c=%s&s=%s", page, encodedCategory, sort))
}

// GetComicDetail 获取漫画详情
func (c *Client) GetComicDetail(bookID string) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("comics/%s", bookID))
}

// GetComicEps 获取漫画章节列表
// bookID: 漫画ID, page: 页码
func (c *Client) GetComicEps(bookID string, page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("comics/%s/eps?page=%d", bookID, page))
}

// GetComicPages 获取章节页面(图片)列表
// bookID: 漫画ID, epsID: 章节ID, page: 页码
func (c *Client) GetComicPages(bookID, epsID string, page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("comics/%s/order/%s/pages?page=%d", bookID, epsID, page))
}

// GetLeaderboard 获取排行榜
// tt: 时间范围(H24=24小时, D7=7天, D30=30天)
func (c *Client) GetLeaderboard(tt string) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("comics/leaderboard?tt=%s&ct=VC", tt))
}

// GetRandomComics 获取随机漫画
func (c *Client) GetRandomComics() (*APIResponse, error) {
	return c.doGet("comics/random")
}

// GetCollections 获取推荐合集
func (c *Client) GetCollections() (*APIResponse, error) {
	return c.doGet("collections")
}

// GetKeywords 获取搜索热词
func (c *Client) GetKeywords() (*APIResponse, error) {
	return c.doGet("keywords")
}

// GetFavourites 获取收藏列表
func (c *Client) GetFavourites(page int, sort string) (*APIResponse, error) {
	if sort == "" {
		sort = "da"
	}
	return c.doGet(fmt.Sprintf("users/favourite?s=%s&page=%d", sort, page))
}

// AddFavourite 添加收藏
func (c *Client) AddFavourite(bookID string) (*APIResponse, error) {
	return c.doPost(fmt.Sprintf("comics/%s/favourite", bookID), nil)
}

// GetComicRecommendation 获取推荐漫画
func (c *Client) GetComicRecommendation(bookID string) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("comics/%s/recommendation", bookID))
}

// GetComments 获取漫画评论
func (c *Client) GetComments(bookID string, page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("comics/%s/comments?page=%d", bookID, page))
}

// SendComment 发送评论
func (c *Client) SendComment(bookID, content string) (*APIResponse, error) {
	return c.doPost(fmt.Sprintf("comics/%s/comments", bookID), map[string]string{"content": content})
}

// GetSubComments 获取子评论
func (c *Client) GetSubComments(commentID string, page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("comics/%s/comments?page=%d", commentID, page))
}

// SendSubComment 发送子评论
func (c *Client) SendSubComment(commentID, content string) (*APIResponse, error) {
	return c.doPost(fmt.Sprintf("comments/%s", commentID), map[string]string{"content": content})
}

// LikeComic 点赞漫画
func (c *Client) LikeComic(bookID string) (*APIResponse, error) {
	return c.doPost(fmt.Sprintf("comics/%s/like", bookID), nil)
}

// LikeComment 点赞评论
func (c *Client) LikeComment(commentID string) (*APIResponse, error) {
	return c.doPost(fmt.Sprintf("comments/%s/like", commentID), nil)
}

// ReportComment 举报评论
func (c *Client) ReportComment(commentID string) (*APIResponse, error) {
	return c.doPost(fmt.Sprintf("comments/%s/report", commentID), nil)
}

// GetCollections 获取推荐合集
func (c *Client) GetCollectionsDetail(collectionID string) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("collections/%s", collectionID))
}
