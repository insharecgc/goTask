package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gotask/task4/service"
	"gotask/task4/util"
)

// CommentHandler 评论控制器
type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID uint `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// Create 创建评论（需登录）
func (h *CommentHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, util.ErrorParam(err.Error()))
		return
	}

	comment, err := h.commentService.Create(req.Content, req.PostID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorWithMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success(comment))
}

// ListByPost 获取文章的所有评论（公开）
func (h *CommentHandler) ListByPost(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorParam(err.Error()))
		return
	}

	comments, err := h.commentService.ListByPostID(uint(postID))
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorWithMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success(comments))
}
