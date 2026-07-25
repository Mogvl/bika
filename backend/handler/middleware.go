package handler

import (
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
// 前端传递图片的 fileServer 和 path，后端代理获取图片数据
func (h *ImageHandler) Proxy(w http.ResponseWriter, r *http.Request) {
	imageURL := r.URL.Query().Get("url")
	if imageURL == "" {
		Error(w, http.StatusBadRequest, "缺少图片 URL 参数")
		return
	}

	// URL 解码
	if decoded, err := url.QueryUnescape(imageURL); err == nil {
		imageURL = decoded
	}

	// 构建完整图片 URL（如果是相对路径）
	fileServer := r.URL.Query().Get("fileServer")
	path := r.URL.Query().Get("path")
	var fullURL string
	if imageURL != "" {
		fullURL = imageURL
	} else if fileServer != "" && path != "" {
		fullURL = pica.GetImageURL(fileServer, path)
	} else {
		Error(w, http.StatusBadRequest, "请提供 url 或 fileServer+path 参数")
		return
	}

	// 下载图片
	data, contentType, err := h.client.RawImage(fullURL)
	if err != nil {
		// 如果图片下载失败，返回占位图
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

// ImageURL 构建图片 URL（返回前端可直接使用的 URL）
// 这个接口返回代理图片的 URL 列表
func (h *ImageHandler) ImageURL(w http.ResponseWriter, r *http.Request) {
	fileServer := r.URL.Query().Get("fileServer")
	path := r.URL.Query().Get("path")

	if fileServer == "" || path == "" {
		Error(w, http.StatusBadRequest, "缺少 fileServer 或 path 参数")
		return
	}

	proxyURL := pica.GetImageURL(fileServer, path)

	// 构建代理 URL
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

// ExtractComicImages 从漫画章节数据中提取所有图片的代理 URL
func (h *ImageHandler) ExtractComicImages(pageData map[string]any, scheme, host string) []string {
	pages, ok := pageData["pages"].([]any)
	if !ok {
		return nil
	}

	var urls []string
	for _, p := range pages {
		page, ok := p.(map[string]any)
		if !ok {
			continue
		}
		media, ok := page["media"].(map[string]any)
		if !ok {
			continue
		}
		fileServer, _ := media["fileServer"].(string)
		path, _ := media["path"].(string)
		if fileServer != "" && path != "" {
			proxyURL := scheme + "://" + host + "/api/image/proxy?url=" +
				url.QueryEscape(pica.GetImageURL(fileServer, path))
			urls = append(urls, proxyURL)
		}
	}
	return urls
}

// CORS 中间件 - 处理跨域请求
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

// StaticFileMiddleware 静态文件中间件 - 处理 SPA 路由
func StaticFileMiddleware(staticDir http.FileSystem, indexPage string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 如果是 API 请求，跳过静态文件处理
			if strings.HasPrefix(r.URL.Path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}

			// 尝试提供静态文件
			path := r.URL.Path
			if path == "/" {
				path = "/" + indexPage
			}

			f, err := staticDir.Open(path)
			if err != nil {
				// 文件不存在，返回 SPA 入口页面（支持客户端路由）
				f, err = staticDir.Open("/" + indexPage)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				defer f.Close()

				stat, _ := f.Stat()
				if stat != nil {
					http.ServeContent(w, r, indexPage, stat.ModTime(), f)
				}
				return
			}
			defer f.Close()

			stat, _ := f.Stat()
			if stat != nil {
				// 根据扩展名设置 Content-Type
				if strings.HasSuffix(path, ".js") {
					w.Header().Set("Content-Type", "application/javascript")
				} else if strings.HasSuffix(path, ".css") {
					w.Header().Set("Content-Type", "text/css")
				}
				http.ServeContent(w, r, path, stat.ModTime(), f)
			}
		})
	}
}
