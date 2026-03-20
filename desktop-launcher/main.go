package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"
)

//go:embed all:dist
var embeddedDist embed.FS

func main() {
	distFS, err := fs.Sub(embeddedDist, "dist")
	if err != nil {
		log.Fatalf("load embedded dist failed: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}
	defer listener.Close()

	server := &http.Server{
		Handler: newSPAHandler(distFS),
	}

	go func() {
		if serveErr := server.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			log.Printf("server stopped with error: %v", serveErr)
		}
	}()

	url := fmt.Sprintf("http://%s/chart", listener.Addr().String())
	fmt.Printf("AI KLine Web 已启动: %s\n", url)
	fmt.Println("关闭这个窗口即可退出。")

	if os.Getenv("AIKLINE_NO_BROWSER") == "" {
		if err := openBrowser(url); err != nil {
			log.Printf("open browser failed: %v", err)
		}
	}

	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func newSPAHandler(staticFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(staticFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean(r.URL.Path)
		if cleanPath == "." || cleanPath == "/" {
			http.ServeFileFS(w, r, staticFS, "index.html")
			return
		}

		filePath := cleanPath[1:]
		if _, err := fs.Stat(staticFS, filePath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFileFS(w, r, staticFS, "index.html")
	})
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
