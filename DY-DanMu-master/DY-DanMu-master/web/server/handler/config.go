package handler

import (
	"DY-DanMu/dbConn/redisConn"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigRidStruct struct {
	Rid string `json:"rid"`
}

// UpdateRid: 更新直播间ID配置
func UpdateRid(ctx *gin.Context) error {
	var data ConfigRidStruct
	err := ctx.BindJSON(&data)
	if err != nil || data.Rid == "" {
		return ParameterError("Invalid Rid")
	}

	// 获取 Redis 客户端
	redisClient := redisConn.DBConn()

	// 设置新的 Rid，过期时间设为永久或较长
	err = redisClient.Set("Spider:TargetRid", data.Rid, 0).Err()
	if err != nil {
		return ServerError()
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "配置已更新，Spider 将在下一次心跳时切换房间", "rid": data.Rid})
	return nil
}

// GetCurrentRid: 获取当前直播间ID
func GetCurrentRid(ctx *gin.Context) error {
	redisClient := redisConn.DBConn()
	rid, err := redisClient.Get("Spider:TargetRid").Result()
	if err != nil {
		// 如果 Redis 中没有，可能还没设置过，或者 Spider 正在使用默认配置
		// 这里简单返回空或者默认值
		rid = ""
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "rid": rid})
	return nil
}
