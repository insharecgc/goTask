package model

import (
	"gorm.io/gorm"
	"gotask/task4/util"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"size:50;unique;not null" json:"username"` // 用户名（唯一）
	Password string `gorm:"size:100;not null" json:"-"`              // 密码（加密存储，不返回给前端）
	Email    string `gorm:"size:100;unique;not null" json:"email"`   // 邮箱（唯一）
}

// BeforeSave 保存前加密密码（钩子函数）
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 仅当密码有变化时加密（避免重复加密）
	if len(u.Password) < 60 { // bcrypt加密后长度固定为60
		hash, err := util.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hash
	}
	return nil
}
