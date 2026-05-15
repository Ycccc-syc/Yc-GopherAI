// @title           GopherAI API
// @version         1.0
// @description     AI 聊天助手后端 API
// @host            localhost:9090
// @BasePath        /api/v1
// @schemes         http

package main

import (
	"GopherAI/common/aihelper"
	"GopherAI/common/mysql"
	"GopherAI/common/rabbitmq"
	"GopherAI/common/redis"
	"GopherAI/config"
	"GopherAI/dao/message"
	_ "GopherAI/docs" //Swagger文档注册
	"GopherAI/router"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func StartServer(addr string, port int) (*http.Server, error) {
	// 生产环境切换 Gin 为 release 模式，减少不必要的日志输出
	gin.SetMode(gin.ReleaseMode)

	r := router.InitRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", addr, port),
		Handler: r,
	}

	// 在 goroutine 中启动 HTTP 服务，不阻塞后续逻辑
	go func() {
		slog.Info("server is listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			panic(err)
		}
	}()

	return srv, nil
}

// readDataFromDB 从数据库加载历史消息到 AIHelperManager 内存
func readDataFromDB() error {
	manager := aihelper.GetGlobalManager()
	msgs, err := message.GetAllMessages()
	if err != nil {
		return err
	}
	for i := range msgs {
		m := &msgs[i]
		// 默认使用 OpenAI 模型
		modelType := "1"
		configMap := make(map[string]interface{})
		helper, err := manager.GetOrCreateAIHelper(m.UserName, m.SessionID, modelType, configMap)
		if err != nil {
			slog.Error("failed to create AIHelper",
				"user", m.UserName,
				"session", m.SessionID,
				"error", err,
			)
			continue
		}
		slog.Debug("loaded session from DB", "session", helper.SessionID)
		helper.AddMessage(m.Content, m.UserName, m.IsUser, false)
	}
	slog.Info("AIHelperManager initialized", "total_messages", len(msgs))
	return nil
}

func main() {
	// ========== 1. 初始化结构化日志 ==========
	// slog 是 Go 1.21+ 内置的结构化日志库，输出 JSON 格式，方便生产环境日志采集
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// ========== 2. 加载配置 ==========
	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port

	// ========== 3. 初始化基础设施 ==========

	// 3.1 MySQL — 数据库连接
	if err := mysql.InitMysql(); err != nil {
		slog.Error("mysql initialization failed", "error", err)
		return
	}
	slog.Info("mysql connected")

	// 3.2 从 MySQL 加载历史消息到内存
	if err := readDataFromDB(); err != nil {
		slog.Error("failed to load messages from DB", "error", err)
	}

	// 3.3 Redis — 缓存（用于验证码等）
	redis.Init()
	slog.Info("redis connected")

	// 3.4 RabbitMQ — 消息队列（异步持久化聊天记录）
	rabbitmq.InitRabbitMQ()
	slog.Info("rabbitmq connected")

	// ========== 4. 启动 HTTP 服务 ==========
	srv, err := StartServer(host, port)
	if err != nil {
		slog.Error("server startup failed", "error", err)
		return
	}

	// ========== 5. 等待退出信号，优雅关闭 ==========
	// 优雅关闭的含义：收到 SIGINT/Ctrl+C 或 SIGTERM/kill 时，
	// 先停止接收新请求，等待正在处理的请求完成（最多 10 秒），再释放资源退出。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("shutting down...", "signal", sig.String())

	// 5.1 关闭 HTTP 服务（等待现有请求完成，最多 10 秒）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("http server forced shutdown", "error", err)
	}

	// 5.2 清理其他资源
	rabbitmq.DestroyRabbitMQ()
	slog.Info("server stopped gracefully")
}
