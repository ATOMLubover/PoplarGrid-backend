package services

// Err 定义了服务错误的接口
type Err interface {
	// ErrorCode 返回错误码
	ErrorCode() int
	// Error 返回错误信息
	Error() string
}

// errno 定义了服务错误的类型
type errno int

func (e errno) ErrorCode() int {
	return int(e)
}

func (e errno) Error() string {
	switch e {
	case ErrParamsLackage:
		return "缺少必要参数"
	case ErrMoetranAPIFailure:
		return "龙译 API 调用失败"
	case ErrPasswordMismatch:
		return "密码不匹配"
	case ErrPassworHashFailure:
		return "密码转换失败"
	case ErrDatabaseFailure:
		return "数据库操作失败"
	case ErrTokenGenerationFailure:
		return "生成令牌失败"
	case ErrInvalidOperator:
		return "冒用当前成员身份的操作"
	case ErrNoSatifiedResults:
		return "没有满足条件的结果"
	case ErrUnacceptedOperation:
		return "无法完成的操作"
	case ErrNoPermission:
		return "没有权限"
	default:
		return "未知错误"
	}
}

// srvError 实现了 Error 接口
// 提供额外的错误信息
type srvError struct {
	code    errno
	message string
}

// newSrvError 创建一个新的服务错误
func newSrvError(code errno, message string) Err {
	e := &srvError{
		code: code,
	}

	if message != "" {
		e.message = code.Error() + ": " + message
	}

	return e
}

func (e *srvError) ErrorCode() int {
	return int(e.code)
}

func (e *srvError) Error() string {
	return e.message
}

// IsServerError 函数判断是否是服务端错误
func IsServerError(e Err) bool {
	return e.ErrorCode() > 1000
}

// IsParamError 函数判断是否是客户端参数错误
func IsParamError(e Err) bool {
	return e.ErrorCode() > 0 && e.ErrorCode() <= 500
}

// 对应 4xx 错误码的错误
const (
	ErrParamsLackage       errno = iota + 1 // 缺少必要参数
	ErrMoetranAPIFailure                    // 龙译 API 调用失败
	ErrPasswordMismatch                     // 密码不匹配
	ErrInvalidOperator                      // 冒用当前成员身份的操作
	ErrNoSatifiedResults                    // 没有满足条件的结果
	ErrUnacceptedOperation                  // 无法完成的操作
	ErrNoPermission                         // 没有权限
)

// 对应 5xx 错误码的错误
const (
	ErrPassworHashFailure     errno = iota + 1001 // 密码记录失败
	ErrDatabaseFailure                            // 数据提取或记录失败
	ErrTokenGenerationFailure                     // 生成令牌失败
)
