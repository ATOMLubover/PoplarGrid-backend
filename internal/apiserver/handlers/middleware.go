package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

type IrisMiddleware func(ctx iris.Context)

// NewCorsMiddleware 生成一个跨域中间件
// 其作用于 handler 之前
// 可以注入要放行的域名、方法等参数
func NewCorsMiddleware(
	allowOrigins []string,
	allowMethods []string,
	allowHeaders []string,
	allowCredentials bool,
	maxAge time.Duration,
) IrisMiddleware {
	if len(allowOrigins) == 0 {
		// 如果没有设置允许的源，直接返回
		return nil
	}

	// 将 origin 组成一个 set
	allowOriginSet := make(map[string]struct{}, len(allowOrigins))
	for _, origin := range allowOrigins {
		if origin != "" {
			allowOriginSet[origin] = struct{}{}
		}
	}

	return func(ctx iris.Context) {
		// 先清除之前所有的 CORS 相关头部，防止重复设置
		ctx.Header("Access-Control-Allow-Origin", "")
		ctx.Header("Access-Control-Allow-Methods", "")
		ctx.Header("Access-Control-Allow-Headers", "")
		ctx.Header("Access-Control-Allow-Credentials", "")
		ctx.Header("Access-Control-Max-Age", "")

		// 动态设置 Access-Control-Allow-Origin
		{
			// 由于实际上 Access-Control-Allow-Origin 只能设置一个源
			// 所以使用遍历查找 origin 来动态确定其值（完全匹配优先于通配符）
			origin := ctx.GetHeader("Origin")

			_, matchExisting := allowOriginSet[origin]
			_, wildcardExisting := allowOriginSet["*"]

			switch {
			case matchExisting:
				// 优先查找是否有严格匹配的 origin
				ctx.Header("Access-Control-Allow-Origin", origin)
			case wildcardExisting:
				// 如果没有严格匹配的 origin，则尝试使用通配符
				ctx.Header("Access-Control-Allow-Origin", "*")
			}
		}

		// 设置 Access-Control-Allow-Methods
		if len(allowMethods) > 0 {
			ctx.Header("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
		}

		// 动态设置 Access-Control-Allow-Headers (针对预检请求)
		if ctx.Method() == iris.MethodOptions {
			requestHeaders := ctx.GetHeader("Access-Control-Request-Headers")
			switch {
			case requestHeaders != "":
				// 如果请求头中有 Access-Control-Request-Headers，则直接使用它
				ctx.Header("Access-Control-Allow-Headers", requestHeaders)
			case len(allowHeaders) > 0:
				// 如果没有指定请求头，则使用预设的允许头部
				ctx.Header("Access-Control-Allow-Headers", strings.Join(allowHeaders, ", "))
			}
		}

		// 设置 Access-Control-Allow-Credentials
		if allowCredentials {
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}

		// 设置 Access-Control-Max-Age (针对预检请求)
		if maxAge > 0 {
			ctx.Header("Access-Control-Max-Age", fmt.Sprintf("%.0f", maxAge.Seconds()))
		}

		if ctx.Method() == iris.MethodOptions {
			// 如果是预检请求，直接返回 200 OK
			ctx.StatusCode(iris.StatusOK)
			return
		}

		// 继续处理后续的请求
		ctx.Next()
	}
}

// NewUserInfoExtractMiddleware 生成一个解析 Gateway Server 解析信息到 ctx 的中间件
// 其作用于 handler 之前
func NewUserInfoExtractMiddleware() IrisMiddleware {
	return func(ctx iris.Context) {
		// 首先解析 Client-Request-User-Id
		if userIdStr := ctx.GetHeader("Client-Request-User-Id"); userIdStr != "" {
			userId, err := strconv.ParseUint(userIdStr, 10, 64)
			if err != nil {
				ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
					"error": "无法解析的 user_id 头",
				})
				return
			}

			ctx.Values().Set("user_id", userId)
		}
	}
}

// NewCheckUserIdMiddleware 生成一个检查 user_id 的中间件
func NewCheckUserIdMiddleware() IrisMiddleware {
	return func(ctx iris.Context) {
		// 从上下文中获取 user_id
		userId, err := ctx.Values().GetUint("user_id")
		if err != nil || userId <= 0 {
			ctx.StopWithJSON(iris.StatusBadRequest, ErrorResponse{
				Error: "无法获取有效的 user_id",
			})
			return
		}

		// 从路径参数获取 user_id 并进行验证
		pathUserId, err := ctx.Params().GetUint("id")
		if err != nil || pathUserId != userId {
			ctx.StopWithJSON(iris.StatusBadRequest, ErrorResponse{
				Error: "无效的 user_id，必须与当前登录用户 ID 匹配",
			})
		}

		// 继续处理请求
		ctx.Next()
	}
}
