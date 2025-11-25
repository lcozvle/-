package handler

import (
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

// Shutdown: 关闭所有相关进程
func Shutdown(ctx *gin.Context) {
	// 先返回响应给前端
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "系统正在关闭..."})

	go func() {
		// 等待 1 秒，确保前端收到响应
		time.Sleep(1 * time.Second)

		// 执行 taskkill 命令关闭所有相关进程
		// 使用 taskkill 比 PowerShell 更直接，且能配合 cmd /c 自动关闭窗口
		exec.Command("taskkill", "/F", "/IM", "persistServerRun.exe").Run()
		exec.Command("taskkill", "/F", "/IM", "spiderRun.exe").Run()
		exec.Command("taskkill", "/F", "/IM", "web.exe").Run()
		// 关闭 go.exe 可能会影响其他 go 程序，但在专用环境下是可接受的
		exec.Command("taskkill", "/F", "/IM", "go.exe").Run()

		// 确保自己也退出
		os.Exit(0)
	}()
}
