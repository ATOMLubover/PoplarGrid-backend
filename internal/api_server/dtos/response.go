package dtos

// MemberLabor 定义了项目成员的分工信息
type MemberLabor struct {
	UserId     uint   `json:"user_id"`     // 用户 ID
	Nickname   string `json:"nickname"`    // 昵称
	LaborRole  uint   `json:"labor_role"`  // 成员在项目中的角色，使用位掩码表示
	JoinedTime string `json:"joined_time"` // 加入的
}

// ProjectBasic 定义了获取项目的基本信息
type ProjectBasic struct {
	Id        uint   `json:"id"`         // 项目 ID
	Title     string `json:"title"`      // 项目名称
	LegacyId  uint   `json:"legacy_id"`  // 历史遗留序号
	MoetranId string `json:"moetran_id"` // 龙译 ID

	WorksetId    uint `json:"workset_id"`    // 所属作品集 ID
	WorksetIndex uint `json:"workset_index"` // 作品集内的序号

	Status      ProjectOverallStatus `json:"status"`       // 项目状态，使用位掩码表示
	IsPublished bool                 `json:"is_published"` // 是否已发布

	AllowAutoJoin bool `json:"allow_auto_join"` // 是否允许自动加入
	IsHidden      bool `json:"is_hidden"`       // 是否是隐藏项目

	Labors *[]MemberLabor `json:"labors,omitempty"` // 补充的成员分工信息
}

// ProjectDetail 定义了获取项目的详细信息
type ProjectDetail struct {
	ProjectBasic `json:",inline"` // 嵌入 ProjectBasic 的字段
	Description  string           `json:"description"` // 项目描述
	CreatedAt    string           `json:"created_at"`  // 创建时间
	UpdatedAt    string           `json:"updated_at"`  // 更新时间
}

// ProjectStats 定义了项目整体的一些统计情况
type ProjectStats struct {
	WorksetId uint `json:"workset_id"` // 对应作品集 ID

	TotalCount     int `json:"total_count"`     // 总项目数量
	PublishedCount int `json:"published_count"` // 已发布项目数量

	NotTranslatingCount int `json:"not_translating_count"` // 未开始翻译的项目数量
	TranslatingCount    int `json:"translating_count"`     // 正在翻译的项目数量
	TranslatedCount     int `json:"translated_count"`      // 已翻译的项目数量

	NotProovingCount int `json:"not_prooving_count"` // 未开始校对的项目数量
	ProovingCount    int `json:"prooving_count"`     // 正在校对的项目数量
	ProovedCount     int `json:"prooved_count"`      // 已校对的项目数量

	NotLetteringCount int `json:"not_lettering_count"` // 未开始嵌字的项目数量
	LetteringCount    int `json:"lettering_count"`     // 正在嵌字的项目数量
	LetteredCount     int `json:"lettered_count"`      // 已嵌字的项目数量

	NotReviewingCount int `json:"not_reviewing_count"` // 未开始审核的项目数量
	ReviewingCount    int `json:"reviewing_count"`     // 正在审核的项目数量
	ReviewedCount     int `json:"reviewed_count"`      // 已审核的项目数量
}

// WorksetBasic 定义了作品集信息
type WorksetBasic struct {
	Id     uint   `json:"id"`      // 作品集 ID
	Name   string `json:"name"`    // 作品集名称
	TeamId uint   `json:"team_id"` // 所属团队 ID
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
	UserBasic     `json:",inline"`
	Email         string `json:"email"`           // 邮箱
	PoplarIsAdmin bool   `json:"poplar_is_admin"` // 是否是 panel 管理员
	Remark        string `json:"remark"`          // 备注
	QqNumber      string `json:"qq_number"`       // QQ 号
	LastActive    string `json:"last_active"`     // 上次活跃时间
}

// MemberBasic 定义了项目成员的基本信息
type MemberBasic struct {
	Id       uint   `json:"id"`       // 成员 ID
	UserId   uint   `json:"user_id"`  // 实际上的用户 ID
	TeamId   uint   `json:"team_id"`  // 所属团队 ID
	Nickname string `json:"nickname"` // 昵称
	Role     uint   `json:"role"`     // 成员在组内的职责
}

// InnerProject 定义了内嵌携带的简单项目信息
type InnerProject struct {
	Id           uint                 `json:"project_id"`    // 所属项目 ID
	Title        string               `json:"project_title"` // 项目标题
	WorksetId    uint                 `json:"workset_id"`    // 所属作品
	WorksetIndex uint                 `json:"workset_index"` // 作品集内的序号
	Status       ProjectOverallStatus `json:"status"`        // 项目状态，使用位
}

// InvitationBasic 定义了邀请的基本信息
type InvitationBasic struct {
	Id        uint `json:"id"`         // 邀请 ID
	InviterId uint `json:"inviter_id"` // 邀请者 ID
	// InviterNickname string `json:"inviter_nickname"` // 邀请者昵称
	InviteeId       uint   `json:"invitee_id"`       // 被邀请者 ID
	InviteeNickname string `json:"invitee_nickname"` // 被邀请者昵称
	InviteRole      uint   `json:"invite_role"`      // 邀请的角色，使用位掩码表示
	Status          int    `json:"status"`           // 邀请状态，0 pending, 1 accepted, 2 rejected

	Project InnerProject `json:"project"` // 所属项目的基本信息
}

// ApplicationBasic 定义了申请的基本信息
type ApplicationBasic struct {
	Id          uint   `json:"id"`           // 申请 ID
	ApplicantId uint   `json:"applicant_id"` // 申请者 ID
	Nickname    string `json:"nickname"`     // 申请者昵称
	Role        uint   `json:"role"`         // 申请的角色，使用位掩码表示
	Status      int    `json:"status"`       // 申请状态，0 pending, 1 accepted, 2 rejected

	Project InnerProject `json:"project"` // 所属项目的基本信息
}

// LaborDivision 定义了项目成员的分工信息
type LaborDivision struct {
	MemberId uint   `json:"member_id"` // 成员 ID
	Nickname string `json:"nickname"`  // 昵称
	Role     uint   `json:"role"`      // 成员在项目中的角色，使用位掩码表示
}

// ProjectCreatedInfo 定义了创建项目的响应结构体
type ProjectCreatedInfo struct {
	Message   string `json:"message"`    // 响应消息
	ProjectId uint   `json:"project_id"` // 创建的项目 ID
	MoetranId string `json:"moetran_id"` // 龙译项目 ID
}
