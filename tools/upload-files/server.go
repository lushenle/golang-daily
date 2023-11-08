package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const (
	uploadPath = "./uploads" // 上传文件存储目录
	username   = "user"      // 基本认证用户名
	password   = "pass"      // 基本认证密码
)

// basicAuth 中间件用于验证HTTP Basic Authentication
func basicAuth(ctx *gin.Context) {
	u, p, hasAuth := ctx.Request.BasicAuth()
	if hasAuth && u == username && p == password {
		ctx.Next()
		return
	}

	// 请求需要HTTP基本认证
	ctx.Header("WWW-Authenticate", `Basic realm="Authorization Required"`)
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

// uploadHandler 处理文件上传
func uploadHandler(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("my_files")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad request")
		return
	}
	filename := header.Filename

	// 确保文件名是安全的，防止目录遍历攻击
	safeFilename := filepath.Base(filename)
	// 创建目标文件
	out, err := os.Create(filepath.Join(uploadPath, safeFilename))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer out.Close()

	// 将上传文件复制到目标文件
	_, err = io.Copy(out, file)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx.String(http.StatusOK, "File uploaded successfully: %s", safeFilename)
}

func main() {
	// 确保上传目录存在
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	// 静态文件服务，设置为上传目录
	router.Static("/files", uploadPath)

	// 上传文件的路由，使用 basicAuth 中间件
	router.POST("/upload", basicAuth, uploadHandler)

	router.Run(":8080")
}
