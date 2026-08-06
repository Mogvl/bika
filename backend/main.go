package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Mogvl/bika/backend/chat"
	"github.com/Mogvl/bika/backend/download"
	"github.com/Mogvl/bika/backend/fried"
	"github.com/Mogvl/bika/backend/handler"
	"github.com/Mogvl/bika/backend/pica"
)

// statusRecorder 记录响应状态码
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

//go:embed all:static
var embeddedStatic embed.FS

func main() {
	client := pica.NewClient()

	// 初始化 API 客户端（获取真实 IP）
	if err := client.Init(); err != nil {
		log.Printf("警告: API 初始化失败: %v", err)
		log.Println("将使用默认 DNS 连接，可能无法正常工作")
	}

	// 初始化下载管理器
	downloadDir := os.Getenv("DOWNLOAD_DIR")
	if downloadDir == "" {
		downloadDir = "downloads"
	}
	os.MkdirAll(downloadDir, 0755)
	downloadManager := download.NewManager(client, downloadDir)
	downloadHandler := handler.NewDownloadHandler(downloadManager)
	localHandler := handler.NewLocalHandler(downloadManager)

	// 初始化聊天客户端
	chatClient := chat.NewChatClient()
	chatHandler := handler.NewChatHandler(chatClient)

	// 初始化锅贴客户端
	friedClient := fried.NewFriedClient()
	friedHandler := handler.NewFriedHandler(friedClient)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(client, friedClient)
	comicsHandler := handler.NewComicsHandler(client)
	imageHandler := handler.NewImageHandler(client)
	gameHandler := handler.NewGameHandler(client)

	mux := http.NewServeMux()

	// ==================== API 路由 ====================

	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/profile", handler.AuthMiddleware(client, authHandler.Profile))
	mux.HandleFunc("/api/auth/punch-in", handler.AuthMiddleware(client, authHandler.PunchIn))
	mux.HandleFunc("/api/auth/change-password", handler.AuthMiddleware(client, authHandler.ChangePassword))
	mux.HandleFunc("/api/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("/api/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("/api/auth/avatar", handler.AuthMiddleware(client, authHandler.SetAvatar))
	mux.HandleFunc("/api/auth/title", handler.AuthMiddleware(client, authHandler.SetTitle))
	mux.HandleFunc("/api/auth/my-comments", handler.AuthMiddleware(client, authHandler.MyComments))

	mux.HandleFunc("/api/categories", handler.AuthMiddleware(client, comicsHandler.Categories))

	mux.HandleFunc("/api/comics", handler.AuthMiddleware(client, comicsHandler.ListByCategory))
	mux.HandleFunc("/api/comics/search", handler.AuthMiddleware(client, comicsHandler.Search))
	mux.HandleFunc("/api/comics/leaderboard", handler.AuthMiddleware(client, comicsHandler.Leaderboard))
	mux.HandleFunc("/api/comics/knight-leaderboard", handler.AuthMiddleware(client, comicsHandler.KnightLeaderboard))
	mux.HandleFunc("/api/comics/random", handler.AuthMiddleware(client, comicsHandler.Random))
	mux.HandleFunc("/api/comics/collections", handler.AuthMiddleware(client, comicsHandler.Collections))
	mux.HandleFunc("/api/comics/keywords", handler.AuthMiddleware(client, comicsHandler.Keywords))

	mux.HandleFunc("/api/comics/{id}", handler.AuthMiddleware(client, comicsHandler.Detail))
	mux.HandleFunc("/api/comics/{id}/eps", handler.AuthMiddleware(client, comicsHandler.Eps))
	mux.HandleFunc("/api/comics/{id}/eps/{epsId}/pages", handler.AuthMiddleware(client, comicsHandler.Pages))
	mux.HandleFunc("/api/comics/{id}/comments", handler.AuthMiddleware(client, comicsHandler.Comments))
	mux.HandleFunc("/api/comics/{id}/comments/send", handler.AuthMiddleware(client, comicsHandler.SendComment))
	mux.HandleFunc("/api/comics/{id}/like", handler.AuthMiddleware(client, comicsHandler.LikeComic))
	mux.HandleFunc("/api/comics/{id}/recommendation", handler.AuthMiddleware(client, comicsHandler.Recommendation))

	mux.HandleFunc("/api/comments/{id}/like", handler.AuthMiddleware(client, comicsHandler.LikeComment))
	mux.HandleFunc("/api/comments/{id}/childrens", handler.AuthMiddleware(client, comicsHandler.SubComments))
	mux.HandleFunc("/api/comments/{id}/childrens/send", handler.AuthMiddleware(client, comicsHandler.SendSubComment))
	mux.HandleFunc("/api/comments/{id}/report", handler.AuthMiddleware(client, comicsHandler.ReportComment))

	// 游戏接口
	mux.HandleFunc("/api/games", handler.AuthMiddleware(client, gameHandler.List))
	mux.HandleFunc("/api/games/{id}", handler.AuthMiddleware(client, gameHandler.Detail))
	mux.HandleFunc("/api/games/{id}/eps", handler.AuthMiddleware(client, gameHandler.Eps))
	mux.HandleFunc("/api/games/{id}/eps/{epsId}/pages", handler.AuthMiddleware(client, gameHandler.Pages))
	mux.HandleFunc("/api/games/{id}/comments", handler.AuthMiddleware(client, gameHandler.Comments))
	mux.HandleFunc("/api/games/{id}/comments/send", handler.AuthMiddleware(client, gameHandler.SendComment))
	mux.HandleFunc("/api/game-comments/{id}/like", handler.AuthMiddleware(client, gameHandler.LikeComment))

	mux.HandleFunc("/api/favourites", handler.AuthMiddleware(client, comicsHandler.Favourites))
	mux.HandleFunc("/api/comics/{id}/favourite", handler.AuthMiddleware(client, comicsHandler.AddFavourite))

	// 下载接口
	mux.HandleFunc("/api/downloads", handler.AuthMiddleware(client, downloadHandler.List))
	mux.HandleFunc("/api/downloads/add", handler.AuthMiddleware(client, downloadHandler.Add))
	mux.HandleFunc("/api/downloads/{id}", handler.AuthMiddleware(client, downloadHandler.Status))
	mux.HandleFunc("/api/downloads/{id}/cancel", handler.AuthMiddleware(client, downloadHandler.Cancel))
	mux.HandleFunc("/api/downloads/{id}/resume", handler.AuthMiddleware(client, downloadHandler.Resume))
	mux.HandleFunc("/api/downloads/{id}/remove", handler.AuthMiddleware(client, downloadHandler.Remove))

	// 本地库接口
	mux.HandleFunc("/api/local/list", handler.AuthMiddleware(client, localHandler.List))
	mux.HandleFunc("/api/local/eps", handler.AuthMiddleware(client, localHandler.Eps))
	mux.HandleFunc("/api/local/image", handler.AuthMiddleware(client, localHandler.Image))

	// 聊天接口
	mux.HandleFunc("/api/chat/login", handler.AuthMiddleware(client, chatHandler.Login))
	mux.HandleFunc("/api/chat/rooms", handler.AuthMiddleware(client, chatHandler.Rooms))
	mux.HandleFunc("/api/chat/messages", handler.AuthMiddleware(client, chatHandler.Messages))
	mux.HandleFunc("/api/chat/send", handler.AuthMiddleware(client, chatHandler.SendMessage))
	mux.HandleFunc("/api/chat/profile", handler.AuthMiddleware(client, chatHandler.Profile))

	// 锅贴接口（好友动态）
	mux.HandleFunc("/api/fried/posts", handler.AuthMiddleware(client, friedHandler.Posts))
	mux.HandleFunc("/api/fried/posts/{id}/comments", handler.AuthMiddleware(client, friedHandler.Comments))
	mux.HandleFunc("/api/fried/posts/{id}/comments/send", handler.AuthMiddleware(client, friedHandler.SendComment))
	mux.HandleFunc("/api/fried/comments/{id}/like", handler.AuthMiddleware(client, friedHandler.LikeComment))

	mux.HandleFunc("/api/image/proxy", imageHandler.Proxy)
	mux.HandleFunc("/api/image/url", imageHandler.ImageURL)

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","version":"1.0.0"}`))
	})

	// ==================== 静态文件服务 ====================

	// 使用嵌入式静态文件
	staticFS, err := fs.Sub(embeddedStatic, "static")
	if err == nil {
		if entries, _ := fs.ReadDir(staticFS, "."); len(entries) > 0 {
			log.Println("使用内置前端文件")
			mux.Handle("/", handler.SPAFileServer(staticFS))
		} else {
			log.Println("内置前端目录为空，尝试外部目录")
			serveExternalStatic(mux)
		}
	} else {
		log.Printf("无法加载内置前端: %v，尝试外部目录", err)
		serveExternalStatic(mux)
	}

	// CORS
	corsHandler := handler.CORSMiddleware(mux)

	// 请求日志中间件
	loggedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		corsHandler.ServeHTTP(rec, r)
		log.Printf("[%s] %s %s → %d (%s)", r.Method, r.URL.Path, r.URL.RawQuery, rec.status, time.Since(start))
	})

	// ==================== 启动服务器 ====================

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	addr := ":" + port
	log.Printf("==========================================")
	log.Printf("  哔咔漫画 Web 版 (PicACG Web)")
	log.Printf("  服务地址: http://localhost:%s", port)
	log.Printf("  API 地址: http://localhost:%s/api", port)
	log.Printf("==========================================")

	if err := http.ListenAndServe(addr, loggedHandler); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func serveExternalStatic(mux *http.ServeMux) {
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		exe, _ := os.Executable()
		frontendDir = filepath.Join(filepath.Dir(exe), "static")
		if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
			frontendDir = "static"
		}
	}

	if stat, err := os.Stat(frontendDir); err == nil && stat.IsDir() {
		if entries, _ := os.ReadDir(frontendDir); len(entries) > 0 {
			log.Printf("从外部目录加载前端文件: %s", frontendDir)
			mux.Handle("/", http.FileServer(http.Dir(frontendDir)))
		} else {
			log.Println("外部静态文件目录为空")
		}
	} else {
		log.Println("未找到前端文件，仅 API 模式运行")
	}
}
