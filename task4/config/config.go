package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// 配置结构体（与 YAML 文件结构对应）
type Config struct {
	ServerConfig ServerConfig `yaml:"server"`
	DBConfig     DBConfig     `yaml:"database"`
	JWTConfig    JWTConfig    `yaml:"jwt"`
}

// 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
}

// 数据库配置
type DBConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         string   `yaml:"port"`
	DbName       string `yaml:"db_name"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	DSN          string // 数据库链接字符串（内部使用）
}

// JWT 配置
type JWTConfig struct {
	SecretKey      string        `yaml:"secret_key"`
	ExpirationTime string        `yaml:"expiration_time"`
	expiration     time.Duration // 解析后的时间（内部使用）
}

func InitConfig() (*Config, error) {
	// 读取 YAML 文件内容
	data, err := os.ReadFile("../config.yaml")
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 YAML 到结构体
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 校验必要配置
	if cfg.DBConfig.User == "" {
		return nil, errors.New("数据库 user 不能为空")
	}
	if cfg.DBConfig.Password == "" {
		return nil, errors.New("数据库 password 不能为空")
	}
	if cfg.DBConfig.Host == "" {
		// 为空设置默认值 127.0.0.1
		cfg.DBConfig.Host = "127.0.0.1"
	}
	if cfg.DBConfig.Port == 0 {
		// 为空设置默认值 3306
		cfg.DBConfig.Port = 3306
	}
	if cfg.DBConfig.DbName == "" {
		return nil, errors.New("数据库 dbName 不能为空")
	}
	if cfg.DBConfig.MaxOpenConns == 0 {
		// 为空设置默认值 100
		cfg.DBConfig.MaxOpenConns = 100
	}
	if cfg.DBConfig.MaxIdleConns == 0 {
		// 为空设置默认值 20
		cfg.DBConfig.MaxOpenConns = 20
	}
	if cfg.JWTConfig.SecretKey == "" {
		return nil, errors.New("JWT 密钥不能为空")
	}
	if cfg.JWTConfig.ExpirationTime == "" {
		// 为空设置默认值 24小时
		cfg.JWTConfig.ExpirationTime = "24h"
	}

	cfg.DBConfig.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DbName
	)

	// 解析 JWT 有效期（字符串转 time.Duration）
	expiration, err := time.ParseDuration(cfg.JWTConfig.ExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("JWT 有效期格式错误（支持 h/m/s）: %v", err)
	}
	cfg.JWTConfig.expiration = expiration

	return &cfg, nil
}

// 获取解析后的有效期（转换为 time.Duration）
func (j *JWTConfig) Expiration() time.Duration {
	return j.expiration
}
