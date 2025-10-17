package service

import (
	"gorm.io/gorm"
	"gotask/task4/model"
	"gotask/task4/util"
)

// CommentService 评论服务
type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

// Create 创建评论（需验证文章存在）
func (s *CommentService) Create(content string, postID, userID uint) (*model.Comment, error) {
	// 检查文章是否存在
	var post model.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, util.ErrPostNotExist
	}

	// 创建评论
	comment := model.Comment{
		Content: content,
		PostID:  postID,
		UserID:  userID,
	}
	if err := s.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListByPostID 获取文章的所有评论（包含评论者信息）
func (s *CommentService) ListByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	if err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username") // 只返回评论者ID和用户名
	}).Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}