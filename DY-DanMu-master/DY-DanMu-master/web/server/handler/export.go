package handler

import (
	"DY-DanMu/DMconfig/config"
	_type2 "DY-DanMu/dbConn/_type"
	"DY-DanMu/web/client"
	_type "DY-DanMu/web/server/_type"
	"DY-DanMu/web/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ExportDanmuForAI: 导出弹幕词频统计数据供 AI 分析
func ExportDanmuForAI(ctx *gin.Context) error {
	result := []_type2.BarrageStatisticsCountResult{}

	// 获取 Rid 参数
	rid := ctx.Query("rid")

	req := _type.DanmuFrequencyRequest{
		Limit: 5000,
		Rid:   rid,
	}

	// 调用 RPC 获取数据
	err := client.ClientRPC.Call(config.DYWebConfig.ExportDanmuFrequency, req, &result)
	if err != nil {
		err = util.RpcClientShutDownErrorhandler(err)
		if err == nil {
			err = client.ClientRPC.Call(config.DYWebConfig.ExportDanmuFrequency, req, &result)
		}
	}

	if err != nil {
		return ServerError()
	}

	// 设置响应头，触发文件下载
	fileName := fmt.Sprintf("danmu_analysis_%s.json", time.Now().Format("20060102_150405"))
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Type", "application/json")

	ctx.JSON(http.StatusOK, result)
	return nil
}
