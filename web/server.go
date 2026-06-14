//go:build web

package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scanner/core/common"
	"scanner/core/common/i18n"
	"scanner/core/web/api"
	"scanner/core/web/ws"
)

//go:embed dist/*
var distFS embed.FS

// StartServer 启动Web服务器
func StartServer(port int) error {
	// 初始化WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// 创建路由
	mux := http.NewServeMux()

	// API路由
	api.RegisterRoutes(mux, hub)

	// WebSocket路由
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	// 静态文件服务
	distContent, err := fs.Sub(distFS, "dist")
	if err != nil {
		return fmt.Errorf("failed to get dist fs: %w", err)
	}
	fileServer := http.FileServer(http.FS(distContent))

	// SPA fallback: 对于非API/WS请求，尝试静态文件，否则返回index.html
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 检查文件是否存在
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// 尝试打开文件
		f, err := distContent.Open(path[1:]) // 移除开头的/
		if err != nil {
			// 文件不存在，返回index.html（SPA路由）
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		f.Close()

		// 文件存在，正常服务
		fileServer.ServeHTTP(w, r)
	})

	// 创建服务器
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:         addr,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 优雅关闭
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		common.LogInfo(i18n.GetText("web_shutting_down"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			common.LogError(fmt.Sprintf("Server shutdown error: %v", err))
		}
		close(done)
	}()

	// 启动服务器
	common.LogSuccess(i18n.Tr("web_server_started", port))
	fmt.Printf("    http://localhost:%d\n", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	<-done
	return nil
}

// corsMiddleware 添加CORS头（开发时需要）
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
