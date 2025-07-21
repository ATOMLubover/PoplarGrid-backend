package dtos

// ProjectStatusQueryParams 定义了项目搜索的参数
// 设置为 nil 意味着不需要查询
type ProjectStatusQueryParams struct {
	TranslateStatus *uint // 翻译状态
	ProofStatus     *uint // 校对状态
	LetterStatus    *uint // 嵌字状态
	ReviewStatus    *uint // 审核状态
	PublishStatus   *bool // 发布状态
}

// ProjectSearchParams 定义了项目搜索的参数
type ProjectSearchParams struct {
	WorksetId *uint                     // 作品集 ID
	UserId    *uint                     // 用户 ID
	Sort      int                       // 排序方式，0：按 ID 倒序，1：按 updated_at 倒序
	Status    *ProjectStatusQueryParams // 项目状态查询参数
}
