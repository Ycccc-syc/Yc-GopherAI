package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCheck(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/healthz", nil)

	Check(c)

	// 测试环境没有 MySQL 和 Redis，应该返回 503
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("状态码 = %d, 期望 %d", w.Code, http.StatusServiceUnavailable)
	}
	expected := `{"mysql":false,"redis":false,"service":"GopherAI","status":"degraded"}`
	if w.Body.String() != expected {
		t.Errorf("返回体 = %s", w.Body.String())
	}
}
