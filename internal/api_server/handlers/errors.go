package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
)

// Err 继承了 services.Err 接口
// 注意，为了保证健壮性，如果发生错误的 API 调用
// 可能会返回错误码为 -1 的错误
type Err services.Err

// errno 实现了 Error 接口
// 其为 handler 层的错误类型
type errno int

// ErrorCode 返回错误码
func (e errno) ErrorCode() int {
	return int(e)
}

// Error 返回错误信息
func (e errno) Error() string {
	switch e {
	case ErrAuthLackage:
		return "缺少认证信息"
	case ErrParamsLackage:
		return "缺少必要参数"
	case ErrHeaderLackage:
		return "缺少必要的请求头信息"
	case ErrBadParams:
		return "非法的请求参数组合"
	case ErrUnmatchedUserId:
		return "路径参数 user_id 与当前登录用户 ID 不匹配"
	default:
		return "未知错误"
	}
}

// IsHandlerError 函数检查是否是 handler 层错误
func IsHandlerError(e Err) bool {
	return e.ErrorCode() > 500 && e.ErrorCode() <= 1000
}

// 在 handler 层就能检出的错误，一定是 4xx 错误码
const (
	ErrAuthLackage     errno = iota + 501 // 缺少 authorization 信息
	ErrHeaderLackage                      // 缺少必要的请求头信息
	ErrParamsLackage                      // 缺少必要参数
	ErrBadParams                          // 非法的请求参数组合
	ErrUnmatchedUserId                    // 路径参数 user_id 与当前登录用户 ID 不匹配
	ErrUnprocessable                      // 无法处理的请求
)

// hdlErr 实现了 Err 接口
// 提供额外的错误信息
type hdlErr struct {
	code    errno
	message string
}

func newHdlErr(code errno, message string) Err {
	e := &hdlErr{
		code: code,
	}

	if message != "" {
		e.message = code.Error() + ": " + message
	}

	return e
}

func (e *hdlErr) ErrorCode() int {
	return int(e.code)
}

func (e *hdlErr) Error() string {
	return e.message
}

// wrapError 将 Error 类型快速转换为错误响应
// 使用 string 作为数据类型，表示没有具体数据返回，方便 Swagger 文档生成
func wrapError(ctx iris.Context, e Err) {
	if e == nil {
		// 如果出现使用 nil 错误调用 wrapError 的情况
		// 返回 500，表示是预期外的操作
		ctx.StatusCode(iris.StatusInternalServerError)

		ctx.Text("预期外的错误")

		return
	}

	// 根据错误类型设置响应状态码
	if services.IsServerError(e) {
		// 如果是服务端错误，则使用 Text 响应
		ctx.WriteString(e.Error())

		return
	}

	// 其余错误都使用 JSON 响应 + 4xx 状态码
	ctx.StatusCode(iris.StatusBadRequest)

	ctx.JSON(&StringFormatResponse{
		ErrorCode: e.ErrorCode(),
		Message:   e.Error(),
	})
}

// wrapSuccess 将成功响应转换为统一格式
func wrapSuccess[T any](ctx iris.Context, data T) {
	ctx.JSON(&FormatResponse[T]{
		Data: data,
	})
}
