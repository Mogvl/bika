package pica

// GetImageURL 构建完整的图片 URL
// 图片服务器可能返回相对路径，需要拼接
func GetImageURL(fileServer, path string) string {
	if fileServer == "" || path == "" {
		return ""
	}

	// 如果 path 已经是完整 URL，直接返回
	if len(path) > 4 && path[:4] == "http" {
		return path
	}

	// 拼接 fileServer 和 path
	// fileServer 可能以 / 结尾，path 可能以 / 开头
	if fileServer[len(fileServer)-1] == '/' {
		fileServer = fileServer[:len(fileServer)-1]
	}
	if path[0] == '/' {
		path = path[1:]
	}
	return fileServer + "/static/" + path
}
