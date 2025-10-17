package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gotask/task4/service"
	"gotask/task4/util"
)

// UserHandler 用户控制器
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, util.ErrorParam(err.Error()))
		return
	}

	if err := h.userService.Register(req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusOK, util.ErrorWithMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success(nil))
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, util.ErrorParam(err.Error()))
		return
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorWithMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success(token))
}
