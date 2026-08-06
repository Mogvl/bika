package pica

import (
	"errors"
	"fmt"
)

// GetCodeErrMsg 根据错误码返回错误信息
func GetCodeErrMsg(code string) string {
	msgMap := map[string]string{
		"1004": "账号或密码错误",
		"1005": "未授权",
		"1006": "账号未激活，请先到邮箱激活",
		"1007": "找不到该账号",
		"1008": "该邮箱已被注册",
		"1009": "该昵称已被使用",
		"1010": "账号或密码错误",
		"1014": "该漫画正在审核中",
		"1019": "等级不够，无法操作",
		"1023": "请求过于频繁",
		"1024": "不支持该邮箱",
		"1025": "哔咔是注册商标，不能使用",
		"1026": "邮箱格式错误",
		"1029": "时间不同步，请调整设备时间",
	}
	if msg, ok := msgMap[code]; ok {
		return msg
	}
	return ""
}

// Login 用户登录，返回 token 和用户信息
func (c *Client) Login(email, password string) (*APIResponse, error) {
	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.doPost("auth/sign-in", req)
	if err != nil {
		// 如果是 API 错误（如 401），提取错误信息
		if resp != nil {
			errMsg := "登录失败，请检查账号密码"
			if resp.Message != "" && resp.Message != "success" {
				errMsg = resp.Message
			}
			if code, ok := resp.Data["code"].(string); ok {
				errMsg = GetCodeErrMsg(code)
				if errMsg == "" {
					errMsg = fmt.Sprintf("登录失败 (错误码: %s)", code)
				}
			}
			return nil, errors.New(errMsg)
		}
		return nil, fmt.Errorf("登录失败: %w", err)
	}

	// 提取 token
	if resp.Data != nil {
		if token, ok := resp.Data["token"]; ok {
			if t, ok := token.(string); ok && t != "" {
				c.token = t
			}
		}
	}

	// 检查是否获取到 token
	if c.token == "" {
		return nil, fmt.Errorf("登录失败，请使用邮箱登录")
	}

	return resp, nil
}

// Register 用户注册
func (c *Client) Register(email, password, name, birthday, gender string, questions []int, answers []string) (*APIResponse, error) {
	req := map[string]any{
		"email":     email,
		"password":  password,
		"name":      name,
		"birthday":  birthday,
		"gender":    gender,
		"question1": questions[0],
		"question2": questions[1],
		"question3": questions[2],
		"answer1":   answers[0],
		"answer2":   answers[1],
		"answer3":   answers[2],
	}
	return c.doPost("auth/register", req)
}

// GetProfile 获取用户信息
func (c *Client) GetProfile() (*APIResponse, error) {
	return c.doGet("users/profile")
}

// PunchIn 签到
func (c *Client) PunchIn() (*APIResponse, error) {
	return c.doPost("users/punch-in", nil)
}

// ChangePassword 修改密码
func (c *Client) ChangePassword(token, oldPassword, newPassword string) (*APIResponse, error) {
	req := map[string]string{
		"old_password": oldPassword,
		"new_password": newPassword,
	}
	return c.doPut("users/password", req)
}

// ForgotPassword 忘记密码
func (c *Client) ForgotPassword(email string) (*APIResponse, error) {
	return c.doPost("auth/forgot-password", map[string]string{"email": email})
}

// ResetPassword 重置密码
func (c *Client) ResetPassword(email string, questionNo int, answer string) (*APIResponse, error) {
	req := map[string]any{
		"email":      email,
		"questionNo": questionNo,
		"answer":     answer,
	}
	return c.doPost("auth/reset-password", req)
}

// SetAvatar 设置头像
func (c *Client) SetAvatar(imageData, picFormat string) (*APIResponse, error) {
	imgData := fmt.Sprintf("data:image/%s;base64,%s", picFormat, imageData)
	return c.doPut("users/avatar", map[string]string{"avatar": imgData})
}

// SetTitle 设置称号
func (c *Client) SetTitle(userID, title string) (*APIResponse, error) {
	return c.doPut(fmt.Sprintf("users/%s/title", userID), map[string]string{"title": title})
}

// GetUserComment 获取我的评论
func (c *Client) GetUserComment(page int) (*APIResponse, error) {
	return c.doGet(fmt.Sprintf("users/my-comments?page=%d", page))
}

// GetMyFavourite 获取收藏
func (c *Client) GetMyFavourite(page int, sort string) (*APIResponse, error) {
	if sort == "" {
		sort = "da"
	}
	return c.doGet(fmt.Sprintf("users/favourite?s=%s&page=%d", sort, page))
}
