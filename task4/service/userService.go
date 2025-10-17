package service

import (
	"gorm.io/gorm"
	"gotask/task4/model"
	"gotask/task4/util"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB // 数据库增删改查
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Register 用户注册
func (s *UserService) Register(username, email, password string) error {
	// 检查用户是否已存在
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err == nil {
		return util.NewErrno(util.ErrUserExist, "username=%s", username)
	}
	if err := s.db.Where("email = ?", email).First(&user).Error; err == nil {
		return util.NewErrno(util.ErrEmailExist, "email=%s", email)
	}

	// 创建用户（密码会在BeforeSave钩子中加密）
	user = model.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	return s.db.Create(&user).Error
}

// Login 用户登录（返回token）
func (s *UserService) Login(username, password string) (string, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", util.ErrUserNotExist
	}

	// 验证密码
	if !util.CheckPassword(user.Password, password) {
		return "", util.ErrInvalidPass
	}

	// 生成JWT token
	return util.GenerateToken(user.ID, user.Username)
}
