package service

import (
	"gorm.io/gorm"
	"gotask/task4/model"
	"gotask/task4/util"
)

// PostService 文章服务
type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

// Create 创建文章
func (s *PostService) Create(title, content string, userID uint) (*model.Post, error) {
	post := model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
	if err := s.db.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// List 分页查询文章列表（按创建时间倒序）
func (s *PostService) Page(page, pageSize int) (*util.PageResult, error) {
	var (
		posts []model.Post
		total int64
	)

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 1. 查询总条数
	if err := s.db.Model(&model.Post{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 2. 分页查询文章（预加载作者信息，只返回用户名）
	if err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username")
	}).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, err
	}

	// 3. 组装分页结果
	return util.CalcPageResult(posts, total, page, pageSize), nil
}

func (s *PostService) ListByUserId(userID uint) ([]model.Post, error) {
	var posts []model.Post
	if err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username") // 只返回作者ID和用户名（保护隐私）
	}).Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// GetByID 获取文章（包含作者信息）
func (s *PostService) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username") // 只返回作者ID和用户名（保护隐私）
	}).Where("id = ?", id).First(&post).Error; err != nil {
		return nil, util.ErrPostNotExist
	}
	return &post, nil
}

// Update 修改文章（仅作者可修改）
func (s *PostService) Update(id uint, title, content string, owerUserId uint) error {
	// 检查文章是否存在且属于当前用户
	var post model.Post
	if err := s.db.Where("id = ? AND user_id = ?", id, owerUserId).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return util.ErrNoPermission
		}
		return err
	}

	// 更新文章
	return s.db.Model(&post).Updates(map[string]interface{}{
		"title":   title,
		"content": content,
	}).Error
}