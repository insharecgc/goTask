package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gotask/task4/util"
)

// JWTAuth 验证JWT token，通过后将userID存入上下文
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取token（格式：Bearer <token>）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, util.ErrNoPermission)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, util.ErrNoPermission)
			c.Abort()
			return
		}

		// 解析token
		data, err := util.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, util.ErrNoPermission)
			c.Abort()
			return
		}

		// 将userID, userName存入上下文，供后续 handler 使用
		c.Set("userID", data["userId"])
		c.Set("userName", data["userName"])
		c.Next()
	}
}
