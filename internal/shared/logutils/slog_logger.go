package logutils

import (
	"log/slog"
	"os"
)

// 构造 logger
// 直接使用标准库的 slog 以简化依赖
func NewLogger(opts *slog.HandlerOptions) *slog.Logger {
	if opts == nil {
		// 如果选项不为空，则使用 opts 构造 text handler
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	// 否则构造默认 text handler
	// 默认状态使用 debug 级别输出日志，且日期格式为 2006-01-02 15:04:05
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				return slog.String(slog.TimeKey, t.Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}))
}
