package fried

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	FriedBaseURL = "https://post-api.wikawika.xyz"
)

// Post 帖子
type Post struct {
	ID        string `json:"_id"`
	Content   string `json:"content"`
	User      User   `json:"_user"`
	Medias    []Media `json:"medias"`
	CreatedAt string `json:"createdAt"`
	TotalComments int `json:"totalComments"`
	TotalLikes    int `json:"totalLikes"`
	Liked         bool `json:"liked"`
}

// User 用户
type User struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	Title     string `json:"title"`
	Character string `json:"character"`
	Avatar    string `json:"avatar"`
}

// Media 媒体
type Media struct {
	ID   string `json:"_id"`
	Path string `json:"path"`
}

// Comment 评论
type Comment struct {
	ID        string `json:"_id"`
	Content   string `json:"content"`
	User      User   `json:"_user"`
	CreatedAt string `json:"createdAt"`
	TotalLikes int  `json:"totalLikes"`
	Liked      bool `json:"liked"`
}

// FriedClient 锅贴客户端
type FriedClient struct {
	token      string
	httpClient *http.Client
}

// NewFriedClient 创建锅贴客户端
func NewFriedClient() *FriedClient {
	return &FriedClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// SetToken 设置 token
func (c *FriedClient) SetToken(token string) {
	c.token = token
}

// GetPosts 获取帖子列表
func (c *FriedClient) GetPosts(page int) ([]Post, int, error) {
	url := fmt.Sprintf("%s/posts?offset=%d", FriedBaseURL, page*10)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Data struct {
			Posts []Post `json:"posts"`
			Total int    `json:"total"`
		} `json:"data"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	json.Unmarshal(body, &result)

	if result.Error.Message != "" {
		return nil, 0, fmt.Errorf("%s", result.Error.Message)
	}

	return result.Data.Posts, result.Data.Total, nil
}

// GetComments 获取评论列表
func (c *FriedClient) GetComments(postID string, page int) ([]Comment, error) {
	url := fmt.Sprintf("%s/posts/%s/comments?offset=%d", FriedBaseURL, postID, page*10)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Data []Comment `json:"data"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	json.Unmarshal(body, &result)

	if result.Error.Message != "" {
		return nil, fmt.Errorf("%s", result.Error.Message)
	}

	return result.Data, nil
}

// SendComment 发送评论
func (c *FriedClient) SendComment(postID, content string) error {
	url := fmt.Sprintf("%s/comments", FriedBaseURL)
	body, err := json.Marshal(map[string]string{
		"content": content,
		"postId":  postID,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("发送失败: %s", string(respBody))
	}

	return nil
}

// LikeComment 点赞评论
func (c *FriedClient) LikeComment(commentID string) error {
	url := fmt.Sprintf("%s/comments/%s/like", FriedBaseURL, commentID)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("点赞失败")
	}

	return nil
}

func (c *FriedClient) setHeaders(req *http.Request) {
	req.Header.Set("Referer", req.URL.String())
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	req.Header.Set("token", c.token)
	req.Header.Set("Content-Type", "application/json")
}
