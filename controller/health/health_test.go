package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCheck(t *testing.T) {
	// 1. 创建测试接收器
	w := httptest.NewRecorder()
	// 2. 创建 Gin 测试上下文
	c, _ := gin.CreateTestContext(w)
	// 3. 设置请求方法
	c.Request = httptest.NewRequest("GET", "/healthz", nil)

	// 4. 调用 handler
	Check(c)

	// 5. 验证状态码
	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %d, 期望 %d", w.Code, http.StatusOK)
	}
}
