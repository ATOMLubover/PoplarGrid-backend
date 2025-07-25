package models

// -------------
// error 定义了一些通用的错误类型
// -------------

// InvalidParameterError 定义了无效参数的错误
type InvalidParameterError struct{}

func (e *InvalidParameterError) Error() string {
	return "参数无效"
}

// LackOfRequiredFieldError 定义了缺少必需字段的错误
type LackOfRequiredFieldError struct{}

func (e *LackOfRequiredFieldError) Error() string {
	return "缺少必需字段"
}

// LackOfPrimaryKeyError 定义了 update/delete 时缺少指定主键的错误
type LackOfPrimaryKeyError struct{}

func (e *LackOfPrimaryKeyError) Error() string {
	return "缺少指定主键"
}
