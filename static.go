package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

func static() {
	// 将嵌入的文件系统转换为http.FileSystem
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Error creating sub filesystem: %v", err)
	}

	// 创建文件服务器
	fileServer := http.FileServer(http.FS(staticFS))

	// 定义HTTP处理函数
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
}
