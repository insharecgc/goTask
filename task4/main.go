package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gotask/task4/config"
	"gotask/task4/handler"
	"gotask/task4/model"
	"gotask/task4/router"
	"gotask/task4/service"
)

func main() {
	// 1. 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 2. 初始化日志（zap）
	logger, err := initZapLogger()
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	// 3. 初始化数据库（GORM）
	db, err := initDB()
	if err != nil {
		logger.Fatal("初始化数据库失败", zap.Error(err))
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 4. 初始化服务和控制器
	userService := service.NewUserService(db)
	postService := service.NewPostService(db)
	commentService := service.NewCommentService(db)

	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	// 5. 初始化路由
	r := gin.New() // 不使用默认中间件（自己手动添加）
	router.Setup(r, userHandler, postHandler, commentHandler, logger)

	// 6. 启动服务器（优雅退出）
	srv := &http.Server{
		Addr:    config.Cfg.ServerConfig.Port,
		Handler: r,
	}
	go func() {
		logger.Info("服务器启动", zap.String("port", config.Cfg.ServerConfig.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号，优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("开始关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器关闭失败", zap.Error(err))
	}

	logger.Info("服务器已关闭")
}

// 初始化zap日志
func initZapLogger() (*zap.Logger, error) {
	var cfg zap.Config
	if config.Cfg.LogLevel == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.OutputPaths = []string{
		"stdout",             // 控制台输出
		"task4/logs/app.log", // 文件输出
	}
	// 2. 自定义 EncoderConfig，修改 ts 字段的编码器
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",                         // 时间字段名（默认就是 "ts"）
		LevelKey:       "level",                        // 日志级别字段名
		NameKey:        "logger",                       // 日志器名称字段名
		CallerKey:      "caller",                       // 调用者字段名（文件:行号）
		FunctionKey:    zapcore.OmitKey,                // 省略函数名字段
		MessageKey:     "msg",                          // 日志消息字段名
		StacktraceKey:  "stacktrace",                   // 堆栈跟踪字段名
		LineEnding:     zapcore.DefaultLineEnding,      // 行结束符
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 级别小写（debug/info/warn/error）
		EncodeTime:     customTimeEncoder,              // 自定义时间编码器（核心）
		EncodeDuration: zapcore.SecondsDurationEncoder, // 耗时格式（秒）
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 调用者格式（短路径，如 main.go:20）
	}
	return cfg.Build()
}

// 自定义带时区的时间编码器（UTC 时区）
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	utcTime := t.UTC() // 转换为 UTC 时区
	enc.AppendString(utcTime.Format("2006-01-02 15:04:05.000"))
}

// 初始化数据库
func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Cfg.DBConfig.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(config.Cfg.DBConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Cfg.DBConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.Cfg.DBConfig.ConnMaxLifetime)

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
