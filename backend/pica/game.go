package pica

import "fmt"

// GetGames 获取游戏列表
func (c *Client) GetGames(page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("games?page=%d", page))
}

// GetGameInfo 获取游戏详情
func (c *Client) GetGameInfo(gameID string) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("games/%s", gameID))
}

// GetGameEps 获取游戏章节
func (c *Client) GetGameEps(gameID string, page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("games/%s/eps?page=%d", gameID, page))
}

// GetGamePages 获取游戏章节页面
func (c *Client) GetGamePages(gameID, epsID string, page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("games/%s/order/%s/pages?page=%d", gameID, epsID, page))
}

// GetGameComments 获取游戏评论
func (c *Client) GetGameComments(gameID string, page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("games/%s/comments?page=%d", gameID, page))
}

// SendGameComment 发送游戏评论
func (c *Client) SendGameComment(gameID, content string) (*APIResponse, error) {
	return c.doPost(fmt.Sprintf("games/%s/comments", gameID), map[string]string{"content": content})
}
