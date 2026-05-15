package health

import (
	"GopherAI/common/mysql"
	"GopherAI/common/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Check 健康检查接口
// @Summary      健康检查
// @Description  检查服务及各依赖组件（MySQL、Redis）的运行状态
// @Tags         系统
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Success      503 {object} map[string]interface{}
// @Router       /healthz [get]

func Check(c *gin.Context) {
	status := "ok"
	httpCode := http.StatusOK

	// 检查 MySQL 连接
	dbOK := true
	if mysql.DB == nil {
		dbOK = false
	} else {
		db, err := mysql.DB.DB()
		if err != nil || db.Ping() != nil {
			dbOK = false
		}
	}
	if !dbOK {
		status = "degraded"
		httpCode = http.StatusServiceUnavailable
	}

	// 检查 Redis 连接
	redisOK := true
	if redis.Rdb != nil {
		if _, err := redis.Rdb.Ping(redis.Ctx).Result(); err != nil {
			redisOK = false
		}
	}
	if !redisOK {
		status = "degraded"
		httpCode = http.StatusServiceUnavailable
	}

	c.JSON(httpCode, gin.H{
		"status":  status,
		"service": "GopherAI",
		"mysql":   dbOK,
		"redis":   redisOK,
	})
}
