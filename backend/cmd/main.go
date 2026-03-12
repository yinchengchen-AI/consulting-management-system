package main

import (
	"os"
	"os/signal"
	"syscall"

	"consulting-system/config"
	"consulting-system/internal/database"
	"consulting-system/internal/handlers"
	"consulting-system/internal/middleware"
	"consulting-system/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置日志
	setupLogger(cfg)
	zapLogger := logger.NewLogger()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		logrus.Fatal("Failed to initialize database: ", err)
	}
	defer database.Close()

	// 执行数据库迁移
	if err := database.Migrate(); err != nil {
		logrus.Fatal("Failed to migrate database: ", err)
	}

	// 初始化基础数据
	if err := database.SeedData(); err != nil {
		logrus.Error("Failed to seed data: ", err)
	}

	// 创建Gin引擎
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.Recovery(zapLogger))
	r.Use(middleware.Logger(zapLogger))
	r.Use(middleware.CORS(&cfg.CORS))
	r.Use(middleware.ErrorLogger())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// API路由组
	api := r.Group("/api")

	// 公开路由（不需要认证）
	public := api.Group("")
	{
		authHandler := handlers.NewAuthHandler()
		authHandler.RegisterRoutes(public)
	}

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.JWTAuth())
	{
		// 用户管理
		userHandler := handlers.NewUserHandler()
		userHandler.RegisterRoutes(protected)

		// 角色管理
		roleHandler := handlers.NewRoleHandler()
		roleHandler.RegisterRoutes(protected)

		// 客户管理
		customerHandler := handlers.NewCustomerHandler()
		customerHandler.RegisterRoutes(protected)

		// 服务类型管理
		serviceTypeHandler := handlers.NewServiceTypeHandler()
		serviceTypeHandler.RegisterRoutes(protected)

		// 服务订单管理
		serviceOrderHandler := handlers.NewServiceOrderHandler()
		serviceOrderHandler.RegisterRoutes(protected)

		// 发票管理
		invoiceHandler := handlers.NewInvoiceHandler()
		invoiceHandler.RegisterRoutes(protected)

		// 收款管理
		paymentHandler := handlers.NewPaymentHandler()
		paymentHandler.RegisterRoutes(protected)

		// 统计管理
		statisticsHandler := handlers.NewStatisticsHandler()
		statisticsHandler.RegisterRoutes(protected)

		// 通知管理
		noticeHandler := handlers.NewNoticeHandler()
		noticeHandler.RegisterRoutes(protected)

		// 文档管理
		documentHandler := handlers.NewDocumentHandler()
		documentHandler.RegisterRoutes(protected)

		// 合同管理
		contractHandler := handlers.NewContractHandler()
		contractHandler.RegisterRoutes(protected)

		// 系统设置
		settingHandler := handlers.NewSettingHandler()
		settingHandler.RegisterRoutes(protected)
	}

	// 启动服务器
	srv := &Server{
		router: r,
		config: cfg,
	}

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			logrus.Fatal("Failed to start server: ", err)
		}
	}()

	logrus.Infof("Server started on port %s", cfg.Server.Port)

	<-quit
	logrus.Info("Shutting down server...")
}

// setupLogger 设置日志
func setupLogger(cfg *config.Config) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
}

// Server HTTP服务器
type Server struct {
	router *gin.Engine
	config *config.Config
}

// Start 启动服务器
func (s *Server) Start() error {
	return s.router.Run(":" + s.config.Server.Port)
}
