package jwt

import (
	"log"
	"net/http"
	"strings"

	"Yc-GopherAI/common/code"
	"Yc-GopherAI/controller"
	"Yc-GopherAI/utils/myjwt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := new(controller.Response)

		var token string
		authHeader := c.GetHeader("Authorization")
		//源代码逻辑不够严谨，可能出现头出错但仍从URL中寻找token（我的修改是逻辑更完善）
		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
				c.Abort()
				return
			}
		} else {
			token = c.Query("token")
		}
		if token == "" {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		//只打印前缀，防止token泄露（我的修改）
		log.Printf("token received, prefix: %s...", token[:10])
		userName, ok := myjwt.ParseToken(token)
		if !ok {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		c.Set("userName", userName)
		c.Next()
	}
}
