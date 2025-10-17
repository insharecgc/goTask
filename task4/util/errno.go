package util

import "fmt"

// Errno 自定义错误类型（包含错误码和消息）
type Errno struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Errno) Error() string {
	return e.Msg
}

// 定义常见错误码
var (
	ErrSuccess       = &Errno{Code: 200, Msg: "success"}
	ErrInvalidParam  = &Errno{Code: 1001, Msg: "请求参数错误"}
	ErrUserExist     = &Errno{Code: 2001, Msg: "用户已存在"}
	ErrEmailExist    = &Errno{Code: 2002, Msg: "邮箱已存在"}
	ErrUserNotExist  = &Errno{Code: 2003, Msg: "用户不存在"}
	ErrInvalidPass   = &Errno{Code: 2003, Msg: "密码错误"}
	ErrPostNotExist  = &Errno{Code: 3001, Msg: "文章不存在"}
	ErrNoPermission  = &Errno{Code: 4001, Msg: "没有权限"}
	ErrInternalError = &Errno{Code: 5001, Msg: "服务器内部错误"}
)

// NewErrno 格式化错误消息
func NewErrno(err *Errno, format string, v ...interface{}) *Errno {
	return &Errno{
		Code: err.Code,
		Msg:  fmt.Sprintf(err.Msg+": "+format, v...),
	}
}
