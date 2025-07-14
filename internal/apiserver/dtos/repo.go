package dtos

// ProjectStatusQueryParams 定义了项目搜索的参数
// 设置为 nil 意味着不需要查询
type ProjectStatusQueryParams struct {
	TranslateStatus *uint `json:"translating"` // 翻译状态
	ProofStatus     *uint `json:"prooving"`    // 校对状态
	LetterStatus    *uint `json:"lettering"`   // 嵌字状态
	ReviewStatus    *uint `json:"reviewing"`   // 审核状态
}
