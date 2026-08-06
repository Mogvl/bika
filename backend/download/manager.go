package download

import (
	"context"
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
	ID         string    `json:"id"`
	BookID     string    `json:"bookId"`
	Title      string    `json:"title"`
	CoverURL   string    `json:"coverUrl"`
	SavePath   string    `json:"savePath"`
	TotalPages int       `json:"totalPages"`
	Downloaded int       `json:"downloaded"`
	Speed      string    `json:"speed,omitempty"`
	Status     string    `json:"status"`
	Error      string    `json:"error,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Manager 下载管理器
type Manager struct {
	client      *pica.Client
	downloadDir string
	tasks       map[string]*DownloadTask
	cancels     map[string]context.CancelFunc
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
		cancels:     make(map[string]context.CancelFunc),
		httpClient:  &http.Client{Timeout: 60 * time.Second},
	}
}

// AddTask 添加下载任务
func (m *Manager) AddTask(bookID, title, coverURL string) *DownloadTask {
	m.mu.Lock()

	if task, exists := m.tasks[bookID]; exists {
		// 已存在且已完成/失败/暂停 → 重新开始
		if task.Status == "completed" || task.Status == "error" || task.Status == "paused" {
			task.Status = "waiting"
			task.Error = ""
			task.Downloaded = 0
			task.Speed = ""
			task.UpdatedAt = time.Now()
			m.mu.Unlock()
			go m.downloadComic(task)
			return task
		}
		m.mu.Unlock()
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
	m.mu.Unlock()

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

// GetDownloadDir 获取下载目录
func (m *Manager) GetDownloadDir() string {
	return m.downloadDir
}

// CancelTask 暂停下载（真正中断下载 goroutine）
func (m *Manager) CancelTask(bookID string) {
	m.mu.Lock()
	if cancel, ok := m.cancels[bookID]; ok {
		cancel()
		delete(m.cancels, bookID)
	}
	if task, exists := m.tasks[bookID]; exists {
		task.Status = "paused"
		task.Speed = ""
		task.UpdatedAt = time.Now()
	}
	m.mu.Unlock()
}

// ResumeTask 恢复暂停的任务
func (m *Manager) ResumeTask(bookID string) {
	m.mu.Lock()
	task, exists := m.tasks[bookID]
	if !exists {
		m.mu.Unlock()
		return
	}
	if task.Status == "paused" || task.Status == "error" {
		task.Status = "waiting"
		task.Error = ""
		task.Speed = ""
		task.UpdatedAt = time.Now()
		m.mu.Unlock()
		go m.downloadComic(task)
		return
	}
	m.mu.Unlock()
}

// RemoveTask 删除下载任务（可选删除文件）
func (m *Manager) RemoveTask(bookID string, deleteFile bool) {
	m.mu.Lock()
	if cancel, ok := m.cancels[bookID]; ok {
		cancel()
		delete(m.cancels, bookID)
	}
	task, exists := m.tasks[bookID]
	delete(m.tasks, bookID)
	m.mu.Unlock()

	if deleteFile && exists && task.SavePath != "" {
		// 等 goroutine 退出后再删
		time.Sleep(300 * time.Millisecond)
		if err := os.RemoveAll(task.SavePath); err != nil {
			log.Printf("[下载] 删除文件失败: %v", err)
		}
	}
}

// downloadComic 下载漫画
func (m *Manager) downloadComic(task *DownloadTask) {
	ctx, cancel := context.WithCancel(context.Background())
	m.mu.Lock()
	m.cancels[task.BookID] = cancel
	task.Status = "downloading"
	task.UpdatedAt = time.Now()
	m.mu.Unlock()
	defer func() {
		m.mu.Lock()
		if c, ok := m.cancels[task.BookID]; ok {
			c()
			delete(m.cancels, task.BookID)
		}
		m.mu.Unlock()
	}()

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

	// 计算总页数（兼容 pagesCount 字段缺失的情况）
	totalPages := 0
	for _, ep := range docs {
		epMap, ok := ep.(map[string]any)
		if !ok {
			continue
		}
		if pc, ok := epMap["pagesCount"].(float64); ok && int(pc) > 0 {
			totalPages += int(pc)
			continue
		}
		// 字段缺失时用已下载页数占位，最后按实际章节逐章累加
		totalPages += -1
	}
	if totalPages < 0 {
		totalPages = 0
	}
	m.mu.Lock()
	task.TotalPages = totalPages
	m.mu.Unlock()

	log.Printf("[下载] 共 %d 个章节, %d 页", len(docs), totalPages)

	// 创建漫画目录
	comicDir := filepath.Join(m.downloadDir, sanitizeFilename(task.Title))
	if err := os.MkdirAll(comicDir, 0755); err != nil {
		m.setError(task, fmt.Sprintf("创建目录失败: %v", err))
		log.Printf("[下载] 创建目录失败: %v", err)
		return
	}
	log.Printf("[下载] 保存目录: %s", comicDir)

	totalDownloaded := 0
	startTime := time.Now()
	var speedBytes float64
	var speedMu sync.Mutex

	// 下载每个章节
	for _, ep := range docs {
		select {
		case <-ctx.Done():
			log.Printf("[下载] 已暂停: %s", task.Title)
			return
		default:
		}

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

		// 累计实际页数（此前 pagesCount 缺失时算不到，这里补上）
		m.mu.Lock()
		task.TotalPages += len(pageDocs)
		m.mu.Unlock()

		// 创建章节目录
		epsDir := filepath.Join(comicDir, sanitizeFilename(epsTitle))
		if err := os.MkdirAll(epsDir, 0755); err != nil {
			log.Printf("[下载] 创建章节目录失败: %v", err)
			continue
		}

		log.Printf("[下载] 章节 %d: %s (%d 页)", int(order), epsTitle, len(pageDocs))

		// 下载每一页
		for j, page := range pageDocs {
			select {
			case <-ctx.Done():
				log.Printf("[下载] 已暂停: %s", task.Title)
				return
			default:
			}

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
			size, err := m.downloadImage(imageURL, filePath)
			if err != nil {
				log.Printf("[下载] 图片失败 %s: %v", filename, err)
				continue
			}

			speedMu.Lock()
			speedBytes += float64(size)
			speedMu.Unlock()
			totalDownloaded++
			m.mu.Lock()
			task.Downloaded = totalDownloaded
			// 计算速度
			elapsed := time.Since(startTime).Seconds()
			if elapsed > 0 {
				speed := speedBytes / elapsed / 1024 // KB/s
				task.Speed = fmt.Sprintf("%.1f KB/s", speed)
			}
			task.UpdatedAt = time.Now()
			m.mu.Unlock()
		}

	}

	_ = comicDir

	m.mu.Lock()
	task.Status = "completed"
	task.Speed = ""
	task.UpdatedAt = time.Now()
	m.mu.Unlock()

	log.Printf("[下载] 完成: %s, 共下载 %d 张图片", task.Title, totalDownloaded)
}

// downloadImage 下载单张图片，返回字节数
func (m *Manager) downloadImage(url, filePath string) (int64, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "okhttp/3.8.1")
	req.Header.Set("Referer", "https://picaapi.picacomic.com/")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// 写入临时文件后重命名
	tmpPath := filePath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		os.Remove(tmpPath)
		return 0, err
	}

	return n, os.Rename(tmpPath, filePath)
}

func (m *Manager) setError(task *DownloadTask, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	task.Status = "error"
	task.Error = msg
	task.Speed = ""
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
