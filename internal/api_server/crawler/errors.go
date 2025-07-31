package crawler

// EOFError 定义了逐步更新项目时的结束错误
type EOFError struct{}

// Error 实现了 error 接口
func (e *EOFError) Error() string {
	return "已到达数据末尾"
}
