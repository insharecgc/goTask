package model

import "gorm.io/gorm"

// Post 文章模型
type Post struct {
	gorm.Model
	Title   string `gorm:"size:200;not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	UserID  uint   `gorm:"not null" json:"userId"`                    // 外键：作者ID
	User    User   `gorm:"foreignKey:UserID" json:"author,omitempty"` // 关联作者（查询时返回）
}
