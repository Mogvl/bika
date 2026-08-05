package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	ChatBaseURL = "https://live-server.bidobido.xyz/"
	APIVersion  = "1.0.3"
)

// Room 聊天室
type Room struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IsAvailable bool   `json:"isAvailable"`
	MinLevel    int    `json:"minLevel"`
}

// Message 聊天消息
type Message struct {
	ID        string    `json:"_id"`
	RoomID    string    `json:"roomId"`
	Message   string    `json:"message"`
	Sender    User      `json:"sender"`
	CreatedAt string    `json:"createdAt"`
	ReplyTo   *Message  `json:"replyTo,omitempty"`
}

// User 聊天用户
type User struct {
	ID      string `json:"_id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Level   int    `json:"level"`
}

// ChatClient 聊天客户端
type ChatClient struct {
	token      string
	httpClient *http.Client
	mu         sync.RWMutex
}

// NewChatClient 创建聊天客户端
func NewChatClient() *ChatClient {
	return &ChatClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// Login 聊天服务器登录
func (c *ChatClient) Login(email, password string) (string, error) {
	url := ChatBaseURL + "auth/signin"
	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	c.setChatHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("聊天服务器连接失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Token string `json:"token"`
		Error string `json:"error"`
	}
	json.Unmarshal(respBody, &result)

	if result.Error != "" {
		return "", fmt.Errorf("聊天登录失败: %s", result.Error)
	}

	if result.Token == "" {
		return "", fmt.Errorf("聊天登录失败: 未获取到token")
	}

	c.mu.Lock()
	c.token = result.Token
	c.mu.Unlock()

	return result.Token, nil
}

// SetToken 设置聊天 token
func (c *ChatClient) SetToken(token string) {
	c.mu.Lock()
	c.token = token
	c.mu.Unlock()
}

// GetRooms 获取聊天室列表
func (c *ChatClient) GetRooms() ([]Room, error) {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	if token == "" {
		return nil, fmt.Errorf("未登录聊天服务器")
	}

	url := ChatBaseURL + "room/list"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setChatHeaders(req)
	req.Header.Set("authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Rooms []Room `json:"rooms"`
		Error string `json:"error"`
	}
	json.Unmarshal(body, &result)

	if result.Error != "" {
		return nil, fmt.Errorf("获取房间失败: %s", result.Error)
	}

	return result.Rooms, nil
}

// GetMessages 获取聊天消息
func (c *ChatClient) GetMessages(roomID string, page int) ([]Message, error) {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	url := fmt.Sprintf("%smessage/list?roomId=%s&page=%d", ChatBaseURL, roomID, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setChatHeaders(req)
	req.Header.Set("authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Messages []Message `json:"messages"`
		Error    string    `json:"error"`
	}
	json.Unmarshal(body, &result)

	if result.Error != "" {
		return nil, fmt.Errorf("获取消息失败: %s", result.Error)
	}

	return result.Messages, nil
}

// SendMessage 发送消息
// referenceId: 回复引用的消息ID（可选）
func (c *ChatClient) SendMessage(roomID, message, referenceID string) error {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	url := ChatBaseURL + "message/send-message"
	if referenceID == "" {
		referenceID = generateUUID()
	}
	body := fmt.Sprintf(`{"roomId":"%s","message":"%s","referenceId":"%s","userMentions":[]}`,
		roomID, message, referenceID)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}

	c.setChatHeaders(req)
	req.Header.Set("authorization", "Bearer "+token)

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

// GetProfile 获取聊天用户信息
func (c *ChatClient) GetProfile() (map[string]any, error) {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	url := ChatBaseURL + "user/profile"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.setChatHeaders(req)
	req.Header.Set("authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]any
	json.Unmarshal(body, &result)

	return result, nil
}

func (c *ChatClient) setChatHeaders(req *http.Request) {
	req.Header.Set("user-agent", "Dart/2.19 (dart:io)")
	req.Header.Set("accept-encoding", "gzip")
	req.Header.Set("api-version", APIVersion)
	req.Header.Set("content-type", "application/json; charset=UTF-8")
}

func generateUUID() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		time.Now().UnixNano()&0xFFFFFFFF,
		time.Now().UnixNano()>>16&0xFFFF,
		time.Now().UnixNano()>>32&0xFFFF|0x4000,
		time.Now().UnixNano()>>48&0xFFFF|0x8000,
		time.Now().UnixNano()&0xFFFFFFFFFFFF)
}
