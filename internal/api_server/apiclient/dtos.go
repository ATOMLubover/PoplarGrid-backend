package apiclient

// 加入项目申请的处理方式
const (
	APPLI_NON_CHECK   = 1 // 允许所有人直接加入
	APPLI_ADMIN_CHECK = 2 // 仅允许管理员通过申请
)

// 项目的可申请范围
const (
	ALLOW_NO_APPLI    = 1 // 禁止任何形式的申请
	ALLOW_ANY_APPLI   = 2 // 允许任何形式的申请
	ALLOW_MEMBER_ONLY = 3 // 只允许本组成员申请
)

// 本地权限编号
const (
	ROLE_ADMIN       = 1 // 管理员
	ROLE_SUPERVISOR  = 2 // 监理
	ROLE_PROOFREADER = 3 // 校对
	ROLE_TRANSLATOR  = 4 // 翻译
	ROLE_EMBEDDER    = 5 // 嵌字
	ROLE_INTERN      = 6 // 实习翻译
)

// SystemRoleIDsMap 包含从本地权限编号到龙译权限 ID 的映射
var SystemRoleIDsMap = map[int]string{
	ROLE_ADMIN:       "63d87c24b8bebd75ff934264", // 管理员
	ROLE_SUPERVISOR:  "63d87c24b8bebd75ff934265", // 监理
	ROLE_PROOFREADER: "63d87c24b8bebd75ff934266", // 校对
	ROLE_TRANSLATOR:  "63d87c24b8bebd75ff934267", // 翻译
	ROLE_EMBEDDER:    "63d87c24b8bebd75ff934268", // 嵌字
	ROLE_INTERN:      "63d87c24b8bebd75ff934269", // 实习翻译
}

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

	AllowApplyType       int // 允许申请类型
	ApplicationCheckType int // 申请审核类型

	DefaultRole string // 默认角色 ID
}

// MoetranCreateProjectRequest 定义了创建项目的请求结构体
type MoetranCreateProjectRequest struct {
	Name                 string   `json:"name"`                   // 项目标题
	Intro                string   `json:"intro"`                  // 项目简介
	SourceLanguage       string   `json:"source_language"`        // 源语言
	TargetLanguages      []string `json:"target_languages"`       // 目标语言列表
	AllowApplyType       int      `json:"allow_apply_type"`       // 允许申请类型
	ApplicationCheckType int      `json:"application_check_type"` // 申请审核类型
	DefaultRole          string   `json:"default_role"`           // 默认角色 ID
	ProjectSet           string   `json:"project_set"`            // 项目集 ID
}

// MoetranCreateProjectResponse 定义了创建项目的响应结构体
type MoetranCreateProjectResponse struct {
	Message string `json:"message"` // 响应消息
	Project struct {
		Id string `json:"id"` // 项目 ID
	} `json:"project"` // 项目详情
}
