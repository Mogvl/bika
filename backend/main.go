package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Mogvl/bika/backend/handler"
	"github.com/Mogvl/bika/backend/pica"
)

//go:embed static/*
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

	mux := http.NewServeMux()

	// ==================== API 路由 ====================

	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/profile", handler.AuthMiddleware(client, authHandler.Profile))
	mux.HandleFunc("/api/auth/punch-in", handler.AuthMiddleware(client, authHandler.PunchIn))

	mux.HandleFunc("/api/categories", comicsHandler.Categories)

	mux.HandleFunc("/api/comics", comicsHandler.ListByCategory)
	mux.HandleFunc("/api/comics/search", comicsHandler.Search)
	mux.HandleFunc("/api/comics/leaderboard", comicsHandler.Leaderboard)
	mux.HandleFunc("/api/comics/random", comicsHandler.Random)
	mux.HandleFunc("/api/comics/collections", comicsHandler.Collections)
	mux.HandleFunc("/api/comics/keywords", comicsHandler.Keywords)

	mux.HandleFunc("/api/comics/{id}", comicsHandler.Detail)
	mux.HandleFunc("/api/comics/{id}/eps", comicsHandler.Eps)
	mux.HandleFunc("/api/comics/{id}/eps/{epsId}/pages", comicsHandler.Pages)
	mux.HandleFunc("/api/comics/{id}/comments", comicsHandler.Comments)

	mux.HandleFunc("/api/favourites", handler.AuthMiddleware(client, comicsHandler.Favourites))
	mux.HandleFunc("/api/comics/{id}/favourite", handler.AuthMiddleware(client, comicsHandler.AddFavourite))

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
