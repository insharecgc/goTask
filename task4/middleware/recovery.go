package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gotask/task4/util"
)

// Recovery 捕获panic，统一返回错误响应
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic日志
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				// 返回内部错误
				c.JSON(http.StatusInternalServerError, util.ErrInternalError)
				c.Abort()
			}
		}()
		
		c.Next()

		// 处理主动返回的错误（在handler中通过c.Error()传递）
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			if e, ok := err.Err.(*util.Errno); ok {
				c.JSON(http.StatusOK, e) // 业务错误（前端根据code处理）
			} else {
				logger.Error("handler error", zap.Error(err.Err))
				c.JSON(http.StatusInternalServerError, util.ErrInternalError)
			}
			c.Abort()
		}
	}
}