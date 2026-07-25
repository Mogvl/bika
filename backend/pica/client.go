package pica

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// picaClient PicACG API 客户端
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	headers    map[string]string
}

// NewClient 创建新的 API 客户端
func NewClient() *Client {
	return &Client{
		baseURL:    BaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SetToken 设置认证 token
func (c *Client) SetToken(token string) {
	c.token = token
}

// GetToken 获取当前 token
func (c *Client) GetToken() string {
	return c.token
}

// buildSignature 计算 API 签名
func (c *Client) buildSignature(path, method, timestamp, nonce string) string {
	// 签名数据: path + timestamp + nonce + method + apiKey
	data := strings.ToLower(path + timestamp + nonce + method + APIKey)

	// 签名密钥
	key := "~d}$Q7$eIni=V)9\\RK/P.RM4;9[7|@/CA}b~OW!3?EV`:<>M7pddUBL5n|0/*Cn"

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// buildHeaders 构建请求头
func (c *Client) buildHeaders(path, method string) map[string]string {
	now := fmt.Sprintf("%d", time.Now().Unix())
	nonce := strings.ReplaceAll(uuid.New().String(), "-", "")
	signature := c.buildSignature(path, method, now, nonce)

	headers := map[string]string{
		"api-key":           APIKey,
		"accept":            Accept,
		"app-channel":       AppChannel,
		"time":              now,
		"app-uuid":          UUID,
		"nonce":             nonce,
		"signature":         signature,
		"app-version":       APIVersion,
		"image-quality":     ImageQuality,
		"app-platform":      Platform,
		"app-build-version": BuildVersion,
		"user-agent":        Agent,
		"version":           AppVersion,
	}

	if method == "POST" || method == "PUT" {
		headers["Content-Type"] = "application/json; charset=UTF-8"
	}

	if c.token != "" {
		headers["authorization"] = c.token
	}

	return headers
}

// doRequest 执行 API 请求
func (c *Client) doRequest(method, path string, body io.Reader) (*APIResponse, error) {
	url := c.baseURL + strings.TrimPrefix(path, "/")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	headers := c.buildHeaders(path, method)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %s, body: %s", err, string(respBody))
	}

	if apiResp.Code != 200 {
		return &apiResp, fmt.Errorf("API 错误: code=%d, message=%s", apiResp.Code, apiResp.Message)
	}

	return &apiResp, nil
}

// doGet 执行 GET 请求
func (c *Client) doGet(path string) (*APIResponse, error) {
	return c.doRequest("GET", path, nil)
}

// doPost 执行 POST 请求
func (c *Client) doPost(path string, body any) (*APIResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}
	return c.doRequest("POST", path, strings.NewReader(string(jsonBody)))
}

// RawImage 获取原始图片数据
func (c *Client) RawImage(url string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	// 使用 Android 的 User-Agent
	req.Header.Set("User-Agent", Agent)
	req.Header.Set("Referer", "https://picaapi.picacomic.com/")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("图片请求失败: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("读取图片数据失败: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	return data, contentType, nil
}
