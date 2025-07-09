package dto

// 成员的简要信息
type MemberBriefInfo struct {
	MemberId uint   `json:"member_id"`
	Nickname string `json:"nickname"`
}

// 成员的完整信息
type MemberFullInfo struct {
	MemberId  uint   `json:"member_id"`
	MoetranId string `json:"moetran_id"`

	Nickname string `json:"nickname"`
	Email    string `json:"email"`

	PoplarIsAdmin bool `json:"poplar_is_admin"`
	Labors        uint `json:"labors"`

	Remark     string `json:"remark"`
	LastActive string `json:"last_active"` // 格式为 "2006-01-02 15:04:05"
}

// 项目人员分工信息
type ProjectMemberLabor struct {
	MemberName string `json:"member_name"` // 成员昵称
	LaborRole  uint   `json:"labor_role"`  // 分工角色，使用位掩码存储
}

// 项目的完整信息
type ProjectFullInfo struct {
	ProjectId uint `json:"project_id"`

	Title string `json:"title"`

	Team    TeamFullInfo    `json:"team"`    // 汉化组信息
	Workset WorksetFullInfo `json:"workset"` // 作品集信息
	Work    WorkFullInfo    `json:"work"`    // 作品信息

	Status  uint  `json:"status"`
	Urgency int16 `json:"urgency"`

	LegacyId uint `json:"legacy_id"`

	LaborDivision []ProjectMemberLabor `json:"labor_division"` // 分工成员列表
}

// 汉化组完整信息
type TeamFullInfo struct {
	TeamId    uint   `json:"team_id"`
	MoetranId string `json:"moetran_id"` // 尨译 ID

	TeamName string `json:"team_name"`
}

// 作品集的完整信息
type WorksetFullInfo struct {
	WorksetId uint   `json:"workset_id"`
	MoetranId string `json:"moetran_id"` // 尨译 ID

	Title string `json:"title"`
}

// 作品的完整信息
type WorkFullInfo struct {
	WorkId    uint   `json:"work_id"`
	MoetranId string `json:"moetran_id"` // 尨译 ID

	Title       string `json:"title"`
	Description string `json:"description"`

	Tags []string `json:"tags"` // 标签列表
}
