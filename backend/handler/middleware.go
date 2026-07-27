package handler

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/Mogvl/bika/backend/pica"
)

// ImageHandler 图片代理处理器
type ImageHandler struct {
	client *pica.Client
}

// NewImageHandler 创建图片处理器
func NewImageHandler(client *pica.Client) *ImageHandler {
	return &ImageHandler{client: client}
}

// Proxy 代理图片请求
func (h *ImageHandler) Proxy(w http.ResponseWriter, r *http.Request) {
	fileServer := r.URL.Query().Get("fileServer")
	path := r.URL.Query().Get("path")
	imageURL := r.URL.Query().Get("url")

	var fullURL string
	if fileServer != "" && path != "" {
		fullURL = pica.GetImageURL(fileServer, path)
	} else if imageURL != "" {
		fullURL = imageURL
	} else {
		Error(w, http.StatusBadRequest, "缺少图片参数")
		return
	}

	data, contentType, err := h.client.RawImage(fullURL)
	if err != nil {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="600" viewBox="0 0 400 600">
			<rect width="400" height="600" fill="#f0f0f0"/>
			<text x="200" y="300" font-family="Arial" font-size="16" fill="#999" text-anchor="middle">图片加载失败</text>
		</svg>`))
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// ImageURL 构建图片 URL
func (h *ImageHandler) ImageURL(w http.ResponseWriter, r *http.Request) {
	fileServer := r.URL.Query().Get("fileServer")
	path := r.URL.Query().Get("path")

	if fileServer == "" || path == "" {
		Error(w, http.StatusBadRequest, "缺少 fileServer 或 path 参数")
		return
	}

	proxyURL := pica.GetImageURL(fileServer, path)

	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	host := r.Host
	proxyPath := scheme + "://" + host + "/api/image/proxy?url=" + url.QueryEscape(proxyURL)

	Success(w, H{
		"url":       proxyURL,
		"proxy_url": proxyPath,
	})
}

// CORS 中间件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SPAFileServer 创建支持 SPA 路由的静态文件服务器
func SPAFileServer(staticFS fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API 请求不处理
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// 尝试打开请求的文件
		filePath := strings.TrimPrefix(r.URL.Path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		// 检查文件是否存在并读取
		if data, err := fs.ReadFile(staticFS, filePath); err == nil {
			// 设置正确的 Content-Type
			contentType := http.DetectContentType(data)
			if strings.HasSuffix(filePath, ".js") {
				contentType = "application/javascript"
			} else if strings.HasSuffix(filePath, ".css") {
				contentType = "text/css"
			} else if strings.HasSuffix(filePath, ".html") {
				contentType = "text/html; charset=utf-8"
			} else if strings.HasSuffix(filePath, ".json") {
				contentType = "application/json"
			} else if strings.HasSuffix(filePath, ".svg") {
				contentType = "image/svg+xml"
			} else if strings.HasSuffix(filePath, ".png") {
				contentType = "image/png"
			} else if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
				contentType = "image/jpeg"
			} else if strings.HasSuffix(filePath, ".gif") {
				contentType = "image/gif"
			} else if strings.HasSuffix(filePath, ".woff2") {
				contentType = "font/woff2"
			} else if strings.HasSuffix(filePath, ".woff") {
				contentType = "font/woff"
			}
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}

		// 文件不存在，提供 index.html（SPA 路由）
		if data, err := fs.ReadFile(staticFS, "index.html"); err == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}

		http.NotFound(w, r)
	})
}
