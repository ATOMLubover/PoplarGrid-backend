package dto

// CreateProjectRequest 定义了创建项目的请求体
type CreateProjectRequest struct {
	Title       string   `json:"title" validate:"required"`      // 项目标题，必填，最大长度 128
	WorksetId   uint     `json:"workset_id" validate:"required"` // 作品集 ID，必填
	Description string   `json:"description"`                    // 项目描述，最大长度 512
	Tags        []string `json:"tags"`                           // 标签列表，最多 5 个标签
}
