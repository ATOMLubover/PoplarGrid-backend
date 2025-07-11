package dto

// ProjectBasic 定义了获取项目的进本信息
type ProjectBasic struct {
	Id        uint   `json:"id"`         // 项目 ID
	Title     string `json:"title"`      // 项目名称
	LegacyId  uint   `json:"legacy_id"`  // 历史遗留序号
	WorksetId uint   `json:"workset_id"` // 所属作品集 ID

	OnTranslating bool `json:"on_translating"` // 是否正在翻译
	IsTranslated  bool `json:"is_translated"`  // 是否已翻译
	OnProoving    bool `json:"on_prooving"`    // 是否正在校对
	IsProofed     bool `json:"is_proofed"`     // 是否已校对
	OnLettering   bool `json:"on_lettering"`   // 是否正在嵌字
	IsLettered    bool `json:"is_lettered"`    // 是否已嵌字
	OnReviewing   bool `json:"on_reviewing"`   // 是否正在审核
	IsReviewed    bool `json:"is_reviewed"`    // 是否已审核
	IsPublished   bool `json:"is_published"`   // 是否已发布
}

// ProjectDetail 定义了获取项目的详细信息
type ProjectDetail struct {
	ProjectBasic
	MoetranId   string   `json:"moetran_id"`  // Moetran ID
	Description string   `json:"description"` // 项目描述
	Status      uint     `json:"status"`      // 项目状态
	Tags        []string `json:"tags"`        // 标签列表
	CreatedAt   string   `json:"created_at"`  // 创建时间
	UpdatedAt   string   `json:"updated_at"`  // 更新时间
}

// ProjectStats 定义了项目整体的一些统计情况
type ProjectStats struct {
	PublishedCount   int `json:"published_count"`   // 已发布项目数量
	TranslatingCount int `json:"translating_count"` // 正在翻译的项目数量
	TranslatedCount  int `json:"translated_count"`  // 已翻译的项目数量
	ProovingCount    int `json:"prooving_count"`    // 正在校对的项目数量
	ProovedCount     int `json:"prooved_count"`     // 已校对的项目数量
	LetteringCount   int `json:"lettering_count"`   // 正在嵌字的项目数量
	LetteredCount    int `json:"lettered_count"`    // 已嵌字的项目数量
	ReviewingCount   int `json:"reviewing_count"`   // 正在审核的项目数量
	ReviewedCount    int `json:"reviewed_count"`    // 已审核的项目数量
	TotalCount       int `json:"total_count"`       // 总项目数量
}

// TagBasic 定义了标签信息
type TagBasic struct {
	Id          uint   `json:"id"`          // 标签 ID
	Name        string `json:"name"`        // 标签名称
	Description string `json:"description"` // 标签描述
}

// WorksetBasic 定义了作品集信息
type WorksetBasic struct {
	Id        uint   `json:"id"`         // 作品集 ID
	Name      string `json:"name"`       // 作品集名称
	MoetranId string `json:"moetran_id"` // Moetran ID
	TeamId    uint   `json:"team_id"`    // 所属团队 ID
}

// TeamBasic 定义了团队信息
type TeamBasic struct {
	Id        uint   `json:"id"`         // 团队 ID
	Name      string `json:"name"`       // 团队名称
	MoetranId string `json:"moetran_id"` // Moetran ID
}
