package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Mogvl/bika/backend/download"
	"github.com/Mogvl/bika/backend/handler"
	"github.com/Mogvl/bika/backend/pica"
)

//go:embed all:static
var embeddedStatic embed.FS

func main() {
	client := pica.NewClient()

	// 初始化 API 客户端（获取真实 IP）
	if err := client.Init(); err != nil {
		log.Printf("警告: API 初始化失败: %v", err)
		log.Println("将使用默认 DNS 连接，可能无法正常工作")
	}

	authHandler := handler.NewAuthHandler(client)
	comicsHandler := handler.NewComicsHandler(client)
	imageHandler := handler.NewImageHandler(client)
	gameHandler := handler.NewGameHandler(client)

	// 初始化下载管理器
	downloadDir := os.Getenv("DOWNLOAD_DIR")
	if downloadDir == "" {
		downloadDir = "downloads"
	}
	os.MkdirAll(downloadDir, 0755)
	downloadManager := download.NewManager(client, downloadDir)
	downloadHandler := handler.NewDownloadHandler(downloadManager)

	mux := http.NewServeMux()

	// ==================== API 路由 ====================

	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/profile", handler.AuthMiddleware(client, authHandler.Profile))
	mux.HandleFunc("/api/auth/punch-in", handler.AuthMiddleware(client, authHandler.PunchIn))
	mux.HandleFunc("/api/auth/change-password", handler.AuthMiddleware(client, authHandler.ChangePassword))
	mux.HandleFunc("/api/auth/forgot-password", authHandler.ForgotPassword)

	mux.HandleFunc("/api/categories", handler.AuthMiddleware(client, comicsHandler.Categories))

	mux.HandleFunc("/api/comics", handler.AuthMiddleware(client, comicsHandler.ListByCategory))
	mux.HandleFunc("/api/comics/search", handler.AuthMiddleware(client, comicsHandler.Search))
	mux.HandleFunc("/api/comics/leaderboard", handler.AuthMiddleware(client, comicsHandler.Leaderboard))
	mux.HandleFunc("/api/comics/random", handler.AuthMiddleware(client, comicsHandler.Random))
	mux.HandleFunc("/api/comics/collections", handler.AuthMiddleware(client, comicsHandler.Collections))
	mux.HandleFunc("/api/comics/keywords", handler.AuthMiddleware(client, comicsHandler.Keywords))

	mux.HandleFunc("/api/comics/{id}", handler.AuthMiddleware(client, comicsHandler.Detail))
	mux.HandleFunc("/api/comics/{id}/eps", handler.AuthMiddleware(client, comicsHandler.Eps))
	mux.HandleFunc("/api/comics/{id}/eps/{epsId}/pages", handler.AuthMiddleware(client, comicsHandler.Pages))
	mux.HandleFunc("/api/comics/{id}/comments", handler.AuthMiddleware(client, comicsHandler.Comments))
	mux.HandleFunc("/api/comics/{id}/comments/send", handler.AuthMiddleware(client, comicsHandler.SendComment))
	mux.HandleFunc("/api/comics/{id}/like", handler.AuthMiddleware(client, comicsHandler.LikeComic))

	mux.HandleFunc("/api/comments/{id}/like", handler.AuthMiddleware(client, comicsHandler.LikeComment))

	// 游戏接口
	mux.HandleFunc("/api/games", handler.AuthMiddleware(client, gameHandler.List))
	mux.HandleFunc("/api/games/{id}", handler.AuthMiddleware(client, gameHandler.Detail))
	mux.HandleFunc("/api/games/{id}/eps", handler.AuthMiddleware(client, gameHandler.Eps))
	mux.HandleFunc("/api/games/{id}/eps/{epsId}/pages", handler.AuthMiddleware(client, gameHandler.Pages))
	mux.HandleFunc("/api/games/{id}/comments", handler.AuthMiddleware(client, gameHandler.Comments))

	mux.HandleFunc("/api/favourites", handler.AuthMiddleware(client, comicsHandler.Favourites))
	mux.HandleFunc("/api/comics/{id}/favourite", handler.AuthMiddleware(client, comicsHandler.AddFavourite))

	// 下载接口
	mux.HandleFunc("/api/downloads", handler.AuthMiddleware(client, downloadHandler.List))
	mux.HandleFunc("/api/downloads/add", handler.AuthMiddleware(client, downloadHandler.Add))
	mux.HandleFunc("/api/downloads/{id}", handler.AuthMiddleware(client, downloadHandler.Status))
	mux.HandleFunc("/api/downloads/{id}/cancel", handler.AuthMiddleware(client, downloadHandler.Cancel))
	mux.HandleFunc("/api/downloads/{id}/remove", handler.AuthMiddleware(client, downloadHandler.Remove))

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

	if err := http.ListenAndServe(addr, corsHandler); err != nil {
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
