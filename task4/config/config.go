package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// 配置结构体（与 YAML 文件结构对应）
type Config struct {
	ServerConfig ServerConfig `mapstructure:"server"`
	DBConfig     DBConfig     `mapstructure:"database"`
	JWTConfig    JWTConfig    `mapstructure:"jwt"`
	LogLevel     string       `mapstructure:"logLevel"` // 日志级别：debug/info/warn/error
}

// 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// 数据库配置
type DBConfig struct {
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	DbName          string        `mapstructure:"dbName"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
	DSN             string        // 数据库链接字符串（内部使用）
}

// JWT 配置
type JWTConfig struct {
	SecretKey      string        `mapstructure:"secretKey"`
	ExpirationTime string        `mapstructure:"expirationTime"`
	ExpireDuration time.Duration // 解析后的时间（内部使用）
}

// 全局配置实例
var Cfg Config

func Init() error {
	// 读取 YAML 文件内容
	viper.SetConfigFile("task4/config.yaml")
	viper.SetConfigType("yaml")   // 配置文件类型
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&Cfg)
	if err != nil {
		return err
	}

	// 校验必要配置
	if Cfg.DBConfig.User == "" {
		return errors.New("数据库 user 不能为空")
	}
	if Cfg.DBConfig.Password == "" {
		return errors.New("数据库 password 不能为空")
	}
	if Cfg.DBConfig.Host == "" {
		// 为空设置默认值 127.0.0.1
		Cfg.DBConfig.Host = "127.0.0.1"
	}
	if Cfg.DBConfig.Port == "" {
		// 为空设置默认值 3306
		Cfg.DBConfig.Port = "3306"
	}
	if Cfg.DBConfig.DbName == "" {
		return errors.New("数据库 dbName 不能为空")
	}
	if Cfg.DBConfig.MaxOpenConns == 0 {
		// 为空设置默认值 100
		Cfg.DBConfig.MaxOpenConns = 100
	}
	if Cfg.DBConfig.MaxIdleConns == 0 {
		// 为空设置默认值 20
		Cfg.DBConfig.MaxOpenConns = 20
	}
	if Cfg.JWTConfig.SecretKey == "" {
		return errors.New("JWT 密钥不能为空")
	}
	if Cfg.JWTConfig.ExpirationTime == "" {
		// 为空设置默认值 24小时
		Cfg.JWTConfig.ExpirationTime = "24h"
	}
	if Cfg.ServerConfig.Port == "" {
		Cfg.ServerConfig.Port = ":8080"
	}

	Cfg.DBConfig.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Cfg.DBConfig.User,
		Cfg.DBConfig.Password,
		Cfg.DBConfig.Host,
		Cfg.DBConfig.Port,
		Cfg.DBConfig.DbName,
	)

	// 解析 JWT 有效期（字符串转 time.Duration）
	expiration, err := time.ParseDuration(Cfg.JWTConfig.ExpirationTime)
	if err != nil {
		return fmt.Errorf("JWT 有效期格式错误（支持 h/m/s）: %v", err)
	}
	Cfg.JWTConfig.ExpireDuration = expiration

	return nil
}
