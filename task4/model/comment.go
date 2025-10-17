package model

import "gorm.io/gorm"

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"`
	PostID  uint   `gorm:"not null" json:"postId"`                       // 外键：文章ID
	UserID  uint   `gorm:"not null" json:"userId"`                       // 外键：评论者ID
	User    User   `gorm:"foreignKey:UserID" json:"commenter,omitempty"` // 关联评论者
	Post    Post   `gorm:"foreignKey:PostID" json:"post,omitempty"`      // 关联文章
}
