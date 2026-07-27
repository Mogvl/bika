package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	SavePath    string    `json:"savePath"`
	TotalPages  int       `json:"totalPages"`
	Downloaded  int       `json:"downloaded"`
	Status      string    `json:"status"`
	Error       string    `json:"error,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Manager 下载管理器
type Manager struct {
	client      *pica.Client
	downloadDir string
	tasks       map[string]*DownloadTask
	mu          sync.RWMutex
	httpClient  *http.Client
}

// NewManager 创建下载管理器
func NewManager(client *pica.Client, downloadDir string) *Manager {
	log.Printf("[下载管理器] 初始化, 目录: %s", downloadDir)
	// 确保目录存在
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		log.Printf("[下载管理器] 创建目录失败: %v", err)
	}
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

	if task, exists := m.tasks[bookID]; exists {
		if task.Status == "completed" || task.Status == "error" {
			task.Status = "waiting"
			task.Error = ""
			task.Downloaded = 0
			task.UpdatedAt = time.Now()
		}
		return task
	}

	comicDir := filepath.Join(m.downloadDir, sanitizeFilename(title))

	task := &DownloadTask{
		ID:        bookID,
		BookID:    bookID,
		Title:     title,
		CoverURL:  coverURL,
		SavePath:  comicDir,
		Status:    "waiting",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.tasks[bookID] = task

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

	log.Printf("[下载] 开始: %s (%s)", task.Title, task.BookID)

	// 获取章节列表 (最多 50 章)
	epsResp, err := m.client.GetComicEps(task.BookID, 1)
	if err != nil {
		m.setError(task, fmt.Sprintf("获取章节失败: %v", err))
		log.Printf("[下载] 获取章节失败: %v", err)
		return
	}

	epsData, _ := epsResp.Data["eps"].(map[string]any)
	if epsData == nil {
		m.setError(task, "章节数据格式错误")
		log.Printf("[下载] 章节数据格式错误")
		return
	}

	docs, _ := epsData["docs"].([]any)
	if len(docs) == 0 {
		m.setError(task, "没有章节")
		log.Printf("[下载] 没有章节")
		return
	}

	log.Printf("[下载] 共 %d 个章节", len(docs))

	// 创建漫画目录
	comicDir := filepath.Join(m.downloadDir, sanitizeFilename(task.Title))
	if err := os.MkdirAll(comicDir, 0755); err != nil {
		m.setError(task, fmt.Sprintf("创建目录失败: %v", err))
		log.Printf("[下载] 创建目录失败: %v", err)
		return
	}
	log.Printf("[下载] 保存目录: %s", comicDir)

	totalDownloaded := 0

	// 下载每个章节
	for _, ep := range docs {
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
			log.Printf("[下载] 获取章节 %d 页面失败: %v", int(order), err)
			continue
		}

		pagesData, _ := pagesResp.Data["pages"].(map[string]any)
		if pagesData == nil {
			log.Printf("[下载] 章节 %d 页面数据为空", int(order))
			continue
		}

		pageDocs, _ := pagesData["docs"].([]any)
		if len(pageDocs) == 0 {
			log.Printf("[下载] 章节 %d 没有页面", int(order))
			continue
		}

		// 创建章节目录
		epsDir := filepath.Join(comicDir, sanitizeFilename(epsTitle))
		if err := os.MkdirAll(epsDir, 0755); err != nil {
			log.Printf("[下载] 创建章节目录失败: %v", err)
			continue
		}

		log.Printf("[下载] 章节 %d: %s (%d 页)", int(order), epsTitle, len(pageDocs))

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
				totalDownloaded++
				m.mu.Lock()
				task.Downloaded = totalDownloaded
				task.UpdatedAt = time.Now()
				m.mu.Unlock()
				continue
			}

			// 下载图片
			if err := m.downloadImage(imageURL, filePath); err != nil {
				log.Printf("[下载] 图片失败 %s: %v", filename, err)
				continue
			}

			totalDownloaded++
			m.mu.Lock()
			task.Downloaded = totalDownloaded
			task.UpdatedAt = time.Now()
			m.mu.Unlock()
		}

	}

	_ = comicDir

	m.mu.Lock()
	task.Status = "completed"
	task.UpdatedAt = time.Now()
	m.mu.Unlock()

	log.Printf("[下载] 完成: %s, 共下载 %d 张图片", task.Title, totalDownloaded)
}

// downloadImage 下载单张图片
func (m *Manager) downloadImage(url, filePath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "okhttp/3.8.1")
	req.Header.Set("Referer", "https://picaapi.picacomic.com/")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// 写入临时文件后重命名
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
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, c := range invalid {
		result = strings.ReplaceAll(result, c, "_")
	}
	if len(result) > 200 {
		result = result[:200]
	}
	return result
}
