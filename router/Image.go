package router

import (
	"github.com/gin-gonic/gin"

	"Yc-GopherAI/controller/image"
)

func ImageRouter(r *gin.RouterGroup) {
	r.POST("/recognize", image.RecognizeImage)
}
