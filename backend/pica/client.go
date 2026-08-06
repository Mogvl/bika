package pica

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Client PicACG API 客户端
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient 创建新的 API 客户端
func NewClient() *Client {
	c := &Client{
		baseURL:    BaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	return c
}

// Init 初始化客户端 - 获取 API 真实 IP
func (c *Client) Init() error {
	// 创建一个独立的 HTTP 客户端用于 IP 探测
	probeClient := &http.Client{Timeout: 10 * time.Second}

	resp, err := probeClient.Get("http://68.183.234.72/init")
	if err != nil {
		return fmt.Errorf("获取 API 服务器 IP 失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var initResp struct {
		Status    string   `json:"status"`
		Addresses []string `json:"addresses"`
	}
	if err := json.Unmarshal(body, &initResp); err != nil {
		return fmt.Errorf("解析 init 响应失败: %w", err)
	}

	if initResp.Status != "ok" || len(initResp.Addresses) == 0 {
		return fmt.Errorf("init 返回无效数据: %s", string(body))
	}

	apiIP := initResp.Addresses[0]
	fmt.Printf("获取到 PicACG API 真实 IP: %s\n", apiIP)

	// 创建自定义 Transport，将 API 域名解析到真实 IP
	transport := &http.Transport{
		ForceAttemptHTTP2: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if strings.Contains(addr, "picaapi.picacomic.com") {
				addr = apiIP + ":443"
			}
			var dialer net.Dialer
			return dialer.DialContext(ctx, network, addr)
		},
	}

	c.httpClient = &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return nil
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
	// 签名数据: path + timestamp + nonce + method + apiKey (和原始 picacg-qt 一致)
	data := strings.ToLower(path + timestamp + nonce + method + APIKey)
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

	// 设置 Host 头为 API 域名（因使用真实 IP 连接，需要保留 Host）
	req.Host = "picaapi.picacomic.com"

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
		// 打日志方便诊断上游返回的原始信息
		log.Printf("[PicACG] %s %s → code=%d message=%q raw=%s", method, path, apiResp.Code, apiResp.Message, string(respBody))
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
	if body == nil {
		// 上游要求 JSON 必须是对象或数组，nil 会序列化成 null 导致 400
		body = map[string]any{}
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}
	return c.doRequest("POST", path, strings.NewReader(string(jsonBody)))
}

// doPut 执行 PUT 请求
func (c *Client) doPut(path string, body any) (*APIResponse, error) {
	if body == nil {
		body = map[string]any{}
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}
	return c.doRequest("PUT", path, strings.NewReader(string(jsonBody)))
}

// RawImage 获取原始图片数据
func (c *Client) RawImage(url string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

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
