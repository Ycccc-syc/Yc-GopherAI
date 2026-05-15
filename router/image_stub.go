//go:build !cgo

package router

import (
	"github.com/gin-gonic/gin"
)

// ImageRouter 空实现 — 没有 CGO 时图片识别功能不可用，不注册路由
func ImageRouter(r *gin.RouterGroup) {
	// 图片识别需要 onnxruntime 依赖，当前环境不支持
}
