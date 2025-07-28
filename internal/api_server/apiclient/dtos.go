package apiclient

// MoetranAppliCheckType 定义加入项目申请的处理方式
type MoetranAppliCheckType int

// 加入项目申请的处理方式
const (
	APPLI_NON_CHECK   MoetranAppliCheckType = iota // 允许所有人直接加入
	APPLI_ADMIN_CHECK                              // 仅允许管理员通过申请
)

// MoetranAllowApplyType 定义了项目的可申请范围
type MoetranAllowApplyType int

// 项目的可申请范围
const (
	ALLOW_NO_APPLI    MoetranAllowApplyType = iota // 禁止任何形式的申请
	ALLOW_ANY_APPLI                                // 允许任何形式的申请
	ALLOW_MEMBER_ONLY                              // 只允许本组成员申请
)

// MoetranRole 定义了龙译中的角色分工
type MoetranRole string

// 龙译角色的常量定义
const (
	ROLE_ADMIN       MoetranRole = "63d87c24b8bebd75ff934264" // 管理员
	ROLE_SUPERVISOR  MoetranRole = "63d87c24b8bebd75ff934265" // 监理
	ROLE_PROOFREADER MoetranRole = "63d87c24b8bebd75ff934266" // 校对
	ROLE_TRANSLATOR  MoetranRole = "63d87c24b8bebd75ff934267" // 翻译
	ROLE_EMBEDDER    MoetranRole = "63d87c24b8bebd75ff934268" // 嵌字
	ROLE_INTERN      MoetranRole = "63d87c24b8bebd75ff934269" // 实习翻译
)

// 语言代码
const (
	// East Asian Languages
	LangJapanese           string = "ja"    // 日语
	LangSimplifiedChinese  string = "zh-CN" // 简体中文
	LangTraditionalChinese string = "zh-TW" // 繁体中文
	LangKorean             string = "ko"    // 韩语

	// Western European Languages
	LangEnglish string = "en" // 英语
)

// CreateProjectInfo 定义了创建项目的基本信息
type CreateProjectInfo struct {
	MoetranAuth string // 龙译 JWT

	Title       string // 项目标题
	Description string // 项目简介

	MoetranProjSetId string // 龙译项目集 ID
	MoetranTeamId    string // 龙译团队 ID
	WorksetIndex     uint   // 工作集索引

	SourceLanguage  string   // 源语言
	TargetLanguages []string // 目标语言列表

	AllowApplyType       MoetranAllowApplyType // 允许申请类型
	ApplicationCheckType MoetranAppliCheckType // 申请审核类型

	DefaultRole MoetranRole // 默认角色 ID
}

// moetranCreateProjectRequest 定义了创建项目的请求结构体
type moetranCreateProjectRequest struct {
	Name                 string                `json:"name"`                   // 项目标题
	Intro                string                `json:"intro"`                  // 项目简介
	SourceLanguage       string                `json:"source_language"`        // 源语言
	TargetLanguages      []string              `json:"target_languages"`       // 目标语言列表
	AllowApplyType       MoetranAllowApplyType `json:"allow_apply_type"`       // 允许申请类型
	ApplicationCheckType MoetranAppliCheckType `json:"application_check_type"` // 申请审核类型
	DefaultRole          MoetranRole           `json:"default_role"`           // 默认角色 ID
	ProjectSet           string                `json:"project_set"`            // 项目集 ID
}

// moetranCreateProjectResponse 定义了创建项目的响应结构体
type moetranCreateProjectResponse struct {
	Message string `json:"message"` // 响应消息
	Project struct {
		Id string `json:"id"` // 项目 ID
	} `json:"project"` // 项目详情
}

// CreateProjectSetInfo 定义了创建项目集的基本信息
type CreateProjectSetInfo struct {
	MoetranAuth string // 龙译 JWT

	MoetranTeamId string // 龙译团队 ID
	Name          string // 项目集名称
}

// moetranCreateProjectSetRequest 定义了创建项目集的请求结构体
type moetranCreateProjectSetRequest struct {
	Name string `json:"name"` // 项目集名称
}

// moetranCreateProjectSetResponse 定义了创建项目集的响应结构体
type moetranCreateProjectSetResponse struct {
	Message    string `json:"message"` // 响应消息
	ProjectSet struct {
		Id string `json:"id"` // 项目集 ID
	} `json:"project_set"` // 项目集详情
}

// InviteMemberInfo 定义了邀请成员加入特定项目的基本信息
type InviteMemberInfo struct {
	MoetranAuth string // 龙译 JWT

	MoetranProjectId string      // 龙译项目 ID
	MoetranInviteeID string      // 龙译被邀请成员 ID
	InviteRole       MoetranRole // 邀请的角色 ID
}

// moetranInviteMemberRequest 定义了邀请成员的请求结构体
type moetranInviteMemberRequest struct {
	UserId  string `json:"user_id"` // 被邀请成员的 ID
	RoleId  string `json:"role"`    // 邀请的角色 ID
	Message string `json:"message"` // 邀请消息
}

// moetranInviteMemberResponse 定义了邀请成员的响应结构体
type moetranInviteMemberResponse struct {
	Message string `json:"message"` // 响应消息
}
