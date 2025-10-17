package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gotask/task4/service"
	"gotask/task4/util"
)

// PostHandler 文章控制器
type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required"`
}

// Create 创建文章（需登录）
func (h *PostHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID") // 从JWT中间件获取用户ID

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, util.ErrorParam(err.Error()))
		return
	}

	post, err := h.postService.Create(req.Title, req.Content, userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, util.Error(util.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, util.Success(post))
}

// List 分页查询文章列表（公开接口）
func (h *PostHandler) Page(c *gin.Context) {
	// 绑定分页参数（默认page=1，pageSize=10）
	var param util.PageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		// 参数验证失败时，设置默认值
		param = util.PageParam{Page: 1, PageSize: 10}
	}

	// 调用服务层分页查询
	pageResult, err := h.postService.Page(param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, util.Error(util.ErrInternalError))
		return
	}

	// 返回统一响应
	c.JSON(http.StatusOK, util.Success(pageResult))
}

// Get 查看用户文章列表（公开，不需要权限）
func (h *PostHandler) ListByUserId(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, util.Error(util.ErrInvalidParam))
		return
	}
	posts, err := h.postService.ListByUserId(uint(userId))
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorWithMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success(posts))
}

// Get 查看文章（公开，不需要权限）
func (h *PostHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, util.Error(util.ErrInvalidParam))
		return
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorWithMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success(post))
}

// UpdatePostRequest 修改文章请求
type UpdatePostRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Title   string `json:"title" binding:"omitempty,min=1,max=200"`
	Content string `json:"content" binding:"omitempty"`
}

// Update 修改文章（仅作者）
func (h *PostHandler) Update(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, util.ErrorParam(err.Error()))
		return
	}

	if err := h.postService.Update(req.ID, req.Title, req.Content, userID.(uint)); err != nil {
		c.JSON(http.StatusOK, util.ErrorWithMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success(nil))
}
