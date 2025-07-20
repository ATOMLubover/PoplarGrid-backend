package handlers

// ErrorResponse 统一定义发生错误时的 JSON 响应格式
type ErrorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"` // 可选，在 service 层发生错误时提供详细信息
}

// SuccessResponse 统一定义成功响应的 JSON 格式
type SuccessResponse struct {
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"` // 可选，提供额外的成功信息
}
