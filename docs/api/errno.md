# 针对返回错误码的整理

## 正常响应码

error_code 为 0 是默认响应，也是正常处理时的错误码

## Service 层产生的错误一览

> 4xx HTTP 响应码对应的 errno 位于 **(0, 500]**
> 5xx HTTP 想听吗对应的 errno 位于 **(1000, infinite)**

```go
// 对应 4xx 错误码的错误
const (
 ErrParamsLackage       errCode = iota + 1 // 缺少必要参数
 ErrMoetranAPIFailure                      // 龙译 API 调用失败
 ErrPasswordMismatch                       // 密码不匹配
 ErrInvalidOperator                        // 冒用当前成员身份的操作
 ErrNoSatifiedResults                      // 没有满足条件的结果
 ErrUnacceptedOperation                    // 无法完成的操作
 ErrNoPermission                           // 没有权限
)

// 对应 5xx 错误码的错误
const (
 ErrPassworHashFailure     errno = iota + 1001 // 密码记录失败
 ErrDatabaseFailure                              // 数据提取或记录失败
 ErrTokenGenerationFailure                       // 生成令牌失败
)
```

## Handler 层产生的错误一览（一定是 4xx 错误码）

> 4xx HTTP 响应码对应的 errno 位于 **(500, 1000]**

```go
// 在 handler 层就能检出的错误，一定是 4xx 错误码
const (
 ErrAuthLackage     errno = iota + 501 // 缺少 authorization 信息
 ErrHeaderLackage                        // 缺少必要的请求头信息
 ErrParamsLackage                        // 缺少必要参数
 ErrBadParams                            // 非法的请求参数组合
 ErrUnmatchedUserId                      // 路径参数 user_id 与当前登录用户 ID 不匹配
 ErrUnprocessable                        // 无法处理的请求
)
```
