package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Mogvl/bika/backend/download"
)

// LocalHandler 本地库处理器
type LocalHandler struct {
	manager *download.Manager
}

func NewLocalHandler(manager *download.Manager) *LocalHandler {
	return &LocalHandler{manager: manager}
}

// LocalComic 本地漫画
type LocalComic struct {
	Title    string   `json:"title"`
	Path     string   `json:"path"`
	CoverURL string   `json:"coverUrl"`
	Eps      []LocalEps `json:"eps"`
	Cover    string   `json:"cover"`
}

// LocalEps 本地章节
// Path 相对下载根目录；Pages 相对下载根目录（均用 posix 分隔符）
type LocalEps struct {
	Title string   `json:"title"`
	Path  string   `json:"path"`
	Pages []string `json:"pages"`
}

// List 获取本地漫画库列表
func (h *LocalHandler) List(w http.ResponseWriter, r *http.Request) {
	root := h.manager.GetDownloadDir()
	comics, err := scanLocalLibrary(root)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, H{"comics": comics})
}

// Eps 获取本地漫画章节（含图片）
func (h *LocalHandler) Eps(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		Error(w, http.StatusBadRequest, "path 不能为空")
		return
	}
	root := h.manager.GetDownloadDir()

	// path 为相对下载根目录的漫画目录（posix 格式），拼回绝对路径
	absRoot, _ := filepath.Abs(root)
	absPath, _ := filepath.Abs(filepath.Join(absRoot, filepath.FromSlash(path)))
	if !strings.HasPrefix(absPath, absRoot) {
		Error(w, http.StatusForbidden, "非法路径")
		return
	}

	eps, err := scanLocalEps(absRoot, absPath)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, H{"eps": eps})
}

// Image 提供本地图片文件（相对下载目录）
func (h *LocalHandler) Image(w http.ResponseWriter, r *http.Request) {
	rel := r.URL.Query().Get("path")
	if rel == "" {
		Error(w, http.StatusBadRequest, "path 不能为空")
		return
	}
	root := h.manager.GetDownloadDir()
	absRoot, _ := filepath.Abs(root)
	absPath, _ := filepath.Abs(filepath.Join(absRoot, filepath.FromSlash(rel)))
	if !strings.HasPrefix(absPath, absRoot) {
		Error(w, http.StatusForbidden, "非法路径")
		return
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		Error(w, http.StatusNotFound, "文件不存在")
		return
	}
	ext := strings.ToLower(filepath.Ext(absPath))
	contentType := "image/jpeg"
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 扫描本地库：一级目录 = 漫画，二级目录 = 章节，图片文件 = 页
func scanLocalLibrary(root string) ([]LocalComic, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var comics []LocalComic
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		comicDir := filepath.Join(root, e.Name())
		eps, _ := scanLocalEps(root, comicDir)
		if len(eps) == 0 {
			continue
		}
		// 相对下载根目录的漫画目录路径（posix 格式，供 /api/local/image 使用）
		comicRel, _ := filepath.Rel(root, comicDir)
		comic := LocalComic{
			Title: e.Name(),
			Path:  filepath.ToSlash(comicRel),
			Eps:   eps,
		}
		// 找封面：第一个章节的第一张图（相对下载根目录）
		if len(eps) > 0 && len(eps[0].Pages) > 0 {
			comic.Cover = eps[0].Pages[0]
		}
		comics = append(comics, comic)
	}

	sort.Slice(comics, func(i, j int) bool {
		return comics[i].Title < comics[j].Title
	})
	return comics, nil
}

// 扫描章节：返回相对下载根目录的章节路径与图片路径（均用于安全访问）
func scanLocalEps(root, comicDir string) ([]LocalEps, error) {
	entries, err := os.ReadDir(comicDir)
	if err != nil {
		return nil, err
	}

	var eps []LocalEps
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		epsDir := filepath.Join(comicDir, e.Name())
		images, _ := filepath.Glob(filepath.Join(epsDir, "*"))
		var pages []string
		imageExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		for _, img := range images {
			ext := strings.ToLower(filepath.Ext(img))
			if imageExts[ext] {
				// 相对下载根目录路径（用于安全访问）
				rel, err := filepath.Rel(root, img)
				if err == nil {
					pages = append(pages, filepath.ToSlash(rel))
				}
			}
		}
		if len(pages) == 0 {
			continue
		}
		epsRel, _ := filepath.Rel(root, epsDir)
		sort.Strings(pages)
		eps = append(eps, LocalEps{
			Title: e.Name(),
			Path:  filepath.ToSlash(epsRel),
			Pages: pages,
		})
	}
	return eps, nil
}
