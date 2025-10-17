package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gotask/task4/config"
)

// JWTClaims JWT载荷
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Cfg.JWTConfig.ExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTConfig.SecretKey))
}

// ParseToken 解析JWT token，返回用户信息
func ParseToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JWTConfig.SecretKey), nil
		},
	)
	retData := make(map[string]interface{})
	if err != nil {
		return retData, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		retData["userId"] = claims.UserID
		retData["username"] = claims.Username
		return retData, nil
	}
	return retData, errors.New("invalid token")
}
