package router

import (
	"GopherAI/controller/health"
	"GopherAI/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	// Swagger API 文档
	r.Static("/docs", "./docs")

	// 健康检查 — 不需要鉴权，负载均衡器和监控系统会定期调用
	r.GET("/api/v1/healthz", health.Check)

	enterRouter := r.Group("/api/v1")
	{
		RegisterUserRouter(enterRouter.Group("/user"))
	}
	// 需要 JWT 鉴权的接口
	{
		AIGroup := enterRouter.Group("/AI")
		AIGroup.Use(jwt.Auth())
		AIRouter(AIGroup)
	}

	{
		ImageGroup := enterRouter.Group("/image")
		ImageGroup.Use(jwt.Auth())
		ImageRouter(ImageGroup)
	}

	return r
}
