package pica

import "fmt"

// Login 用户登录，返回 token 和用户信息
func (c *Client) Login(email, password string) (*APIResponse, error) {
	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.doPost("auth/sign-in", req)
	if err != nil {
		return nil, fmt.Errorf("登录失败: %w", err)
	}

	// 提取 token
	if token, ok := resp.Data["token"]; ok {
		if t, ok := token.(string); ok {
			c.token = t
		}
	}

	return resp, nil
}

// GetProfile 获取用户信息
func (c *Client) GetProfile() (*APIResponse, error) {
	return c.doGet("users/profile")
}

// PunchIn 签到
func (c *Client) PunchIn() (*APIResponse, error) {
	return c.doPost("users/punch-in", nil)
}
