package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Mogvl/bika/backend/pica"
)

// DownloadTask 下载任务
type DownloadTask struct {
	ID          string    `json:"id"`
	BookID      string    `json:"bookId"`
	Title       string    `json:"title"`
	CoverURL    string    `json:"coverUrl"`
	TotalPages  int       `json:"totalPages"`
	Downloaded  int       `json:"downloaded"`
	Status      string    `json:"status"` // waiting, downloading, completed, error, paused
	Error       string    `json:"error,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Manager 下载管理器
type Manager struct {
	client     *pica.Client
	downloadDir string
	tasks      map[string]*DownloadTask
	mu         sync.RWMutex
	httpClient *http.Client
}

// NewManager 创建下载管理器
func NewManager(client *pica.Client, downloadDir string) *Manager {
	return &Manager{
		client:      client,
		downloadDir: downloadDir,
		tasks:       make(map[string]*DownloadTask),
		httpClient:  &http.Client{Timeout: 60 * time.Second},
	}
}

// AddTask 添加下载任务
func (m *Manager) AddTask(bookID, title, coverURL string) *DownloadTask {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	if task, exists := m.tasks[bookID]; exists {
		if task.Status == "completed" || task.Status == "error" {
			task.Status = "waiting"
			task.Error = ""
			task.Downloaded = 0
			task.UpdatedAt = time.Now()
		}
		return task
	}

	task := &DownloadTask{
		ID:        bookID,
		BookID:    bookID,
		Title:     title,
		CoverURL:  coverURL,
		Status:    "waiting",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.tasks[bookID] = task

	// 异步开始下载
	go m.downloadComic(task)

	return task
}

// GetTasks 获取所有下载任务
func (m *Manager) GetTasks() []*DownloadTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]*DownloadTask, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetTask 获取单个下载任务
func (m *Manager) GetTask(bookID string) *DownloadTask {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tasks[bookID]
}

// CancelTask 取消下载任务
func (m *Manager) CancelTask(bookID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if task, exists := m.tasks[bookID]; exists {
		task.Status = "paused"
		task.UpdatedAt = time.Now()
	}
}

// RemoveTask 删除下载任务
func (m *Manager) RemoveTask(bookID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tasks, bookID)
}

// downloadComic 下载漫画
func (m *Manager) downloadComic(task *DownloadTask) {
	m.mu.Lock()
	task.Status = "downloading"
	task.UpdatedAt = time.Now()
	m.mu.Unlock()

	// 获取章节列表
	epsResp, err := m.client.GetComicEps(task.BookID, 1)
	if err != nil {
		m.setError(task, fmt.Sprintf("获取章节失败: %v", err))
		return
	}

	epsData, _ := epsResp.Data["eps"].(map[string]any)
	if epsData == nil {
		m.setError(task, "章节数据格式错误")
		return
	}

	docs, _ := epsData["docs"].([]any)
	if len(docs) == 0 {
		m.setError(task, "没有章节")
		return
	}

	// 创建漫画目录
	comicDir := filepath.Join(m.downloadDir, sanitizeFilename(task.Title))
	os.MkdirAll(comicDir, 0755)

	// 下载每个章节
	for i, ep := range docs {
		epMap, ok := ep.(map[string]any)
		if !ok {
			continue
		}
		order, _ := epMap["order"].(float64)
		epsTitle, _ := epMap["title"].(string)
		if epsTitle == "" {
			epsTitle = fmt.Sprintf("第%d话", int(order))
		}

		// 获取章节页面
		pagesResp, err := m.client.GetComicPages(task.BookID, fmt.Sprintf("%d", int(order)), 1)
		if err != nil {
			continue
		}

		pagesData, _ := pagesResp.Data["pages"].(map[string]any)
		if pagesData == nil {
			continue
		}

		pageDocs, _ := pagesData["docs"].([]any)
		if len(pageDocs) == 0 {
			continue
		}

		// 创建章节目录
		epsDir := filepath.Join(comicDir, sanitizeFilename(epsTitle))
		os.MkdirAll(epsDir, 0755)

		// 下载每一页
		for j, page := range pageDocs {
			pageMap, ok := page.(map[string]any)
			if !ok {
				continue
			}

			media, _ := pageMap["media"].(map[string]any)
			if media == nil {
				continue
			}

			fileServer, _ := media["fileServer"].(string)
			path, _ := media["path"].(string)
			if fileServer == "" || path == "" {
				continue
			}

			imageURL := pica.GetImageURL(fileServer, path)
			filename := fmt.Sprintf("%03d.jpg", j+1)
			filePath := filepath.Join(epsDir, filename)

			// 检查文件是否已存在
			if _, err := os.Stat(filePath); err == nil {
				m.mu.Lock()
				task.Downloaded++
				task.UpdatedAt = time.Now()
				m.mu.Unlock()
				continue
			}

			// 下载图片
			if err := m.downloadImage(imageURL, filePath); err != nil {
				// 单个图片失败不影响整体
				continue
			}

			m.mu.Lock()
			task.Downloaded++
			task.UpdatedAt = time.Now()
			m.mu.Unlock()
		}

		// 保存章节信息
		infoPath := filepath.Join(epsDir, "info.json")
		infoData := fmt.Sprintf(`{"comic_id":"%s","eps_order":%d,"title":"%s","pages":%d}`,
			task.BookID, int(order), epsTitle, len(pageDocs))
		os.WriteFile(infoPath, []byte(infoData), 0644)

		_ = i // suppress unused warning
	}

	// 保存漫画信息
	comicInfoPath := filepath.Join(comicDir, "info.json")
	comicInfo := fmt.Sprintf(`{"id":"%s","title":"%s","downloaded_at":"%s"}`,
		task.BookID, task.Title, time.Now().Format(time.RFC3339))
	os.WriteFile(comicInfoPath, []byte(comicInfo), 0644)

	m.mu.Lock()
	task.Status = "completed"
	task.UpdatedAt = time.Now()
	m.mu.Unlock()
}

// downloadImage 下载单张图片
func (m *Manager) downloadImage(url, filePath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "okhttp/3.8.1")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// 写入临时文件后重命名，防止部分写入
	tmpPath := filePath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, filePath)
}

func (m *Manager) setError(task *DownloadTask, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	task.Status = "error"
	task.Error = msg
	task.UpdatedAt = time.Now()
}

func sanitizeFilename(name string) string {
	// 替换非法文件名字符
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, c := range invalid {
		result = replaceAll(result, c, "_")
	}
	if len(result) > 200 {
		result = result[:200]
	}
	return result
}

func replaceAll(s, old, new string) string {
	result := ""
	for _, c := range s {
		if string(c) == old {
			result += new
		} else {
			result += string(c)
		}
	}
	return result
}
