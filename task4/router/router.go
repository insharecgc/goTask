package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gotask/task4/handler"
	"gotask/task4/middleware"
)

// Setup 初始化路由
func Setup(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler,
	logger *zap.Logger,
) {
	// 全局中间件
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))

	// 公开路由
	public := r.Group("/api/v1")
	{
		// 用户相关
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)

		// 文章相关（公开访问）
		public.GET("/posts/:id", postHandler.Get)
		public.GET("/posts/list/:userId", postHandler.ListByUserId)
		public.GET("/posts/page", postHandler.Page)

		// 评论相关（公开访问）
		public.GET("/posts/comments/:postId", commentHandler.ListByPost)
	}

	// 需要认证的路由（JWT验证）
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth())
	{
		// 文章相关（需登录）
		auth.POST("/posts", postHandler.Create)
		auth.PUT("/posts/update", postHandler.Update)

		// 评论相关（需登录）
		auth.POST("/posts/addComments", commentHandler.Create)
	}
}