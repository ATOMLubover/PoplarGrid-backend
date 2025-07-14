package dtos

// ProjectBasic 定义了获取项目的进本信息
type ProjectBasic struct {
	Id       uint   `json:"id"`        // 项目 ID
	Title    string `json:"title"`     // 项目名称
	LegacyId uint   `json:"legacy_id"` // 历史遗留序号

	WorksetId    uint `json:"workset_id"`    // 所属作品集 ID
	WorksetIndex uint `json:"workset_index"` // 作品集内的序号

	Status      ProjectOverallStatus `json:"status"`       // 项目状态，使用位掩码表示
	IsPublished bool                 `json:"is_published"` // 是否已发布

	AllowAutoJoin bool `json:"allow_auto_join"` // 是否允许自动加入
	IsHidden      bool `json:"is_hidden"`       // 是否是隐藏项目
}

// ProjectDetail 定义了获取项目的详细信息
type ProjectDetail struct {
	ProjectBasic
	Description string `json:"description"` // 项目描述
	CreatedAt   string `json:"created_at"`  // 创建时间
	UpdatedAt   string `json:"updated_at"`  // 更新时间
}

// MyProjectBasic 定义了获取用户参与的项目的基本信息
type MyProjectBasic struct {
	ProjectBasic
	MemberId uint `json:"member_id"` // 成员 ID
	Role     uint `json:"role"`      // 成员在项目中的角色，使用掩码计算多重身份
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

// WorksetBasic 定义了作品集信息
type WorksetBasic struct {
	Id        uint   `json:"id"`         // 作品集 ID
	Name      string `json:"name"`       // 作品集名称
	MoetranId string `json:"moetran_id"` // Moetran ID
	TeamId    uint   `json:"team_id"`    // 所属团队 ID
}

// TeamBasic 定义了团队信息
type TeamBasic struct {
	Id   uint   `json:"id"`   // 团队 ID
	Name string `json:"name"` // 团队名称
}

// UserBasic 定义了用户的基本信息
type UserBasic struct {
	Id       uint   `json:"id"`       // 用户 ID
	Nickname string `json:"nickname"` // 昵称
}

// UserDetail 定义了用户的详细信息
type UserDetail struct {
	UserBasic
	Email         string `json:"email"`           // 邮箱
	PoplarIsAdmin bool   `json:"poplar_is_admin"` // 是否是 panel 管理员
	Remark        string `json:"remark"`          // 备注
	QqNumber      string `json:"qq_number"`       // QQ 号
	LastActive    string `json:"last_active"`     // 上次活跃时间
}

// MemberBasic 定义了项目成员的基本信息
type MemberBasic struct {
	Id     uint `json:"id"`      // 成员 ID
	UserId uint `json:"user_id"` // 实际上的用户 ID
	TeamId uint `json:"team_id"` // 所属团队 ID
	Role   uint `json:"role"`    // 成员在组内的职责
}

// InvitationBasic 定义了邀请的基本信息
type InvitationBasic struct {
	Id         uint `json:"id"`          // 邀请 ID
	InviterId  uint `json:"inviter_id"`  // 邀请者 ID
	InviteeId  uint `json:"invitee_id"`  // 被邀请者 ID
	ProjectId  uint `json:"project_id"`  // 所属项目 ID
	InviteRole uint `json:"invite_role"` // 邀请的角色，使用位掩码表示
}
