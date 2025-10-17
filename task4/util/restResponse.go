package util

// RestResponse 通用响应结构体
type RestResponse struct {
	Code int         `json:"code"` // 错误码（0表示成功）
	Msg  string      `json:"msg"`  // 消息
	Data interface{} `json:"data"` // 业务数据（成功时返回，失败时为nil）
}

// Success 成功响应（data为具体业务数据）
func Success(data interface{}) *RestResponse {
	return &RestResponse{
		Code: ErrSuccess.Code,
		Msg:  ErrSuccess.Msg,
		Data: data,
	}
}

// Error 错误响应（接收自定义Errno）
func Error(err *Errno) *RestResponse {
	return &RestResponse{
		Code: err.Code,
		Msg:  err.Msg,
		Data: nil,
	}
}

func ErrorParam(msg string) *RestResponse {
	return Error(NewErrno(ErrInvalidParam, msg))
}

// ErrorWithCodeMsg 错误响应（接收自定义错误码和消息）
func ErrorWithCodeMsg(code int, msg string) *RestResponse {
	return &RestResponse{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

// ErrorWithMsg 错误响应（接收错误消息，不关心错误码）
func ErrorWithMsg(msg string) *RestResponse {
	return &RestResponse{
		Code: -1,
		Msg:  msg,
		Data: nil,
	}
}
